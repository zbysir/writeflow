package writeflow

import (
	"github.com/zbysir/writeflow/internal/cmd"
	"github.com/zbysir/writeflow/internal/model"
	"github.com/zbysir/writeflow/pkg/schema"
)

func ComponentFromModel(m *model.Component, builtinCmd map[string]schema.CMDer) (cmder schema.CMDer, sc cmd.Schema, err error) {
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

	switch m.Data.Source.CmdType {
	case "go_script":
		cmder, err = cmd.NewGoScript(nil, "", m.Data.Source.GoScript.Script)
		if err != nil {
			return
		}
	case "builtin":
		cmder = builtinCmd[m.Data.Source.BuiltinCmd]
	}

	sc = cmd.Schema{
		Key:     m.Key,
		Inputs:  input,
		Outputs: output,
	}

	return
}
