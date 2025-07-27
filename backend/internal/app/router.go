package app

import (
	"Voice-Assistant/internal/handler"
	"Voice-Assistant/internal/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// middleware
	router.Use(middleware.CORS())

	// api router
	api := router.Group("/api/voice-assistant/v1")
	{
		assistant := api.Group("/assistant")
		{
			assistant.POST("/", handler.Registry.AssistantHandler.SaveAssistant)
			assistant.GET("/", handler.Registry.AssistantHandler.GetAllAssistants)
			assistant.PUT("/:id", handler.Registry.AssistantHandler.UpdateAssistantById)
			assistant.DELETE("/:id", handler.Registry.AssistantHandler.DeleteAssistantById)
		}

		history := api.Group("/history")
		{
			history.POST("/", handler.Registry.HistoryHandler.SaveNewHistory)
			history.GET("/:id", handler.Registry.HistoryHandler.GetAllHistoriesByAssistantId)
			history.DELETE("/history_id/:id", handler.Registry.HistoryHandler.DeleteHistoryByHistoryId)
			history.DELETE("/assistant_id/:id", handler.Registry.HistoryHandler.DeleteHistoryByAssistantId)
		}

		llm := api.Group("/llm")
		{
			llm.POST("/chat", handler.Registry.LLMHandler.ChatLLM)
		}
	}
	return router
}
