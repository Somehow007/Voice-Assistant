package handler

import (
	"Voice-Assistant/internal/repository"
	"Voice-Assistant/internal/service"
	"gorm.io/gorm"
)

type HandlerRegistry struct {
	AssistantHandler *AssistantHandler
	HistoryHandler   *HistoryHandler
}

var Registry *HandlerRegistry

func InitHandlers(db *gorm.DB) {

	assistantRepo := repository.NewAssistantRepository(db)
	historyRepo := repository.NewHistoryRepo(db)

	assistantService := service.NewAssistantService(assistantRepo)
	historyService := service.NewHistoryService(historyRepo)

	Registry = &HandlerRegistry{
		AssistantHandler: NewAssistantHandler(assistantService),
		HistoryHandler:   NewHistoryHandler(historyService),
	}

}
