package util

import (
	"context"
	"github.com/sashabaranov/go-openai"
)

type Message struct {
	// Role: Either of "system", "user", "assistant".
	Role string `json:"role"`

	// Content: A content of the message.
	Content string `json:"content"`

	// FunctionCall requested by ChatGPT.
	// Only appears in a response from ChatGPT in which ChatGPT wants to call a function.
	FunctionCall *FunctionCall `json:"function_call,omitempty"`

	// Name of the function called, to tell this message is a result of function_call.
	// Only appears in a request from us when the previous message is "function_call" requested by ChatGPT.
	Name string `json:"name,omitempty"`
}

type FunctionCall struct {
	Name string `json:"name,omitempty"`
	// call function with arguments in JSON format
	Arguments string `json:"arguments,omitempty"`
}

type Messages = []Message

type ChatMemory interface {
	GetHistory(ctx context.Context) Messages
	AppendHistory(ctx context.Context, message Message)
}

var history map[string][]openai.ChatCompletionMessage

func init() {
	history = map[string][]openai.ChatCompletionMessage{}
}

type MemoryChatMemory struct {
	sessionId string
	maxSize   int
}

func NewMemoryChatMemory(sessionId string) *MemoryChatMemory {
	return &MemoryChatMemory{
		sessionId: sessionId,
	}
}

func (m *MemoryChatMemory) GetHistory(ctx context.Context) []openai.ChatCompletionMessage {
	if m.sessionId == "" {
		return nil
	}
	return history[m.sessionId]
}

func (m *MemoryChatMemory) AppendHistory(ctx context.Context, message openai.ChatCompletionMessage) {
	if m.sessionId == "" {
		return
	}

	history[m.sessionId] = append(history[m.sessionId], message)
	if m.maxSize != 0 {
		if len(history[m.sessionId]) > m.maxSize {
			history[m.sessionId] = history[m.sessionId][1:]
		}
	}
}
