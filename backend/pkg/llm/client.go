package llm

import (
	"Voice-Assistant/tools"
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type LLMClientInterface interface {
	GenerateResponse(userInput string, config *LLMConfig) (*DashScopeResponse, error)
	GenerateStreamResponse(userInput string, config *LLMConfig) (<-chan *DashScopeResponse, error)
	GenerateResponseWithContext(config *LLMConfig) (*DashScopeResponse, error)
	GenerateStreamResponseWithHistory(config *LLMConfig) (<-chan *DashScopeResponse, error)
	GenerateResponseWithFunction(cfg *LLMConfig) (*FunctionResp, error)
}

func NewLLMClient() *LLMClient {
	_ = godotenv.Load(".env.development")
	return &LLMClient{
		apiKey:     os.Getenv("LLM_API_KEY"),
		httpClient: &http.Client{Timeout: 30 * time.Second},
		baseURL:    "https://dashscope.aliyuncs.com/api/v1/services/aigc/text-generation/generation",
	}
}

func (c *LLMClient) GenerateResponse(userInput string, config *LLMConfig) (*DashScopeResponse, error) {
	req := &LLMRequest{
		Model: config.Model,
		Input: Input{
			LLMMessages: []LLMMessage{
				{Role: "system", Content: config.Prompt},
				{Role: "user", Content: userInput},
			},
		},
		Temperature: config.Temperature,
		MaxTokens:   config.MaxTokens,
		Stream:      false,
	}
	return c.sendRequest(req)
}

func (c *LLMClient) GenerateStreamResponse(userInput string, config *LLMConfig) (<-chan *DashScopeResponse, error) {
	req := &LLMRequest{
		Model: config.Model,
		Input: Input{
			LLMMessages: []LLMMessage{
				{Role: "system", Content: config.Prompt},
				{Role: "user", Content: userInput},
			},
		},
		Temperature: config.Temperature,
		MaxTokens:   config.MaxTokens,
		Stream:      true,
		Parameters: Parameters{
			ResultFormat:      "text",
			IncrementalOutput: true,
		},
	}
	return c.sendStreamRequest(req)
}

func (c *LLMClient) GenerateResponseWithContext(config *LLMConfig) (*DashScopeResponse, error) {
	req := &LLMRequest{
		Model:       config.Model,
		Input:       Input{LLMMessages: config.LLMMessages},
		Temperature: config.Temperature,
		MaxTokens:   config.MaxTokens,
		Stream:      false,
	}
	return c.sendRequest(req)
}

func (c *LLMClient) GenerateStreamResponseWithHistory(config *LLMConfig) (<-chan *DashScopeResponse, error) {
	req := &LLMRequest{
		Model:       config.Model,
		Input:       Input{LLMMessages: config.LLMMessages},
		Temperature: config.Temperature,
		MaxTokens:   config.MaxTokens,
		Stream:      true,
		Parameters: Parameters{
			ResultFormat:      "text",
			IncrementalOutput: true,
		},
	}
	return c.sendStreamRequest(req)
}

func (c *LLMClient) GenerateResponseWithFunction(cfg *LLMConfig) (*FunctionResp, error) {
	req := &FunctionRequest{
		Model:    cfg.Model,
		Messages: cfg.LLMMessages,
		Tools:    tools.AllTools,
		Stream:   false,
	}
	return c.sendFunctionRequest(req)
}

func (c *LLMClient) sendRequest(req *LLMRequest) (*DashScopeResponse, error) {
	body, _ := json.Marshal(req)
	httpReq, _ := http.NewRequest("POST", c.baseURL, bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)
	var dashResp DashScopeResponse
	if err := json.Unmarshal(respBody, &dashResp); err != nil {
		return nil, err
	}
	if dashResp.Output.Text == "" && dashResp.Output.Choices == nil {
		return nil, errors.New("API返回了空的选择")
	}
	return &dashResp, nil
}

func (c *LLMClient) sendFunctionRequest(req *FunctionRequest) (*FunctionResp, error) {
	url := "https://dashscope.aliyuncs.com/compatible-mode/v1/chat/completions"
	body, _ := json.Marshal(req)
	log.Printf("运行到这，请求参数: %s", string(body))

	httpReq, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)
	log.Printf("原始响应体: %s", string(respBody))
	var functionResp FunctionResp
	if err := json.Unmarshal(respBody, &functionResp); err != nil {
		return nil, err
	}
	if len(functionResp.Choices) == 0 {
		return nil, errors.New("API返回了空的选择")
	}
	return &functionResp, nil
}

func (c *LLMClient) sendStreamRequest(req *LLMRequest) (<-chan *DashScopeResponse, error) {
	body, _ := json.Marshal(req)
	httpReq, _ := http.NewRequest("POST", c.baseURL, bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	httpReq.Header.Set("Accept", "text/event-stream")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("API返回错误状态码 %d: %s", resp.StatusCode, string(respBody))
	}

	streamChan := make(chan *DashScopeResponse)
	go func() {
		defer resp.Body.Close()
		defer close(streamChan)
		reader := bufio.NewReader(resp.Body)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				break
			}
			if strings.HasPrefix(line, "data:") {
				data := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
				if data != "[DONE]" {
					var dashResp DashScopeResponse
					if err := json.Unmarshal([]byte(data), &dashResp); err == nil {
						streamChan <- &dashResp
					}
				}
			}
		}
	}()
	return streamChan, nil
}
