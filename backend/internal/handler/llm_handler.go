package handler

import (
	"Voice-Assistant/internal/model"
	"Voice-Assistant/internal/service"
	"Voice-Assistant/pkg/llm"
	"github.com/gin-gonic/gin"
	"net/http"
)

type LLMHandler struct {
	llmService       service.LLMService
	assistantService service.AssistantService
	historyService   *service.HistoryService
}

func NewLLMHandler(llmService service.LLMService, assistantService service.AssistantService, historyService service.HistoryService) *LLMHandler {
	return &LLMHandler{llmService: llmService, assistantService: assistantService, historyService: &historyService}
}

// ChatLLM POST /api/voice-assistant/v1/llm/chat
func (h *LLMHandler) ChatLLM(c *gin.Context) {
	var req struct {
		AssistantId string `json:"assistant_id"`
		UserInput   string `json:"user_input"`
		Stream      bool   `json:"stream"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数格式错误"})
		return
	}

	assistant, _ := h.assistantService.GetAssistantById(req.AssistantId)
	histories, _ := h.historyService.GetAllHistoryByAssistantId(req.AssistantId)
	if assistant == nil || histories == nil {
		c.JSON(http.StatusBadRequest, &model.Result{
			Code:    http.StatusBadRequest,
			Message: "the  history is nil when start to chat with llm",
			Data:    nil,
		})
		return
	}

	systemPrompt := assistant.Prompt
	var llmMessages []llm.LLMMessage
	if len(histories) > 0 {
		maxHistory := 20
		total := len(histories)
		start := 0
		if total > maxHistory {
			start = total - maxHistory
		}
		recentHistory := histories[start:]
		llmMessages = model.HistoryToLLMMessages(recentHistory)
		systemPrompt = histories[total-1].Input[0].Content
	}

	if req.Stream {
		streamChan, err := h.llmService.GenerateStreamResponseWithContext(systemPrompt, req.UserInput, llmMessages)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		for resp := range streamChan {
			c.SSEvent("message", resp)
		}
	} else {
		resp, err := h.llmService.GenerateResponseWithFunctionCalling(systemPrompt, req.UserInput, llmMessages)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"result": resp})
	}
}
