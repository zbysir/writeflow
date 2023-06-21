package writeflow

import (
	"context"
	"errors"
	"fmt"
	"github.com/samber/lo"
	"github.com/spf13/cast"
	"github.com/zbysir/writeflow/internal/pkg/keylock"
	"github.com/zbysir/writeflow/internal/pkg/log"
	"io"
	"sort"
	"strings"
	"sync"
	"time"
)

type WriteFlowCore struct {
	cmds map[string]CMDer
}

func NewWriteFlowCore() *WriteFlowCore {
	return &WriteFlowCore{
		cmds: map[string]CMDer{},
	}
}

func (f *WriteFlowCore) RegisterCmd(key string, cmd CMDer) {
	f.cmds[key] = cmd
}

type NodeInputType = string

const (
	NodeInputAnchor  NodeInputType = "anchor"
	NodeInputLiteral NodeInputType = "literal"
)

type NodeInput struct {
	Key      string
	Type     NodeInputType // anchor, literal
	Literal  interface{}   // 字面量
	List     bool
	Required bool

	Anchors []NodeAnchorTarget
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
	BuiltCmd CMDer // go script, js script
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

type Map = map[string]interface{}

func (f *WriteFlowCore) ExecFlowAsync(ctx context.Context, flow *Flow, initParams map[string]interface{}, parallel int) (results chan NodeStatusLog, err error) {
	// use params node to get init params
	f.RegisterCmd("_params", NewFun(func(ctx context.Context, _ Map) (Map, error) {
		return NewMap(initParams), nil
	}))

	fr := newRunner(f.cmds, flow, parallel)
	rootNodes := flow.Nodes.GetRootNodes()

	results = make(chan NodeStatusLog, 100)
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
				_, _ = fr.ExecNode(ctx, node.Id, false, func(result NodeStatusLog) {
					results <- result
				})
			}()
		}

		wg.Wait()
	}()

	return
}

func (f *WriteFlowCore) ExecNode(ctx context.Context, flow *Flow, initParams map[string]interface{}, parallel int) (rsp Map, err error) {
	// use params node to get init params
	f.RegisterCmd("_params", NewFun(func(ctx context.Context, _ Map) (Map, error) {
		return NewMap(initParams), nil
	}))

	fr := newRunner(f.cmds, flow, parallel)
	_, ok := flow.Nodes[flow.OutputNodeId]
	if !ok {
		return Map{}, fmt.Errorf("output node %s not found", flow.OutputNodeId)
	}

	return fr.ExecNode(ctx, flow.OutputNodeId, false, nil)
}

type runner struct {
	flowDef     *Flow
	cmd         map[string]CMDer                  // id -> cmder
	cmdRspCache map[string]*runnerRsp             // nodeId->key->value
	inject      map[string]map[string]interface{} // nodeId->key->value
	l           sync.RWMutex                      // lock for map
	keyLock     *keylock.KeyLock                  // lock for cmdRspCache (防止并发下缓存穿透)
	limitChan   chan struct{}
}
type runnerRsp struct {
	rsp Map
	err error
}

