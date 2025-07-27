package service

import (
	"Voice-Assistant/pkg/llm"
	"log"
)

type LLMService interface {
	GenerateLLMResponse(userInput, systemPrompt string) (*llm.DashScopeResponse, error)
	GenerateStreamResponse(userInput, systemPrompt string) (<-chan *llm.DashScopeResponse, error)
	GenerateResponseWithContext(systemPrompt, userInput string, history []llm.LLMMessage) (*llm.DashScopeResponse, error)
	GenerateStreamResponseWithContext(systemPrompt, userInput string, history []llm.LLMMessage) (<-chan *llm.DashScopeResponse, error)
	GenerateResponseWithFunctionCalling(systemPrompt, userInput string, history []llm.LLMMessage) (string, error)
}

type llmServiceImpl struct {
	client llm.LLMClientInterface
}

func NewLLMService(client llm.LLMClientInterface) LLMService {
	return &llmServiceImpl{client: client}
}

func (l *llmServiceImpl) GenerateLLMResponse(userInput, systemPrompt string) (*llm.DashScopeResponse, error) {
	cfg := &llm.LLMConfig{
		Model:       "qwen-plus",
		Temperature: 0.7,
		MaxTokens:   1000,
		Prompt:      systemPrompt,
	}
	return l.client.GenerateResponse(userInput, cfg)
}

func (l *llmServiceImpl) GenerateStreamResponse(userInput, systemPrompt string) (<-chan *llm.DashScopeResponse, error) {
	cfg := &llm.LLMConfig{
		Model:       "qwen-plus",
		Temperature: 0.7,
		MaxTokens:   1000,
		Prompt:      systemPrompt,
	}
	return l.client.GenerateStreamResponse(userInput, cfg)
}

func (l *llmServiceImpl) GenerateResponseWithContext(systemPrompt, userInput string, history []llm.LLMMessage) (*llm.DashScopeResponse, error) {
	cfg := &llm.LLMConfig{
		Model:       "qwen-plus",
		Temperature: 0.7,
		MaxTokens:   1000,
		LLMMessages: BuildLLMContext(systemPrompt, userInput, history),
	}
	return l.client.GenerateResponseWithContext(cfg)
}

func (l *llmServiceImpl) GenerateStreamResponseWithContext(systemPrompt, userInput string, history []llm.LLMMessage) (<-chan *llm.DashScopeResponse, error) {
	cfg := BuildLLMConfig(systemPrompt, userInput, history)
	return l.client.GenerateStreamResponseWithHistory(cfg)
}

func (l *llmServiceImpl) GenerateResponseWithFunctionCalling(systemPrompt, userInput string, history []llm.LLMMessage) (string, error) {
	full := make([]llm.LLMMessage, 0, len(history)+2)
	if systemPrompt != "" {
		full = append(full, llm.LLMMessage{Role: "system", Content: systemPrompt})
	}
	if userInput != "" {
		full = append(full, llm.LLMMessage{Role: "user", Content: userInput})
	}
	cfg := &llm.LLMConfig{
		Model:       "qwen-plus",
		Temperature: 0.7,
		MaxTokens:   1000,
		LLMMessages: full,
	}
	log.Println("运行到1")
	resp, err := l.client.GenerateResponseWithFunction(cfg)
	if err != nil || resp == nil {
		return "第一次调用失败", err
	}
	var choices = resp.Choices
	secondFull := make([]llm.LLMMessage, 0, len(history)+len(choices))
	secondFull = append(secondFull, full...)
	for _, choice := range choices {
		var toolCalls = choice.Message.ToolCalls
		if toolCalls != nil && len(*toolCalls) > 0 {
			for _, toolCall := range *toolCalls {
				var functionCall = toolCall.Function
				result, err := CallToolByName(functionCall.Name, functionCall)
				if err != nil {
					log.Printf("获取工具函数结果出错: %v", err)
				}
				secondFull = append(secondFull, llm.LLMMessage{
					Role:    "assistant",
					Content: "",
					ToolCalls: []llm.ToolCall{
						{
							Id:       toolCall.Id,
							Function: functionCall,
							Type:     "function",
						},
					}})
				secondFull = append(secondFull, llm.LLMMessage{
					Role:       "tool",
					Content:    result,
					ToolCallId: toolCall.Id,
				})
			}
		}
	}
	secondCfg := &llm.LLMConfig{
		Model:       "qwen-plus",
		Temperature: 0.7,
		MaxTokens:   1000,
		LLMMessages: secondFull,
	}
	secondResp, err := l.client.GenerateResponseWithFunction(secondCfg)
	if err != nil {
		log.Printf("第二次调用失败: %+v", err)
	}
	text := secondResp.Choices[0].Message.Content
	return text, nil
}

func BuildLLMContext(systemPrompt, userInput string, history []llm.LLMMessage) []llm.LLMMessage {
	full := make([]llm.LLMMessage, 0, len(history)+2)
	if systemPrompt != "" {
		full = append(full, llm.LLMMessage{Role: "system", Content: systemPrompt})
	}
	full = append(full, history...)
	if userInput != "" {
		full = append(full, llm.LLMMessage{Role: "user", Content: userInput})
	}
	return full
}

func BuildLLMConfig(systemPrompt, userInput string, history []llm.LLMMessage) *llm.LLMConfig {
	return &llm.LLMConfig{
		Model:       "qwen-plus",
		Temperature: 0.7,
		MaxTokens:   1000,
		Prompt:      systemPrompt,
		LLMMessages: BuildLLMContext(systemPrompt, userInput, history),
	}
}
