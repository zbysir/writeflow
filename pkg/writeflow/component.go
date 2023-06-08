package writeflow

import (
	"context"
	"github.com/zbysir/writeflow/internal/cmd"
	"github.com/zbysir/writeflow/internal/model"
	"github.com/zbysir/writeflow/pkg/schema"
)

type nothingCMD struct {
}

func (n nothingCMD) Exec(ctx context.Context, params map[string]interface{}) (rsp map[string]interface{}, err error) {
	return nil, nil
}

var _nothingCMD = nothingCMD{}

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
	case model.GoScriptCmd:
		key := m.Data.Source.GoScript.InputKey
		if key == "" {
			key = "script"
		}
		script := m.Data.GetInputValue(key)
		cmder, err = cmd.NewGoScript(nil, "", script)
		if err != nil {
			return
		}
	case model.BuiltInCmd:
		cmder = builtinCmd[m.Data.Source.BuiltinCmd]
	case model.NothingCmd:
		cmder = _nothingCMD
	}

	sc = cmd.Schema{
		Key:     m.Key,
		Inputs:  input,
		Outputs: output,
	}

	return
}
