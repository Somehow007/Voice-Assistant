package handler

import (
	"Voice-Assistant/internal/repository"
	"Voice-Assistant/internal/service"
	"Voice-Assistant/pkg/llm"
	"gorm.io/gorm"
)

type HandlerRegistry struct {
	AssistantHandler *AssistantHandler
	HistoryHandler   *HistoryHandler
	LLMHandler       *LLMHandler
}

var Registry *HandlerRegistry

func InitHandlers(db *gorm.DB) {

	assistantRepo := repository.NewAssistantRepository(db)
	historyRepo := repository.NewHistoryRepo(db)

	assistantService := service.NewAssistantService(assistantRepo)
	historyService := service.NewHistoryService(historyRepo)
	llmClient := llm.NewLLMClient()
	llmService := service.NewLLMService(llmClient)

	Registry = &HandlerRegistry{
		AssistantHandler: NewAssistantHandler(assistantService),
		HistoryHandler:   NewHistoryHandler(historyService),
		LLMHandler:       NewLLMHandler(llmService, *assistantService, *historyService),
	}

}
