package writeflow

import (
	"context"
	"errors"
	"fmt"
	"github.com/samber/lo"
	"github.com/spf13/cast"
	"github.com/zbysir/writeflow/internal/cmd"
	"github.com/zbysir/writeflow/internal/model"
	"github.com/zbysir/writeflow/pkg/schema"
	"sort"
	"sync"
	"time"
)

type WriteFlowCore struct {
	cmds map[string]schema.CMDer
}

func NewWriteFlowCore() *WriteFlowCore {
	return &WriteFlowCore{
		cmds: map[string]schema.CMDer{},
	}
}

func (f *WriteFlowCore) RegisterCmd(key string, cmd schema.CMDer) {
	f.cmds[key] = cmd
}

type NodeInputType = string

const (
	NodeInputAnchor  NodeInputType = "anchor"
	NodeInputLiteral NodeInputType = "literal"
)

type NodeInput struct {
	Key     string
	Type    NodeInputType // anchor, literal
	Literal interface{}   // 字面量
	List    bool

	Anchors []model.NodeAnchorTarget
}

type NodeInputs []NodeInput

func (n NodeInputs) PopKey(key string) (NodeInput, NodeInputs, bool) {
	for i, v := range n {
		if v.Key == key {
			nn := make([]NodeInput, len(n)-1)
			copy(nn, n[:i])
			copy(nn[i:], n[i+1:])
			return v, nn, true
		}
	}

	return NodeInput{}, n, false
}

type NodeAnchorTarget struct {
	NodeId    string `json:"node_id"`    // 关联的节点 id
	OutputKey string `json:"output_key"` // 关联的节点输出 key
}

type Node struct {
	Id       string
	Cmd      string
	BuiltCmd schema.CMDer
	Inputs   NodeInputs
}

type ForItemNode struct {
	NodeId    string
	InputKey  string
	OutputKey string // outputKey 可不填，默认等于 inputKey
}

type inputKeysKey struct{}

func WithInputKeys(ctx context.Context, inputKeys []string) context.Context {
	return context.WithValue(ctx, inputKeysKey{}, inputKeys)
}

func GetInputKeys(ctx context.Context) []string {
	if v, ok := ctx.Value(inputKeysKey{}).([]string); ok {
		return v
	}
	return nil
}

type Nodes map[string]Node
type Flow struct {
	Nodes        Nodes // node id -> node
	OutputNodeId string
}

// GetRootNodes Get root nodes that need run
func (d Nodes) GetRootNodes() (nodes []Node) {
	nds := map[string]Node{}
	for k, v := range d {
		nds[k] = v
	}

	for _, v := range d {
		for _, input := range v.Inputs {
			if input.Type == NodeInputAnchor {
				for _, v := range input.Anchors {
					delete(nds, v.NodeId)
				}
			}
		}
	}

	// sort for stable
	var keys []string
	for k := range nds {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		nodes = append(nodes, nds[key])
	}
	return nodes
}

func (d *Flow) UsedComponents() (componentType []string) {
	for _, v := range d.Nodes {
		componentType = append(componentType, v.Cmd)
	}
	componentType = lo.Uniq(componentType)

	return componentType
}

