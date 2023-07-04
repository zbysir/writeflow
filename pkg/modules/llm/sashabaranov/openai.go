package sashabaranov

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"github.com/spf13/cast"
	"github.com/zbysir/writeflow/pkg/export"
	"github.com/zbysir/writeflow/pkg/modules/llm/util"
	"io"
)

// Plugin implement PluginLLM
type Plugin struct {
}

func NewPlugin() *Plugin {
	return &Plugin{}
}

func (p *Plugin) NewOpenAICmd() export.CMDer {
	return util.NewFun(func(ctx context.Context, params map[string]interface{}) (rsp map[string]interface{}, err error) {
		key := params["api_key"].(string)
		baseUrl := cast.ToString(params["base_url"])
		config := openai.DefaultConfig(key)
		if baseUrl != "" {
			config.BaseURL = baseUrl
		}
		client := openai.NewClientWithConfig(config)
		return map[string]interface{}{"default": client}, nil
	})
}

func (p *Plugin) CallOpenAICmd() export.CMDer {
	return util.NewFun(func(ctx context.Context, params map[string]interface{}) (rsp map[string]interface{}, err error) {
		//log.Infof("langchain_call")
		openaiClient := params["llm"].(*openai.Client)
		promptI := params["prompt"]
		functionI := params["functions"]
		if promptI == nil {
			return nil, fmt.Errorf("prompt is nil")
		}
		enableSteam := cast.ToBool(params["stream"])
		prompt := promptI.(string)
		var functions []*openai.FunctionDefine
		if functionI != nil {
			function := functionI.(string)
			err = json.Unmarshal([]byte(function), &functions)
			if err != nil {
				return nil, err
			}
		}

		var messages []openai.ChatCompletionMessage
		var chatMemory util.ChatMemory
		if params["chat_memory"] != nil {
			chatMemory = params["chat_memory"].(util.ChatMemory)
		}

		if chatMemory != nil {
			messages = append(messages, coverMessageListToSDK(chatMemory.GetHistory(ctx))...)
		}

		userMsg := openai.ChatCompletionMessage{Content: prompt, Role: openai.ChatMessageRoleUser}
		if chatMemory != nil {
			chatMemory.AppendHistory(ctx, coverMessageToBase(userMsg))
		}
		messages = append(messages, userMsg)

		if enableSteam {
			s, err := openaiClient.CreateChatCompletionStream(ctx, openai.ChatCompletionRequest{
				Model:            "gpt-3.5-turbo-0613",
				Messages:         messages,
				MaxTokens:        2000,
				Temperature:      0,
				TopP:             0,
				N:                0,
				Stream:           true,
				Stop:             nil,
				PresencePenalty:  0,
				FrequencyPenalty: 0,
				LogitBias:        nil,
				User:             "",
				Functions:        functions,
				FunctionCall:     "",
			})
			if err != nil {
				return nil, err
			}

			steam := util.NewSteamResponse()
			go func() {
				defer s.Close()
				var content string
				for {
					recv, err := s.Recv()
					if err != nil {
						if err == io.EOF {
							break
						}
						steam.Close(err)
						break
					}
					if len(recv.Choices) == 0 {
						// 心跳，通常是 30s 一次。
						steam.Close(fmt.Errorf("recv.Choices is empty"))
						break
					}

					c := recv.Choices[0].Delta.Content
					if len(c) != 0 {
						content += c
						steam.Append(c)
					}
				}
				steam.Close(nil)

				if chatMemory != nil {
					if content != "" {
						chatMemory.AppendHistory(ctx, coverMessageToBase(openai.ChatCompletionMessage{
							Role:    openai.ChatMessageRoleAssistant,
							Content: content,
						}))
					}
				}
			}()

			return map[string]interface{}{"default": steam, "function_call": ""}, nil
		} else {
			rsp, err := openaiClient.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
				Model:            "gpt-3.5-turbo-0613",
				Messages:         messages,
				MaxTokens:        2000,
				Temperature:      0,
				TopP:             0,
				N:                0,
				Stream:           false,
				Stop:             nil,
				PresencePenalty:  0,
				FrequencyPenalty: 0,
				LogitBias:        nil,
				User:             "",
				Functions:        functions,
				FunctionCall:     "",
			})
			if err != nil {
				return nil, err
			}

			content := rsp.Choices[0].Message.Content
			if chatMemory != nil {
				chatMemory.AppendHistory(ctx, coverMessageToBase(rsp.Choices[0].Message))
			}

			return map[string]interface{}{"default": content, "function_call": rsp.Choices[0].Message.FunctionCall}, nil
		}
	})
}

func (p *Plugin) SupportStream() bool {
	return true
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
