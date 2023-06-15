package model

import (
	"fmt"
	"github.com/zbysir/writeflow/internal/cmd"
	"github.com/zbysir/writeflow/pkg/schema"
	"github.com/zbysir/writeflow/pkg/writeflow"
)

func FlowFromModel(m *Flow) (*writeflow.Flow, error) {
	nodes := map[string]writeflow.Node{}

	for _, node := range m.Graph.Nodes {
		var inputs []writeflow.NodeInput
		for _, input := range node.Data.InputParams {
			switch input.InputType {
			case writeflow.NodeInputAnchor:
				anchors, list := node.Data.GetInputAnchorValue(input.Key)
				inputs = append(inputs, writeflow.NodeInput{
					Key:      input.Key,
					Type:     writeflow.NodeInputAnchor,
					Literal:  "",
					Anchors:  anchors,
					List:     list,
					Required: !input.Optional,
				})
			default:
				inputs = append(inputs, writeflow.NodeInput{
					Key:      input.Key,
					Type:     writeflow.NodeInputLiteral,
					Literal:  node.Data.GetInputValue(input.Key),
					Required: !input.Optional,
				})
			}
		}

		cmdName := node.Type
		var cmder schema.CMDer
		switch node.Data.Source.CmdType {
		case writeflow.NothingCmd:
			cmdName = string(writeflow.NothingCmd)
		case writeflow.GoScriptCmd:
			var script string
			if node.Data.Source.Script.Source != "" {
				script = node.Data.Source.Script.Source
			} else {
				script = node.Data.GetInputValue("script").(string)
			}

			var err error
			cmder, err = cmd.NewGoScript(nil, "", script)
			if err != nil {
				return nil, writeflow.NewExecNodeError(fmt.Errorf("parse script error: %v", err), node.Id)
			}
		case writeflow.JavaScriptCmd:
			var script string
			if node.Data.Source.Script.Source != "" {
				script = node.Data.Source.Script.Source
			} else {
				script = node.Data.GetInputValue("script").(string)
			}

			var err error
			cmder, err = cmd.NewJavaScript(script)
			if err != nil {
				return nil, writeflow.NewExecNodeError(fmt.Errorf("parse script error: %v", err), node.Id)
			}
		case writeflow.BuiltInCmd:
			cmdName = node.Data.Source.BuiltinCmd
		}

		nodes[node.Id] = writeflow.Node{
			Id:       node.Id,
			Cmd:      cmdName,
			BuiltCmd: cmder,
			Inputs:   inputs,
		}
	}
	return &writeflow.Flow{
		Nodes:        nodes,
		OutputNodeId: m.Graph.GetOutputNodeId(),
	}, nil
}
