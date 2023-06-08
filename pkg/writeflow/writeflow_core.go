package writeflow

import (
	"context"
	"errors"
	"fmt"
	"github.com/samber/lo"
	"github.com/zbysir/writeflow/internal/cmd"
	"github.com/zbysir/writeflow/internal/model"
	"github.com/zbysir/writeflow/pkg/schema"
	"sort"
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
	Key       string
	Type      NodeInputType // anchor, literal
	Literal   string        // 字面量
	NodeId    string        // anchor node id
	OutputKey string
}

type Node struct {
	Id     string
	Cmd    string
	Inputs []NodeInput
	Enable []NodeInput // 只有当 enable 为 true 时，才会执行该节点
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
				delete(nds, input.NodeId)
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
			inputs = append(inputs, NodeInput{
				Key:       input.Key,
				Type:      "literal",
				Literal:   node.Data.GetInputValue(input.Key),
				NodeId:    "",
				OutputKey: "",
			})
		}

		for _, input := range node.Data.InputAnchors {
			nodeId, outputKey := node.Data.GetInputAnchorValue(input.Key)
			if nodeId == "" && !input.Optional {
				return nil, fmt.Errorf("input '%v' for node '%v' is not defined", input.Key, node.Id)
			}

			inputs = append(inputs, NodeInput{
				Key:       input.Key,
				Type:      "anchor",
				Literal:   "",
				NodeId:    nodeId,
				OutputKey: outputKey,
			})
		}

		cmd := node.Type
		if node.Data.Source.CmdType == model.NothingCmd {
			cmd = model.NothingCmd
		} else if node.Data.Source.BuiltinCmd != "" {
			cmd = node.Data.Source.BuiltinCmd
		}

		nodes[node.Id] = Node{
			Id:     node.Id,
			Cmd:    cmd,
			Inputs: inputs,
		}
	}
	return &Flow{
		Nodes:        nodes,
		OutputNodeId: m.Graph.GetOutputNodeId(),
	}, nil
}

func (f *WriteFlowCore) ExecFlow(ctx context.Context, flow *Flow, initParams map[string]interface{}) (rsp map[string]interface{}, err error) {
	result := make(chan *model.NodeStatus, len(flow.Nodes))

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		err = f.ExecFlowAsync(ctx, flow, initParams, result)
		if err != nil {
			close(result)
			return
		}
		close(result)
	}()

	for r := range result {
		if r.NodeId == flow.OutputNodeId && r.Status == model.StatusSuccess {
			rsp = r.Result
			break
		}
	}

	return
}

func (f *WriteFlowCore) ExecFlowAsync(ctx context.Context, flow *Flow, initParams map[string]interface{}, results chan *model.NodeStatus) (err error) {
	// use INPUT node to get init params
	f.RegisterCmd("INPUT", cmd.NewFun(func(ctx context.Context, _ map[string]interface{}) (map[string]interface{}, error) {
		return initParams, nil
	}))

	fr := newRunner(f.cmds, flow)
	runNodes := flow.Nodes.GetRootNodes()

	for _, node := range runNodes {
		_, err := fr.ExecJob(ctx, node.Id, func(result model.NodeStatus) {
			results <- &result
		})
		if err != nil {
			return err
		}

		//log.Infof("node: %v, rsp: %+v", node.Id, rsp)
	}

	return
}

type runner struct {
	cmd         map[string]schema.CMDer // id -> cmder
	flowDef     *Flow
	cmdRspCache map[string]map[string]interface{}
}

func newRunner(cmd map[string]schema.CMDer, flowDef *Flow) *runner {
	return &runner{cmd: cmd, flowDef: flowDef, cmdRspCache: map[string]map[string]interface{}{}}
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

func (f *runner) ExecJob(ctx context.Context, nodeId string, onNodeRun func(result model.NodeStatus)) (rsp map[string]interface{}, err error) {
	start := time.Now()
	defer func() {
		if onNodeRun != nil {
			if err != nil {
				if err == ErrNodeUnreachable || err.Error() == ErrNodeUnreachable.Error() {
					onNodeRun(model.NodeStatus{
						NodeId: nodeId,
						Status: model.StatusUnreachable,
						Error:  "",
						Result: nil,
						RunAt:  start,
						EndAt:  time.Now(),
					})
				} else {
					onNodeRun(model.NodeStatus{
						NodeId: nodeId,
						Status: model.StatusFailed,
						Error:  err.Error(),
						Result: nil,
						RunAt:  start,
						EndAt:  time.Now(),
					})
				}
			} else {
				onNodeRun(model.NodeStatus{
					NodeId: nodeId,
					Status: model.StatusSuccess,
					Error:  "",
					Result: rsp,
					RunAt:  start,
					EndAt:  time.Now(),
				})
			}
		}
	}()

	nodeDef := f.flowDef.Nodes[nodeId]

	inputs := nodeDef.Inputs
	//log.Infof("input %v: %+v", nodeId, inputs)

	dependValue := map[string]interface{}{}
	var inputKeys []string
	for _, i := range inputs {
		inputKeys = append(inputKeys, i.Key)
		var value interface{}
		switch i.Type {
		case NodeInputLiteral:
			value = i.Literal
		case NodeInputAnchor:
			if i.NodeId == "" {
				// 如果节点 id 为空，则说明是非必填字段。
				continue
			}
			if f.cmdRspCache[i.NodeId] != nil {
				//log.Printf("i.NodeId %v: %+v", i.NodeId, f.cmdRspCache[i.NodeId])
				value = f.cmdRspCache[i.NodeId][i.OutputKey]
			} else {
				rsps, err := f.ExecJob(ctx, i.NodeId, onNodeRun)
				if err != nil {
					return nil, fmt.Errorf("exec task '%s' err: %w", i.NodeId, err)
				}

				value = rsps[i.OutputKey]

				f.cmdRspCache[i.NodeId] = rsps
			}
		}

		dependValue[i.Key] = value
	}

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

	//log.Printf("dependValue: %+v", dependValue)
	cmd := nodeDef.Cmd
	if cmd == "" {
		return nil, NewExecNodeError(fmt.Errorf("cmd is not defined"), nodeDef.Id)
	}
	c, ok := f.cmd[cmd]
	if !ok {
		if cmd == model.NothingCmd {
			return nil, nil
		}
		return nil, NewExecNodeError(fmt.Errorf("cmd '%s' not found", cmd), nodeDef.Id)
	}
	rsp, err = c.Exec(WithInputKeys(ctx, inputKeys), dependValue)
	if err != nil {
		return nil, NewExecNodeError(err, nodeDef.Id)
	}
	return rsp, err
}