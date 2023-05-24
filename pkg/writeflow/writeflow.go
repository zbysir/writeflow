package writeflow

import (
	"context"
	"fmt"
	"github.com/samber/lo"
	"github.com/zbysir/writeflow/internal/cmd"
	"github.com/zbysir/writeflow/internal/model"
	"github.com/zbysir/writeflow/pkg/schema"
	"strings"
)

type WriteFlow struct {
	cmds map[string]*Component
}

func NewWriteFlow() *WriteFlow {
	return &WriteFlow{
		cmds: map[string]*Component{},
	}
}

func (f *WriteFlow) RegisterComponent(cmd *Component) {
	key := cmd.Schema.Key
	f.cmds[key] = cmd
}

// 所有的依赖可以并行计算。
// 这是通过代码逻辑不好描述的
//
// appendName-1:
//   cmd: appendName
//   input:
//     - _args[0]
//     - _args[1]
// hello:
//   cmd: hello
//   input:
//     - appendName-1[0]
//
// END:
//   input:
//     - hello[0]

// flow: 流程定义
// job: flow 由 多个 job 组成
// cmd: job 可以调用 cmd

type NodeInput struct {
	Key         string
	Type        string // anchor, literal
	Literal     string // format: node_id.response_key
	NodeId      string
	ResponseKey string
}

type Node struct {
	Id           string
	ComponentKey string
	Inputs       []NodeInput
}

type Flow struct {
	Nodes map[string]Node // node id -> node
}

func (d *Flow) UsedComponents() (componentKeys []string) {
	for _, v := range d.Nodes {
		componentKeys = append(componentKeys, v.ComponentKey)
	}
	componentKeys = lo.Uniq(componentKeys)

	return componentKeys
}

func FromModelFlow(m *model.Flow) (*Flow, error) {
	nodes := map[string]Node{}

	for _, node := range m.Graph.Nodes {
		var inputs []NodeInput
		for _, input := range node.Data.InputParams {
			inputs = append(inputs, NodeInput{
				Key:         input.Key,
				Type:        "literal",
				Literal:     node.Data.Inputs[input.Key],
				NodeId:      "",
				ResponseKey: "",
			})
		}

		for _, input := range node.Data.InputAnchors {
			ss := strings.Split(node.Data.Inputs[input.Key], ".")
			var nodeId string
			var responseKey string
			if len(ss) > 1 {
				nodeId = ss[0]
				responseKey = ss[1]
			}

			inputs = append(inputs, NodeInput{
				Key:         input.Key,
				Type:        "anchor",
				Literal:     "",
				NodeId:      nodeId,
				ResponseKey: responseKey,
			})
		}

		nodes[node.Id] = Node{
			Id:           node.Id,
			ComponentKey: node.Type,
			Inputs:       inputs,
		}
	}
	return &Flow{
		Nodes: nodes,
	}, nil
}

// SpanInterface 特殊语法，返回值
type SpanInterface []interface{}

type YFlow struct {
	Version string          `yaml:"version"`
	Flow    map[string]YJob `yaml:"flow"`
}

type YJob struct {
	Cmd     string                 `yaml:"cmd"`
	Inputs  map[string]interface{} `yaml:"inputs"`
	Depends []string               `yaml:"depends"`
}

func (j *YJob) ToJobDef(name string) Node {
	var inputs []NodeInput
	for key, item := range j.Inputs {
		switch item := item.(type) {
		case string:
			// _args[1]
			ss := strings.Split(item, ".")
			taskName := ""
			var respField string
			if len(ss) == 2 {
				taskName = ss[0]
				respField = ss[1]
			} else {
				taskName = ss[0]
				respField = "default"
			}

			inputs = append(inputs, NodeInput{
				NodeId:      taskName,
				ResponseKey: respField,
				Key:         key,
			})
		case map[string]interface{}:
			// {name: args[0]}
			// TODO object
		}
	}

	return Node{
		Id:           name,
		ComponentKey: j.Cmd,
		Inputs:       inputs,
	}
}

func (f *WriteFlow) ExecFlow(ctx context.Context, flow *Flow, outputNodeId string, initParams map[string]interface{}) (rsp map[string]interface{}, err error) {
	// use INPUT node to get init params
	f.RegisterComponent(NewComponent(cmd.NewFun(func(ctx context.Context, _ map[string]interface{}) (map[string]interface{}, error) {
		return initParams, nil
	}), cmd.Schema{
		Key: "INPUT",
	}))
	cmds := map[string]schema.CMDer{}
	for k, v := range f.cmds {
		cmds[k] = v.Cmder
	}

	if outputNodeId == "" {
		outputNodeId = "OUTPUT"
	}

	fr := newRunner(cmds, flow)
	rsp, err = fr.ExecJob(ctx, outputNodeId)
	if err != nil {
		return
	}

	return
}

func (f *WriteFlow) GetCMDs(ctx context.Context, names []string) (rsp []cmd.Schema, err error) {
	for _, cmd := range f.cmds {
		rsp = append(rsp, cmd.Schema)
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

func (f *runner) ExecJob(ctx context.Context, jobName string) (rsp map[string]interface{}, err error) {
	jobDef := f.flowDef.Nodes[jobName]

	inputs := jobDef.Inputs
	//log.Printf("input %v: %+v", jobName, inputs)

	dependValue := map[string]interface{}{}
	for _, i := range inputs {
		var value interface{}
		switch i.Type {
		case "literal":
			value = i.Literal
		case "anchor":
			if f.cmdRspCache[i.NodeId] != nil {
				//log.Printf("i.NodeId %v: %+v", i.NodeId, f.cmdRspCache[i.NodeId])
				value = f.cmdRspCache[i.NodeId][i.ResponseKey]
			} else {
				rsps, err := f.ExecJob(ctx, i.NodeId)
				if err != nil {
					return nil, fmt.Errorf("exec task '%s' err: %v", i.NodeId, err)
				}

				value = rsps[i.ResponseKey]

				f.cmdRspCache[i.NodeId] = rsps
			}
		}

		dependValue[i.Key] = value
	}

	//log.Printf("dependValue: %+v", dependValue)
	cmd := jobDef.ComponentKey
	if cmd == "" {
		cmd = jobName
	}
	c, ok := f.cmd[cmd]
	if !ok {
		return nil, fmt.Errorf("cmd '%s' not found", cmd)
	}
	rsp, err = c.Exec(ctx, dependValue)
	if err != nil {
		return nil, fmt.Errorf("exec cmd '%s' err: %v", cmd, err)
	}
	return rsp, err
}
