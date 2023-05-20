package writeflow

import (
	"github.com/zbysir/writeflow/internal/cmd"
	"github.com/zbysir/writeflow/internal/model"
	"github.com/zbysir/writeflow/pkg/schema"
)

type Component struct {
	Cmder  schema.CMDer
	Schema cmd.Schema
}

func NewComponent(cmder schema.CMDer, schema cmd.Schema) *Component {
	return &Component{Cmder: cmder, Schema: schema}
}

func ComponentFromModel(m *model.Component) (c Component, err error) {
	var input []cmd.SchemaParams
	for _, a := range m.Data.InputAnchors {
		input = append(input, cmd.SchemaParams{
			Key:         a.Key,
			Type:        "anchor",
			NameLocales: a.Name,
			DescLocales: nil,
		})
	}
	for _, a := range m.Data.InputParams {
		input = append(input, cmd.SchemaParams{
			Key:         a.Key,
			Type:        "literal",
			NameLocales: a.Name,
			DescLocales: nil,
		})
	}

	var output []cmd.SchemaParams
	for _, a := range m.Data.OutputAnchors {
		output = append(output, cmd.SchemaParams{
			Key:         a.Key,
			Type:        "anchor",
			NameLocales: a.Name,
			DescLocales: nil,
		})
	}

	var cmder schema.CMDer
	switch m.Data.Source.CmdType {
	case "go_script":
		cmder, err = cmd.NewGoScript(nil, "", m.Data.Source.GoScript)
		if err != nil {
			return
		}
	}

	c = Component{
		Schema: cmd.Schema{
			Key:     m.Key,
			Inputs:  input,
			Outputs: output,
		},
		Cmder: cmder,
	}

	return
}
