package handler

import (
	"Voice-Assistant/internal/model"
	"Voice-Assistant/internal/service"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type AssistantHandler struct {
	service *service.AssistantService
}

func NewAssistantHandler(service *service.AssistantService) *AssistantHandler {
	return &AssistantHandler{service: service}
}

type CreateAssistantReq struct {
	Name        string `json:"name" binding:"required"` // gin
	Prompt      string `json:"prompt" binding:"required"`
	Description string `json:"description" binding:"required"`
}

type UpdateAssistantReq struct {
	Name        string `json:"name" binding:"required"` // gin
	Prompt      string `json:"prompt" binding:"required"`
	Description string `json:"description" binding:"required"`
}

func (h *AssistantHandler) SaveAssistant(c *gin.Context) {
	var req CreateAssistantReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	if req.Prompt == "" || req.Name == "" {
		log.Printf("save new assiatant failed.The parametar is empty. name: %s, prompt: %s", req.Name, req.Prompt)
		c.JSON(http.StatusBadRequest, &model.Result{
			Code:    http.StatusBadRequest,
			Message: "save new assistant failed. The parameter is empty.",
			Data:    nil,
		})
		return
	}

	assistant, err := h.service.SaveAssistant(req.Name, req.Prompt, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &model.Result{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, &model.Result{
		Code:    http.StatusOK,
		Message: "save assistant successfully!",
		Data:    assistant,
	})
}

func (h *AssistantHandler) GetAllAssistants(c *gin.Context) {
	assistants, err := h.service.GetAllAssistant()
	if err != nil {
		c.JSON(http.StatusInternalServerError, &model.Result{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	if len(assistants) == 0 {
		c.JSON(http.StatusOK, &model.Result{
			Code:    http.StatusOK,
			Message: "there is no assistant",
			Data:    assistants,
		})
		return
	}
	c.JSON(http.StatusOK, &model.Result{
		Code:    http.StatusOK,
		Message: "select all successfully!",
		Data:    assistants,
	})
}

func (h *AssistantHandler) UpdateAssistantById(c *gin.Context) {
	id := c.Param("id")

	var req UpdateAssistantReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, &model.Result{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}
	updatedAssistant, err := h.service.UpdateAssistantById(id, req.Name, req.Prompt, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &model.Result{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, &model.Result{
		Code:    http.StatusOK,
		Message: "update successfully!",
		Data:    updatedAssistant,
	})
}

func (h *AssistantHandler) DeleteAssistantById(c *gin.Context) {
	id := c.Param("id")

	err := h.service.DeleteAssistantById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &model.Result{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}
	err = Registry.HistoryHandler.service.DeleteByAssistantId(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &model.Result{
			Code:    http.StatusBadRequest,
			Message: "delete assistant failed, because delete histories failed",
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, &model.Result{
		Code:    http.StatusOK,
		Message: "delete assistant successfully!",
		Data:    nil,
	})
}