func (r *runner) getRspCache(nodeId string, key string) (v interface{}, exist bool, err error) {
	r.l.RLock()
	defer r.l.RUnlock()

	if r.cmdRspCache[nodeId] == nil {
		return nil, false, nil
	}
	err = r.cmdRspCache[nodeId].err
	if r.cmdRspCache[nodeId].rsp == nil {
		return nil, true, err
	}
	v, exist = r.cmdRspCache[nodeId].rsp[key]
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

func (r *runner) setRspCache(nodeId string, rsp Map, err error) {
	r.l.Lock()
	defer r.l.Unlock()

	r.cmdRspCache[nodeId] = &runnerRsp{rsp: rsp, err: err}
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

func newRunner(cmd map[string]CMDer, flowDef *Flow, parallel int) *runner {
	return &runner{
		flowDef:     flowDef,
		cmd:         cmd,
		cmdRspCache: map[string]*runnerRsp{},
		inject:      map[string]map[string]interface{}{},
		l:           sync.RWMutex{},
		keyLock:     keylock.NewKeyLock(),
		limitChan:   make(chan struct{}, parallel),
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

func NewMap(s map[string]interface{}) Map {
	if s == nil {
		return map[string]interface{}{}
	}
	return s
}

type Stack struct {
	Nodes string

	// stack
}

func (s Stack) Push(nodeId string) Stack {
	return Stack{
		Nodes: s.Nodes + "->" + nodeId,
	}
}

func (f *runner) ExecNode(ctx context.Context, nodeId string, nocache bool, onNodeStatusChange func(result NodeStatusLog)) (rsp Map, err error) {
	start := time.Now()
	skipEmitChange := false
	defer func() {
		if !skipEmitChange && onNodeStatusChange != nil {
			if err != nil {
				onNodeStatusChange(NewNodeStatusLog(nodeId, StatusFailed, err.Error(), Map{}, start, time.Now()))
			} else {
				onNodeStatusChange(NewNodeStatusLog(nodeId, StatusSuccess, "", rsp, start, time.Now()))
			}
		}
	}()

	nodeDef, ok := f.flowDef.Nodes[nodeId]
	if !ok {
		// 可能是前段没有删除干净，所以有空，先忽略错误
		return Map{}, nil
	}
	inputs := nodeDef.Inputs

	var calcInput = func(i NodeInput, nocache bool) (interface{}, error) {
		switch i.Type {
		case NodeInputLiteral:
			return i.Literal, nil
		case NodeInputAnchor:
			if len(i.Anchors) == 0 && i.Required {
				return nil, NewExecNodeError(fmt.Errorf("params '%v' is required", i.Key), nodeId)
			}

			var values []interface{}

			for _, i := range i.Anchors {
				v, ok := f.getInject(i.NodeId, i.OutputKey)
				if ok {
					return v, nil
				}

				lockKey := fmt.Sprintf("%s", i.NodeId)

				// 加锁防止缓存穿透，这里有递归调用, 只能使用 tryLock
				//log.Infof("----lock %s", lockKey)
				for {
					lock := f.keyLock.TryLock(lockKey)
					if lock {
						//log.Infof("---- lock %s", lockKey)
						break
					}
					//log.Infof("----wait lock %s", lockKey)
					time.Sleep(time.Millisecond * 100)
				}

				if !nocache {
					v, ok, err = f.getRspCache(i.NodeId, i.OutputKey)
					if err != nil {
						//log.Infof("----Unlock %s", lockKey)
						f.keyLock.Unlock(lockKey)
						return nil, err
					}
					if ok {
						values = append(values, v)
						f.keyLock.Unlock(lockKey)
					} else {
						rsps, err := f.ExecNode(ctx, i.NodeId, nocache, onNodeStatusChange)
						if err != nil {
							//log.Infof("----Unlock %s", lockKey)
							f.setRspCache(i.NodeId, Map{}, err)
							f.keyLock.Unlock(lockKey)
							return nil, err
						}
						if rsps == nil {
							f.setRspCache(i.NodeId, rsps, nil)
							f.keyLock.Unlock(lockKey)
							continue
						}

						value := rsps[i.OutputKey]
						values = append(values, value)

						f.setRspCache(i.NodeId, rsps, nil)
						//log.Infof("----Unlock %s", lockKey)
						f.keyLock.Unlock(lockKey)
					}
				} else {
					rsps, err := f.ExecNode(ctx, i.NodeId, nocache, onNodeStatusChange)
					if err != nil {
						//log.Infof("----Unlock %s", lockKey)
						f.setRspCache(i.NodeId, Map{}, err)
						f.keyLock.Unlock(lockKey)
						return nil, err
					}
					if rsps == nil {
						f.setRspCache(i.NodeId, rsps, nil)
						f.keyLock.Unlock(lockKey)
						continue
					}

					value, _ := rsps[i.OutputKey]
					values = append(values, value)
					f.keyLock.Unlock(lockKey)
				}
			}

			if i.List {
				return values, nil
			} else if len(values) >= 1 {
				return values[0], nil
			} else {
				return nil, nil
			}
		}

		return nil, nil
	}

	// 当 _enable 为 false 时，才会跳过节点。
	enableInput, inputs, ok := inputs.PopKey("_enable")
	if ok {
		enable, err := calcInput(enableInput, nocache)
		if err != nil {
			return Map{}, err
		}
		if cast.ToBool(cast.ToString(enable)) == false {
			return Map{}, nil
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
				return Map{}, err
			}
		}

		for _, input := range inputs {
			condition := input.Key

			v, err := LookInterface(map[string]interface{}{"data": data}, condition)
			if err != nil {
				return Map{}, NewExecNodeError(fmt.Errorf("exec condition %s error: %w", condition, err), nodeId)
			}

			// ToBool can't convert int64 to bool
			if cast.ToBool(cast.ToString(v)) {
				r, err := calcInput(input, nocache)
				if err != nil {
					return Map{}, err
				}

				return NewMap(map[string]interface{}{"default": r, "branch": input.Key}), nil
			}
		}

		rsp = NewMap(map[string]interface{}{"default": nil, "branch": ""})
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
					return Map{}, err
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
			return Map{}, NewExecNodeError(fmt.Errorf("for %T error: %w", data, err), nodeId)
		}

		if forError != nil {
			return Map{}, NewExecNodeError(fmt.Errorf("for %T error: %w", data, forError), nodeId)
		}

		rsp = NewMap(map[string]interface{}{"default": rsps})
	default:
		dependValue := NewMap(nil)
		var inputKeys []string

		var wg sync.WaitGroup
		var ml sync.Mutex
		var calcErr error
		for _, i := range inputs {
			inputKeys = append(inputKeys, i.Key)

			// 并发执行
			// 由于是递归，不方便控制节点执行数量，而是控制协程数量（不包括主协程）。
			select {
			case f.limitChan <- struct{}{}:
				wg.Add(1)
				go func(i NodeInput) {
					defer func() {
						wg.Done()
						<-f.limitChan
					}()
					if calcErr != nil {
						return
					}

					r, err := calcInput(i, nocache)
					if err != nil {
						calcErr = err
					} else {
						ml.Lock()
						dependValue[i.Key] = r
						ml.Unlock()
					}
				}(i)
			default:
				r, err := calcInput(i, nocache)
				if err != nil {
					log.Errorf("calcInput %v error: %v", nodeId, err)
					return Map{}, err
				}
				ml.Lock()
				dependValue[i.Key] = r
				ml.Unlock()
			}
		}
		wg.Wait()

		if calcErr != nil {
			return Map{}, calcErr
		}

		// 流式输出特殊处理，异步读取输入的流并同步状态。
		if nodeDef.Cmd == "_output" {
			var wg sync.WaitGroup
			dependValuex := map[string]interface{}{}
			valueLock := sync.Mutex{}
			for k, v := range dependValue {
				dependValuex[k] = v
			}
			for k, v := range dependValue {
				if steam, ok := v.(*StreamResponse[string]); ok {
					wg.Add(1)
					valueLock.Lock()
					dependValuex[k] = ""
					valueLock.Unlock()
					reader := steam.NewReader()
					go func() {
						defer wg.Done()
						var allContent string
						for {
							var t string
							t, err = reader.Read()
							if err != nil {
								if err == io.EOF {
									err = nil
								}
								break
							} else {
								allContent += t

								valueLock.Lock()
								dependValuex[k] = allContent
								valueLock.Unlock()

								onNodeStatusChange(NewNodeStatusLog(nodeId, StatusRunning, "", dependValuex, start, time.Now()))
							}
						}
					}()
				}
			}

			wg.Wait()
			if err == nil {
				onNodeStatusChange(NewNodeStatusLog(nodeId, StatusSuccess, "", dependValuex, start, time.Now()))
			} else {
				onNodeStatusChange(NewNodeStatusLog(nodeId, StatusFailed, err.Error(), dependValuex, start, time.Now()))
			}

			skipEmitChange = true
			return nil, err
		} else {
			// 只有自定义 cmd 才需要报告 running 状态，特殊的 _for, _switch 不需要。
			if onNodeStatusChange != nil {
				onNodeStatusChange(NewNodeStatusLog(nodeId, StatusRunning, "", Map{}, start, time.Time{}))
			}

			// 其他命令需要等待直到流完成
			for k, v := range dependValue {
				if steam, ok := v.(*StreamResponse[string]); ok {
					var ts []string
					ts, err = steam.NewReader().ReadAll()
					if err != nil {
						break
					}
					if len(ts) != 0 {
						dependValue[k] = strings.Join(ts, "")
					} else {
						dependValue[k] = ""
					}
				}
			}
			if err != nil {
				return Map{}, NewExecNodeError(err, nodeDef.Id)
			}
		}

		cmdName := nodeDef.Cmd
		if cmdName == "" {
			return Map{}, NewExecNodeError(fmt.Errorf("cmd is not defined"), nodeDef.Id)
		}
		cmder := nodeDef.BuiltCmd
		if cmder == nil {
			var ok bool
			cmder, ok = f.cmd[cmdName]
			if !ok {
				if cmdName == "nothing" {
					// 如果不需要执行任何命令，则直接返回 input
					return dependValue, nil
				}
				return Map{}, NewExecNodeError(fmt.Errorf("cmd '%s' not found", cmdName), nodeDef.Id)
			}
		}

		rsp, err = HandlePanicCmd(cmder).Exec(WithInputKeys(ctx, inputKeys), dependValue)
		if err != nil {
			return Map{}, NewExecNodeError(err, nodeDef.Id)
		}
	}

	//log.Printf("dependValue: %+v", dependValue)

	return rsp, nil
}
