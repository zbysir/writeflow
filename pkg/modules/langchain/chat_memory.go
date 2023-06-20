package langchain

import (
	"context"
	"github.com/sashabaranov/go-openai"
)

type ChatMemory interface {
	GetHistory(ctx context.Context) []openai.ChatCompletionMessage
	AppendHistory(ctx context.Context, message openai.ChatCompletionMessage)
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