func FlowFromModel(m *model.Flow) (*Flow, error) {
	nodes := map[string]Node{}

	for _, node := range m.Graph.Nodes {
		var inputs []NodeInput
		for _, input := range node.Data.InputParams {
			switch input.InputType {

			case model.NodeInputTypeAnchor:
				anchors, list := node.Data.GetInputAnchorValue(input.Key)
				if len(anchors) == 0 && !input.Optional {
					return nil, fmt.Errorf("input '%v' for node '%v' is not defined", input.Key, node.Id)
				}

				inputs = append(inputs, NodeInput{
					Key:     input.Key,
					Type:    "anchor",
					Literal: "",
					Anchors: anchors,
					List:    list,
				})
			default:
				inputs = append(inputs, NodeInput{
					Key:     input.Key,
					Type:    NodeInputLiteral,
					Literal: node.Data.GetInputValue(input.Key),
				})
			}
		}

		// 废弃后可以删掉
		for _, input := range node.Data.InputAnchors {
			anchors, list := node.Data.GetInputAnchorValue(input.Key)
			if len(anchors) == 0 && !input.Optional {
				return nil, fmt.Errorf("input '%v' for node '%v' is not defined", input.Key, node.Id)
			}

			inputs = append(inputs, NodeInput{
				Key:     input.Key,
				Type:    "anchor",
				Literal: "",
				Anchors: anchors,
				List:    list,
			})
		}

		cmdName := node.Type
		var cmder schema.CMDer
		switch node.Data.Source.CmdType {
		case model.NothingCmd:
			cmdName = model.NothingCmd
		case model.GoScriptCmd:
			key := node.Data.Source.Script.InputKey
			if key == "" {
				key = "script"
			}
			script := node.Data.GetInputValue(key)
			var err error
			cmder, err = cmd.NewGoScript(nil, "", script)
			if err != nil {
				return nil, fmt.Errorf("parse node '%v' script '%v' error: %v", node.Id, script, err)
			}
		case model.JavaScriptCmd:
			key := node.Data.Source.Script.InputKey
			if key == "" {
				key = "script"
			}
			script := node.Data.GetInputValue(key)
			var err error
			cmder, err = cmd.NewJavaScript(script)
			if err != nil {
				return nil, fmt.Errorf("parse node '%v' script '%v' error: %v", node.Id, script, err)
			}
		}
		if node.Data.Source.CmdType == model.NothingCmd {
			cmdName = model.NothingCmd
		} else if node.Data.Source.BuiltinCmd != "" {
			cmdName = node.Data.Source.BuiltinCmd
		}

		nodes[node.Id] = Node{
			Id:       node.Id,
			Cmd:      cmdName,
			BuiltCmd: cmder,
			Inputs:   inputs,
		}
	}
	return &Flow{
		Nodes:        nodes,
		OutputNodeId: m.Graph.GetOutputNodeId(),
	}, nil
}

func (f *WriteFlowCore) ExecFlow(ctx context.Context, flow *Flow, initParams map[string]interface{}) (rsp map[string]interface{}, err error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	result, err := f.ExecFlowAsync(ctx, flow, initParams)
	if err != nil {
		return nil, err
	}

	for r := range result {
		if r.NodeId == flow.OutputNodeId && r.Status == model.StatusSuccess {
			rsp = r.Result
			break
		}
	}

	return
}

func (f *WriteFlowCore) ExecFlowAsync(ctx context.Context, flow *Flow, initParams map[string]interface{}, ) (results chan *model.NodeStatus, err error) {
	// use params node to get init params
	f.RegisterCmd("_params", cmd.NewFun(func(ctx context.Context, _ map[string]interface{}) (map[string]interface{}, error) {
		return initParams, nil
	}))

	fr := newRunner(f.cmds, flow)
	rootNodes := flow.Nodes.GetRootNodes()

	results = make(chan *model.NodeStatus, 100)
	go func() {
		defer func() {
			close(results)
		}()
		var wg sync.WaitGroup
		for _, node := range rootNodes {
			node := node
			wg.Add(1)
			go func() {
				defer wg.Done()
				_, _ = fr.ExecNode(ctx, node.Id, false, func(result model.NodeStatus) {
					results <- &result
				})
			}()
		}

		wg.Wait()
	}()

	return
}

type runner struct {
	flowDef     *Flow
	cmd         map[string]schema.CMDer           // id -> cmder
	cmdRspCache map[string]map[string]interface{} // nodeId->key->value
	inject      map[string]map[string]interface{} // nodeId->key->value
	l           sync.RWMutex                      // lock for map
}

func (r *runner) getRspCache(nodeId string, key string) (v interface{}, exist bool) {
	r.l.RLock()
	defer r.l.RUnlock()

	if r.cmdRspCache[nodeId] == nil {
		return nil, false
	}
	v, exist = r.cmdRspCache[nodeId][key]
	return
}

func (r *runner) getInject(nodeId string, key string) (v interface{}, exist bool) {
	r.l.RLock()
	defer r.l.RUnlock()

	if r.inject[nodeId] == nil {
		return nil, false
	}
	v, exist = r.inject[nodeId][key]
	return
}

func (r *runner) setRspCache(nodeId string, rsp map[string]interface{}) {
	r.l.Lock()
	defer r.l.Unlock()

	r.cmdRspCache[nodeId] = rsp
	return
}

func (r *runner) setRspItemCache(nodeId string, k string, v interface{}) {
	r.l.Lock()
	defer r.l.Unlock()

	if r.cmdRspCache[nodeId] == nil {
		r.cmdRspCache[nodeId] = map[string]interface{}{}
	}

	r.cmdRspCache[nodeId][k] = v
	return
}

