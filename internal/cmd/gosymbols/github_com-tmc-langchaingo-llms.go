// Code generated by 'yaegi extract github.com/tmc/langchaingo/llms'. DO NOT EDIT.

package gosymbols

import (
	"context"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/schema"
	"reflect"
)

func init() {
	Symbols["github.com/tmc/langchaingo/llms/llms"] = map[string]reflect.Value{
		// function, constant and variable definitions
		"WithMaxTokens":   reflect.ValueOf(llms.WithMaxTokens),
		"WithModel":       reflect.ValueOf(llms.WithModel),
		"WithOptions":     reflect.ValueOf(llms.WithOptions),
		"WithStopWords":   reflect.ValueOf(llms.WithStopWords),
		"WithTemperature": reflect.ValueOf(llms.WithTemperature),

		// type definitions
		"CallOption":     reflect.ValueOf((*llms.CallOption)(nil)),
		"CallOptions":    reflect.ValueOf((*llms.CallOptions)(nil)),
		"ChatGeneration": reflect.ValueOf((*llms.ChatGeneration)(nil)),
		"ChatLLM":        reflect.ValueOf((*llms.ChatLLM)(nil)),
		"Generation":     reflect.ValueOf((*llms.Generation)(nil)),
		"LLM":            reflect.ValueOf((*llms.LLM)(nil)),

		// interface wrapper definitions
		"_ChatLLM": reflect.ValueOf((*_github_com_tmc_langchaingo_llms_ChatLLM)(nil)),
		"_LLM":     reflect.ValueOf((*_github_com_tmc_langchaingo_llms_LLM)(nil)),
	}
}

// _github_com_tmc_langchaingo_llms_ChatLLM is an interface wrapper for ChatLLM type
type _github_com_tmc_langchaingo_llms_ChatLLM struct {
	IValue interface{}
	WChat  func(ctx context.Context, messages []schema.ChatMessage, options ...llms.CallOption) (*llms.ChatGeneration, error)
}

func (W _github_com_tmc_langchaingo_llms_ChatLLM) Chat(ctx context.Context, messages []schema.ChatMessage, options ...llms.CallOption) (*llms.ChatGeneration, error) {
	return W.WChat(ctx, messages, options...)
}

// _github_com_tmc_langchaingo_llms_LLM is an interface wrapper for LLM type
type _github_com_tmc_langchaingo_llms_LLM struct {
	IValue    interface{}
	WCall     func(ctx context.Context, prompt string, options ...llms.CallOption) (string, error)
	WGenerate func(ctx context.Context, prompts []string, options ...llms.CallOption) ([]*llms.Generation, error)
}

func (W _github_com_tmc_langchaingo_llms_LLM) Call(ctx context.Context, prompt string, options ...llms.CallOption) (string, error) {
	return W.WCall(ctx, prompt, options...)
}
func (W _github_com_tmc_langchaingo_llms_LLM) Generate(ctx context.Context, prompts []string, options ...llms.CallOption) ([]*llms.Generation, error) {
	return W.WGenerate(ctx, prompts, options...)
}
