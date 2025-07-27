package model

import "Voice-Assistant/pkg/llm"

type History struct {
	ID          string    `gorm:"primaryKey"` // primary key
	AssistantID string    `gorm:"index"`      // to bind assistant
	Input       []Message `gorm:"serializer:json"`
	Output      Output    `gorm:"serializer:json"`
	Usage       Usage     `gorm:"embedded"`
	CreatedAt   string
}

type Output struct {
	Content      string `json:"content"`
	FinishReason string `json:"finish_reason"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Usage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
	TotalTokens  int `json:"total_tokens"`
}

func HistoryToLLMMessages(histories []History) []llm.LLMMessage {
	var messages []llm.LLMMessage
	for _, h := range histories {
		messages = append(messages, llm.LLMMessage{
			Role:    "user",
			Content: h.Input[1].Content,
		})
		messages = append(messages, llm.LLMMessage{
			Role:    "assistant",
			Content: h.Output.Content,
		})
	}

	return messages
}