func (r *runner) setInject(nodeId string, k string, v interface{}) {
	r.l.Lock()
	defer r.l.Unlock()

	if r.inject[nodeId] == nil {
		r.inject[nodeId] = map[string]interface{}{}
	}

	r.inject[nodeId][k] = v
	return
}

func newRunner(cmd map[string]schema.CMDer, flowDef *Flow) *runner {
	return &runner{
		cmd:         cmd,
		flowDef:     flowDef,
		cmdRspCache: map[string]map[string]interface{}{},
		inject:      map[string]map[string]interface{}{},
	}
}

type ExecNodeError struct {
	Cause  error
	NodeId string
}

func NewExecNodeError(cause error, nodeId string) *ExecNodeError {
	return &ExecNodeError{Cause: cause, NodeId: nodeId}
}

func (e *ExecNodeError) Error() string {
	return fmt.Sprintf("exec node '%s' err: %v", e.NodeId, e.Cause)
}
func (e *ExecNodeError) Unwrap() error {
	return e.Cause
}

var ErrNodeUnreachable = errors.New("node unreachable")

// DGA 可以转为一个表达式。
// 例如：
//
// LangChain({
//   openai: OpenAI({ key: "xx" }).default,
//   prompt: InputStr({ default: "Hi!" }).default
// })
// 这个函数可以通过递归执行依赖的函数（节点）来得到结果。
// 足够简单并且只会执行依赖的函数。
//
// 但要描述 If，和 For 语句就不如 脚本 方便。
//
// Switch:
//
// LangChain({
//   openai: Switch({
//     data: InputStr({default: 'openai'}).default,
//     conditions: [
//        { exp: "data==openai", value: OpenAI({ key: "xx" }).default },
//        { exp: "data==local", value: LocalLLM({ path: "xx" }).default },
//     ],
//   }),
//   prompt: InputStr({ default: "Hi!" }).default
// })
//
// For:
//
// LangChain({
//   openai: OpenAI({ key: "xx" }).default,
//   prompts: For({
//     data: GetList().default,
//     item: AddPrefix({item: <item>, prefix: "Hi: "}).default,
//   })
// })

// Switch 和 For 不能使用 Cmd(params map[string]interface{}) map[string]interface{} 实现，而是需要内置。
// Cmd 依赖的是已经处理好的值，而 Switch 和 For 需要依赖懒值（函数），只有当需要的时候才会执行。
// 如果让 Cmd 处理懒值会导致 Cmd 的编写逻辑变得复杂，同时还需要处理函数执行异常，不方便用户编写。
// 而逻辑分支相对固定，可以内置实现。

