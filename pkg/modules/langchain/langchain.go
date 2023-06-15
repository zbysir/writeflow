package langchain

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/zbysir/writeflow/internal/cmd"
	"github.com/zbysir/writeflow/internal/model"
	llms2 "github.com/zbysir/writeflow/internal/pkg/langchaingo/llms"
	openai2 "github.com/zbysir/writeflow/internal/pkg/langchaingo/llms/openai"
	schema2 "github.com/zbysir/writeflow/internal/pkg/langchaingo/schema"
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
			Type:     "new_openai",
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
					Script:     model.ComponentScript{},
				},
				InputParams: []model.NodeInputParam{
					{
						Name: map[string]string{
							"zh-CN": "ApiKey",
						},
						Key:      "api_key",
						Type:     "string",
						Optional: false,
					},
				},
				OutputAnchors: []model.NodeOutputAnchor{
					{
						Name: map[string]string{
							"zh-CN": "Default",
						},
						Key:  "default",
						Type: "langchain/llm",
					},
				},
			},
		},
		{
			Type:     "langchain_call",
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
				InputParams: []model.NodeInputParam{
					{
						InputType: model.NodeInputTypeAnchor,
						Name: map[string]string{
							"zh-CN": "LLM",
						},
						Key:  "llm",
						Type: "langchain/llm",
					},
					{
						InputType: model.NodeInputTypeAnchor,
						Name: map[string]string{
							"zh-CN": "Functions",
						},
						Key:      "functions",
						Type:     "string",
						Optional: true,
					},
					{
						InputType: model.NodeInputTypeAnchor,
						Name: map[string]string{
							"zh-CN": "Prompt",
						},
						Key:  "prompt",
						Type: "string",
					},
				},
				OutputAnchors: []model.NodeOutputAnchor{
					{
						Key:  "default",
						Type: "string",
					},
					{
						Key:  "function_call",
						Type: "any",
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
			ll, err := openai2.New(openai2.WithToken(key), openai2.WithModel("gpt-3.5-turbo-0613"))
			if err != nil {
				return nil, err
			}

			return map[string]interface{}{"default": ll}, nil
		}),
		"langchain_call": cmd.NewFun(func(ctx context.Context, params map[string]interface{}) (rsp map[string]interface{}, err error) {
			llm := params["llm"].(llms2.ChatLLM)
			promptI := params["prompt"]
			functionI := params["functions"]
			if promptI == nil {
				return nil, fmt.Errorf("prompt is nil")
			}
			prompt := promptI.(string)
			var functions []llms2.Function
			if functionI != nil {
				function := functionI.(string)
				err = json.Unmarshal([]byte(function), &functions)
				if err != nil {
					return nil, err
				}
			}

			s, err := llm.Chat(ctx, []schema2.ChatMessage{
				schema2.HumanChatMessage{Text: prompt},
			}, llms2.WithFunctions(functions), llms2.WithModel("gpt-3.5-turbo-0613"))
			if err != nil {
				return nil, err
			}

			return map[string]interface{}{"default": s.Message.Text, "function_call": s.Message.FunctionCall}, nil
		}),
	}
}

func (l *LangChain) GoSymbols() map[string]map[string]reflect.Value {
	return nil
}

var _ modules.Module = (*LangChain)(nil)
