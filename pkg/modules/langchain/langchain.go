package langchain

import (
	"context"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/zbysir/writeflow/internal/cmd"
	"github.com/zbysir/writeflow/internal/model"
	"github.com/zbysir/writeflow/pkg/modules"
	"github.com/zbysir/writeflow/pkg/schema"
	"reflect"
)

type LangChain struct {
}

func NewLangChain() *LangChain {
	return &LangChain{}
}

func (l *LangChain) Info() modules.ModuleInfo {
	return modules.ModuleInfo{
		NameSpace: "langchain",
	}
}

func (l *LangChain) Categories() []model.Category {
	return []model.Category{
		{
			Key: "llm",
			Name: map[string]string{
				"zh-CN": "LLM",
			},
			Desc: nil,
		},
	}
}

func (l *LangChain) Components() []model.Component {
	return []model.Component{
		{
			Id:       0,
			Key:      "new_openai",
			Category: "llm",
			Data: model.ComponentData{
				Name: map[string]string{
					"zh-CN": "OpenAI",
				},
				Icon:        "",
				Description: map[string]string{},
				Source: model.ComponentSource{
					CmdType:    model.BuiltInCmd,
					BuiltinCmd: "new_openai",
					GoPackage:  model.ComponentGoPackage{},
					GoScript:   model.ComponentGoScript{},
				},
				InputAnchors: nil,
				InputParams: []model.NodeInputParam{
					{
						Id: "",
						Name: map[string]string{
							"zh-CN": "ApiKey",
						},
						Key:      "api_key",
						Type:     "string",
						Optional: false,
					},
				},
				OutputAnchors: []model.NodeAnchor{
					{
						Id: "",
						Name: map[string]string{
							"zh-CN": "Default",
						},
						Key:      "default",
						Type:     "langchain/llm",
						Optional: false,
					},
				},
			},
		},
		{
			Key:      "langchain_call",
			Category: "llm",
			Data: model.ComponentData{
				Name: map[string]string{
					"zh-CN": "LangChain",
				},
				Icon:        "",
				Description: map[string]string{},
				Source: model.ComponentSource{
					CmdType:    model.BuiltInCmd,
					BuiltinCmd: "langchain_call",
				},
				InputAnchors: []model.NodeAnchor{
					{
						Name: map[string]string{
							"zh-CN": "LLM",
						},
						Key:  "llm",
						Type: "langchain/llm",
					},
					{
						Name: map[string]string{
							"zh-CN": "Prompt",
						},
						Key:  "prompt",
						Type: "string",
					},
				},
				InputParams: []model.NodeInputParam{},
				OutputAnchors: []model.NodeAnchor{
					{
						Key:  "default",
						Type: "string",
					},
				},
			},
		},
	}
}

func (l *LangChain) Cmd() map[string]schema.CMDer {
	return map[string]schema.CMDer{
		"new_openai": cmd.NewFun(func(ctx context.Context, params map[string]interface{}) (rsp map[string]interface{}, err error) {
			key := params["api_key"].(string)
			ll, err := openai.New(openai.WithToken(key))
			if err != nil {
				return nil, err
			}

			return map[string]interface{}{"default": ll}, nil
		}),
		"langchain_call": cmd.NewFun(func(ctx context.Context, params map[string]interface{}) (rsp map[string]interface{}, err error) {
			llm := params["llm"].(llms.LLM)
			prompt := params["prompt"].(string)
			s, err := llm.Call(ctx, prompt)
			if err != nil {
				return nil, err
			}

			return map[string]interface{}{"default": s}, nil
		}),
	}
}

func (l *LangChain) GoSymbols() map[string]map[string]reflect.Value {
	return nil
}

var _ modules.Module = (*LangChain)(nil)
