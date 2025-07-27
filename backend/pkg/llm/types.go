package llm

import (
	"Voice-Assistant/tools"
	"net/http"
)

type Input struct {
	LLMMessages []LLMMessage `json:"messages"`
}

// LLMClient 提供与LLM API交互的客户端
type LLMClient struct {
	apiKey     string
	httpClient *http.Client
	baseURL    string
}

// LLMConfig 表示LLM的配置
type LLMConfig struct {
	Model       string       `json:"model"`
	Temperature float64      `json:"temperature"`
	MaxTokens   int          `json:"max_tokens"`
	Prompt      string       `json:"prompt"`
	LLMMessages []LLMMessage `json:"messages"`
}

// LLMRequest 表示发送给LLM的请求
type LLMRequest struct {
	Model       string     `json:"model"`
	Input       Input      `json:"input"`
	Temperature float64    `json:"temperature"` // 可选
	MaxTokens   int        `json:"max_tokens"`
	Stream      bool       `json:"stream,omitempty"`
	Parameters  Parameters `json:"parameters"`
}

type FunctionRequest struct {
	Model    string       `json:"model"`
	Messages []LLMMessage `json:"messages"`
	Tools    []tools.Tool `json:"tools"`
	Stream   bool         `json:"stream"`
}

// LLMMessage 表示对话中的一条消息
type LLMMessage struct {
	Role       string     `json:"role"`
	Content    string     `json:"content"`
	ToolCalls  []ToolCall `json:"tool_calls,omitempty"`
	ToolCallId string     `json:"tool_call_id,omitempty"`
}

type Parameters struct {
	EnableThinking    bool   `json:"enable_thinking,omitempty"`
	ResultFormat      string `json:"result_format"`
	IncrementalOutput bool   `json:"incremental_output,omitempty"`
	EnableSearch      bool   `json:"enable_search,omitempty"`
}

// DashScopeResponse 表示 DashScope 平台返回的完整响应
type DashScopeResponse struct {
	Output    Output `json:"output"`
	Usage     Usage  `json:"usage"`
	RequestId string `json:"request_id,omitempty"`
}

// Output 包含模型生成的文本和状态
type Output struct {
	Choices      *[]Choice `json:"choices"`
	Text         string    `json:"text,omitempty"`
	FinishReason string    `json:"finish_reason,omitempty"`
	GmtCreate    string    `json:"gmt_create,omitempty"`
	//LogProbs     *LogProbs `json:"logprobs,omitempty"` // 可选字段
}

// LogProbs 是可选字段，包含 token 概率信息
type LogProbs struct {
	// 如果需要可以扩展具体内容
}

// Usage 包含请求消耗的 token 数量
type Usage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
	TotalTokens  int `json:"total_tokens"`
}

type FunctionResp struct {
	Choices []Choice `json:"choices"`
	Object  string   `json:"object"`
	Usage   Usage    `json:"usage"`
	Created uint     `json:"created"`
	Model   string   `json:"model"`
	Id      string   `json:"id"`
}

// Choice 表示LLM返回的选择
type Choice struct {
	Message      ChoiceMessage `json:"message"`
	FinishReason string        `json:"finish_reason,omitempty"`
	Index        int           `json:"index,omitempty"`
	LogProbs     LogProbs      `json:"logprobs,omitempty"`
}

type ChoiceMessage struct {
	Content          string      `json:"content"`
	Role             string      `json:"role"`
	ReasoningContent string      `json:"reasoning_content,omitempty"`
	ToolCalls        *[]ToolCall `json:"tool_calls,omitempty"`
}

type ToolCall struct {
	Index    int          `json:"index,omitempty"`
	Id       string       `json:"id"`
	Type     string       `json:"type"`
	Function FunctionCall `json:"function"`
}

type FunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}
