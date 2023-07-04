package llm

import (
	"context"
	"github.com/sashabaranov/go-openai"
	"github.com/spf13/cast"
	"github.com/zbysir/writeflow/pkg/export"
	"github.com/zbysir/writeflow/pkg/modules/llm/sashabaranov"
	"github.com/zbysir/writeflow/pkg/modules/llm/util"
	"reflect"
)

type PluginLLM interface {
	NewOpenAICmd() export.CMDer
	CallOpenAICmd() export.CMDer
	SupportStream() bool
}

type LangChain struct {
	pluginLLM          PluginLLM
	libraryVectorStore VectorStoreFactory
}

func NewLangChain(libraryVectorStore VectorStoreFactory) export.Plugin {
	return &LangChain{pluginLLM: sashabaranov.NewPlugin(), libraryVectorStore: libraryVectorStore}
}

func (l *LangChain) Info() export.PluginInfo {
	return export.PluginInfo{
		NameSpace: "langchain",
	}
}

func (l *LangChain) Categories() []export.Category {
	return []export.Category{
		{
			Key: "llm",
			Name: map[string]string{
				"zh-CN": "LLM",
			},
			Desc: nil,
		},
	}
}

func (l *LangChain) Components() []export.Component {
	var langchainCallInputParams []export.NodeInputParam
	if l.pluginLLM.SupportStream() {
		langchainCallInputParams = append(langchainCallInputParams, export.NodeInputParam{
			Name: map[string]string{
				"zh-CN": "流式返回",
			},
			Value: true,
			Key:   "stream",
			Type:  "bool",
		})
	}

	return []export.Component{
		{
			Id:       0,
			Type:     "new_openai",
			Category: "llm",
			Data: export.ComponentData{
				Name: map[string]string{
					"zh-CN": "OpenAI",
				},
				Source: export.ComponentSource{
					CmdType:    "builtin",
					BuiltinCmd: "new_openai",
				},
				InputParams: []export.NodeInputParam{
					{
						Name: map[string]string{"zh-CN": "ApiKey"},
						Key:  "api_key",
						Type: "string",
					},
					{
						Name: map[string]string{"zh-CN": "BaseURL"},
						Key:  "base_url",
						Type: "string",
					},
				},
				OutputAnchors: []export.NodeOutputAnchor{
					{
						Name: map[string]string{"zh-CN": "Default"},
						Key:  "default",
						Type: "llm.llm",
					},
				},
			},
		},
		{
			Id:       0,
			Type:     "chat_memory",
			Category: "llm",
			Data: export.ComponentData{
				Name: map[string]string{"zh-CN": "ChatMemory"},
				Source: export.ComponentSource{
					CmdType:    "builtin",
					BuiltinCmd: "chat_memory",
				},
				InputParams: []export.NodeInputParam{
					{
						Name:     map[string]string{"zh-CN": "SessionID"},
						Key:      "session_id",
						Type:     "string",
						Optional: true,
					},
				},
				OutputAnchors: []export.NodeOutputAnchor{
					{
						Name: map[string]string{"zh-CN": "Default"},
						Key:  "default",
						Type: "llm.chat_memory",
					},
				},
			},
		},
		{
			Type:     "call_openai",
			Category: "llm",
			Data: export.ComponentData{
				Name:        map[string]string{"zh-CN": "LangChain"},
				Icon:        "",
				Description: map[string]string{},
				Source: export.ComponentSource{
					CmdType:    "builtin",
					BuiltinCmd: "call_openai",
				},
				InputParams: append([]export.NodeInputParam{
					{
						InputType: "anchor",
						Name: map[string]string{
							"zh-CN": "LLM",
						},
						Key:  "llm",
						Type: "llm.llm",
					},
					{
						InputType: "anchor",
						Name: map[string]string{
							"zh-CN": "ChatMemory",
						},
						Key:      "chat_memory",
						Type:     "llm.chat_memory",
						Optional: true,
					},
					{
						InputType: "anchor",
						Name:      map[string]string{"zh-CN": "Functions"},
						Key:       "functions",
						Type:      "string",
						Optional:  true,
					},
					{
						InputType: "anchor",
						Name: map[string]string{
							"zh-CN": "Prompt",
						},
						Key:  "prompt",
						Type: "string",
					},
				}, langchainCallInputParams...),
				OutputAnchors: []export.NodeOutputAnchor{
					{
						Name: map[string]string{
							"zh-CN": "Default",
						},
						Key:  "default",
						Type: "string",
					},
					{
						Name: map[string]string{
							"zh-CN": "FunctionCall",
						},
						Key:  "function_call",
						Type: "any",
					},
				},
			},
		},
		{
			Type:     "embedding_openai",
			Category: "llm",
			Data: export.ComponentData{
				Name:        map[string]string{"zh-CN": "OpenAIEmbedding"},
				Icon:        "",
				Description: map[string]string{},
				Source: export.ComponentSource{
					CmdType:    "builtin",
					BuiltinCmd: "openai_create_embedding",
				},
				InputParams: []export.NodeInputParam{
					{
						InputType: "anchor",
						Name: map[string]string{
							"zh-CN": "Query",
						},
						Key:  "query",
						Type: "string",
					},
					{
						InputType: "anchor",
						Name: map[string]string{
							"zh-CN": "OpenAI",
						},
						Key:  "llm",
						Type: "llm.llm",
					},
				},
				OutputAnchors: []export.NodeOutputAnchor{
					{
						Name: map[string]string{
							"zh-CN": "Default",
						},
						Key:  "default",
						Type: "llm.vector",
					},
				},
			},
		},
		{
			Type:     "similarity_search",
			Category: "llm",
			Data: export.ComponentData{
				Name:        map[string]string{"zh-CN": "相似度搜索", "en": "SimilaritySearch"},
				Icon:        "",
				Description: map[string]string{},
				Source: export.ComponentSource{
					CmdType:    "builtin",
					BuiltinCmd: "similarity_search",
				},
				InputParams: []export.NodeInputParam{
					{
						InputType: "anchor",
						Name: map[string]string{
							"zh-CN": "Embedding",
						},
						Key:  "embedding",
						Type: "llm.vector",
					},
					{
						InputType: "anchor",
						Name: map[string]string{
							"zh-CN": "向量数据库",
						},
						Key:  "vector_store",
						Type: "llm.vector_store",
					},
					{
						InputType: "input",
						Name: map[string]string{
							"zh-CN": "数量",
						},
						Key:  "number",
						Type: "int",
					},
				},
				OutputAnchors: []export.NodeOutputAnchor{
					{
						Name: map[string]string{
							"zh-CN": "Default",
						},
						Key:  "default",
						Type: "[]llm.fragment",
					},
				},
			},
		},
		{
			Type:     "vector_store_library",
			Category: "llm",
			Data: export.ComponentData{
				Name:        map[string]string{"zh-CN": "资料库", "en": "Library"},
				Icon:        "",
				Description: map[string]string{},
				Source: export.ComponentSource{
					CmdType:    "builtin",
					BuiltinCmd: "vector_store_library",
				},
				InputParams: []export.NodeInputParam{
					{
						InputType: "input",
						Name: map[string]string{
							"zh-CN": "书籍 ID",
						},
						Key:  "book_id",
						Type: "int",
					},
				},
				OutputAnchors: []export.NodeOutputAnchor{
					{
						Name: map[string]string{
							"zh-CN": "Default",
						},
						Key:  "default",
						Type: "llm.vector_store",
					},
				},
			},
		},
	}
}