func (f *runner) ExecNode(ctx context.Context, nodeId string, nocache bool, onNodeRun func(result model.NodeStatus)) (rsp map[string]interface{}, err error) {
	start := time.Now()
	defer func() {
		if onNodeRun != nil {
			if err != nil {
				if errors.Is(err, ErrNodeUnreachable) || err.Error() == ErrNodeUnreachable.Error() {
					onNodeRun(model.NodeStatus{
						NodeId: nodeId,
						Status: model.StatusUnreachable,
						RunAt:  start,
						EndAt:  time.Now(),
					})
					err = nil
				} else {
					onNodeRun(model.NodeStatus{
						NodeId: nodeId,
						Status: model.StatusFailed,
						Error:  err.Error(),
						RunAt:  start,
						EndAt:  time.Now(),
					})
				}
			} else {
				onNodeRun(model.NodeStatus{
					NodeId: nodeId,
					Status: model.StatusSuccess,
					Result: rsp,
					RunAt:  start,
					EndAt:  time.Now(),
				})
			}
		} else {
			if err == ErrNodeUnreachable || err.Error() == ErrNodeUnreachable.Error() {
				err = nil
			}
		}
	}()

	nodeDef := f.flowDef.Nodes[nodeId]
	inputs := nodeDef.Inputs

	var calcInput = func(i NodeInput, nocache bool) (interface{}, error) {
		switch i.Type {
		case NodeInputLiteral:
			return i.Literal, nil
		case NodeInputAnchor:
			if len(i.Anchors) == 0 {
				// 如果节点 id 为空，则说明是非必填字段。
				return nil, nil
			}
			var values []interface{}

			for _, i := range i.Anchors {
				v, ok := f.getInject(i.NodeId, i.OutputKey)
				if ok {
					return v, nil
				}

				v, ok = f.getRspCache(i.NodeId, i.OutputKey)
				if ok && !nocache {
					values = append(values, v)
				} else {
					rsps, err := f.ExecNode(ctx, i.NodeId, nocache, onNodeRun)
					if err != nil {
						return nil, err
					}
					values = append(values, rsps[i.OutputKey])

					f.setRspCache(i.NodeId, rsps)
				}
			}

			if i.List {
				return values, nil
			} else {
				return values[0], nil
			}
		}

		return nil, nil
	}

	// 当 _enable 为 false 时，才会跳过节点。
	enableInput, inputs, ok := inputs.PopKey("_enable")
	if ok {
		enable, err := calcInput(enableInput, nocache)
		if err != nil {
			return nil, err
		}
		if cast.ToBool(cast.ToString(enable)) == false {
			return nil, nil
		}
	}

	//log.Infof("input %v: %+v", nodeId, inputs)
	// switch 和 for 内置实现，不使用 cmd 逻辑。
	switch nodeDef.Cmd {
	case "_switch":
		// get data
		var data interface{}
		dataInput, inputs, ok := inputs.PopKey("data")
		if ok {
			data, err = calcInput(dataInput, nocache)
			if err != nil {
				return nil, err
			}
		}

		for _, input := range inputs {
			condition := input.Key

			v, err := LookInterface(data, condition)
			if err != nil {
				return nil, NewExecNodeError(fmt.Errorf("exec condition %s error: %w", condition, err), nodeId)
			}

			// ToBool can't convert int64 to bool
			if cast.ToBool(cast.ToString(v)) {
				r, err := calcInput(input, nocache)
				if err != nil {
					return nil, err
				}

				return map[string]interface{}{"default": r, "branch": input.Key}, nil
			}
		}

		rsp = map[string]interface{}{"default": nil, "branch": ""}
	case "_for":
		// for 的执行逻辑：
		//  找到 input 依赖 forOutput.item 的节点，然后执行它，会递归执行到 forInput 节点，然后 forInput 会返回迭代值。
		// get data
		var data interface{}
		var itemInput NodeInput
		for _, input := range inputs {
			if input.Key == "data" {
				data, err = calcInput(input, nocache)
				if err != nil {
					return nil, err
				}
			} else if input.Key == "item" {
				itemInput = input
			}
		}

		var rsps []interface{}
		var forError error
		err := ForInterface(data, func(i interface{}) {
			if forError != nil {
				return
			}
			f.setInject(nodeId, "item", i)
			r, err := calcInput(itemInput, true)
			if err != nil {
				forError = err
			} else {
				rsps = append(rsps, r)
			}
		})

		if err != nil {
			return nil, NewExecNodeError(fmt.Errorf("for %T error: %w", data, err), nodeId)
		}

		if forError != nil {
			return nil, NewExecNodeError(fmt.Errorf("for %T error: %w", data, forError), nodeId)
		}

		rsp = map[string]interface{}{"default": rsps}
	default:
		dependValue := map[string]interface{}{}
		var inputKeys []string
		for _, i := range inputs {
			inputKeys = append(inputKeys, i.Key)

			r, err := calcInput(i, nocache)
			if err != nil {
				return nil, err
			}
			dependValue[i.Key] = r
		}

		cmdName := nodeDef.Cmd
		if cmdName == "" {
			return nil, NewExecNodeError(fmt.Errorf("cmd is not defined"), nodeDef.Id)
		}
		c := nodeDef.BuiltCmd
		if c == nil {
			var ok bool
			c, ok = f.cmd[cmdName]
			if !ok {
				if cmdName == model.NothingCmd {
					// 如果不需要执行任何命令，则直接返回 input
					return dependValue, nil
				}
				return nil, NewExecNodeError(fmt.Errorf("cmd '%s' not found", cmdName), nodeDef.Id)
			}
		}

		// 只有自定义 cmd 才需要报告 running 状态，特殊的 _for, _switch 不需要。
		if onNodeRun != nil {
			onNodeRun(model.NodeStatus{
				NodeId: nodeId,
				Status: model.StatusRunning,
				Error:  "",
				Result: nil,
				RunAt:  start,
				EndAt:  time.Time{},
			})
		}

		rsp, err = c.Exec(WithInputKeys(ctx, inputKeys), dependValue)
		if err != nil {
			return nil, NewExecNodeError(err, nodeDef.Id)
		}
	}

	//log.Printf("dependValue: %+v", dependValue)

	return rsp, err
}