func coverMessageToBase(a openai.ChatCompletionMessage) util.Message {
	var fc *util.FunctionCall
	if a.FunctionCall != nil {
		fc = &util.FunctionCall{
			Name:      a.FunctionCall.Name,
			Arguments: a.FunctionCall.Arguments,
		}
	}
	return util.Message{
		Role:         a.Role,
		Content:      a.Content,
		FunctionCall: fc,
		Name:         a.Name,
	}
}

func coverMessageToSDK(a util.Message) openai.ChatCompletionMessage {
	var fc *openai.FunctionCall
	if a.FunctionCall != nil {
		fc = &openai.FunctionCall{
			Name:      a.FunctionCall.Name,
			Arguments: a.FunctionCall.Arguments,
		}
	}
	return openai.ChatCompletionMessage{
		Role:         a.Role,
		Content:      a.Content,
		FunctionCall: fc,
		Name:         a.Name,
	}
}

func coverMessageListToSDK(as []util.Message) []openai.ChatCompletionMessage {
	var bs []openai.ChatCompletionMessage
	for _, a := range as {
		bs = append(bs, coverMessageToSDK(a))
	}
	return bs
}

type Vector = []float32

func (l *LangChain) Cmd() map[string]export.CMDer {
	return map[string]export.CMDer{
		"new_openai":     l.pluginLLM.NewOpenAICmd(),
		"langchain_call": l.pluginLLM.NewOpenAICmd(), // 废弃
		"call_openai":    l.pluginLLM.CallOpenAICmd(),
		"similarity_search": util.NewFun(func(ctx context.Context, params map[string]interface{}) (rsp map[string]interface{}, err error) {
			embedding := params["embedding"].(Vector)
			vs := params["vector_store"].(VectorStore)
			number := cast.ToInt(params["number"])
			if number == 0 {
				number = 10
			}
			fs, err := vs.SimilaritySearch(ctx, SimilaritySearchParams{
				Vector: embedding,
				Number: number,
			})
			if err != nil {
				return map[string]interface{}{}, err
			}
			return map[string]interface{}{"default": fs}, nil
		}),
		"vector_store_library": util.NewFun(func(ctx context.Context, params map[string]interface{}) (rsp map[string]interface{}, err error) {
			bookId := cast.ToInt(params["book_id"])
			vs, err := l.libraryVectorStore.NewVectorStore(ctx, map[string]interface{}{
				"book_id": bookId,
			})
			if err != nil {
				return nil, err
			}

			return map[string]interface{}{"default": vs}, nil
		}),
		"openai_create_embedding": util.NewFun(func(ctx context.Context, params map[string]interface{}) (rsp map[string]interface{}, err error) {
			openaiClient := params["llm"].(*openai.Client)
			query := params["query"].(string)
			eb := NewOpenAIEmbedding(openaiClient)
			rr, err := eb.Embedding([]string{query})
			if err != nil {
				return map[string]interface{}{}, nil
			}
			return map[string]interface{}{"default": rr[0]}, nil
		}),
		// chat_memory 存储对话记录
		"chat_memory": util.NewFun(func(ctx context.Context, params map[string]interface{}) (rsp map[string]interface{}, err error) {
			idi := params["session_id"]
			if idi == nil {
				return map[string]interface{}{"default": util.NewMemoryChatMemory("")}, nil
			}
			id := idi.(string)

			memory := util.NewMemoryChatMemory(id)
			return map[string]interface{}{"default": memory}, nil
		}),
	}
}

func (l *LangChain) GoSymbols() map[string]map[string]reflect.Value {
	return nil
}

var _ export.Plugin = (*LangChain)(nil)
