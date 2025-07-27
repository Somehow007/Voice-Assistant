package handler

import (
	"Voice-Assistant/internal/model"
	"Voice-Assistant/internal/service"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type HistoryHandler struct {
	service *service.HistoryService
}

type CreateHistoryReq struct {
	AssistantId string          `json:"assistant_id"`
	Input       []model.Message `json:"input"`
	Output      model.Output    `json:"output"`
	Usage       model.Usage     `json:"usage"`
}

func NewHistoryHandler(service *service.HistoryService) *HistoryHandler {
	return &HistoryHandler{service: service}
}

func (h *HistoryHandler) SaveNewHistory(c *gin.Context) {
	var req CreateHistoryReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, &model.Result{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}
	history := &model.History{
		AssistantID: req.AssistantId,
		Input:       req.Input,
		Output:      req.Output,
		Usage:       req.Usage,
	}
	log.Printf("history's assistant_id: %s", req.AssistantId)
	if err := h.service.SaveHistory(history); err != nil {
		c.JSON(http.StatusInternalServerError, &model.Result{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		})
	} else {
		c.JSON(http.StatusOK, &model.Result{
			Code:    http.StatusOK,
			Message: "save history successfully!",
			Data:    history,
		})
	}
}

func (h *HistoryHandler) GetAllHistoriesByAssistantId(c *gin.Context) {
	assistantId := c.Param("id")
	log.Printf("request id: %s", assistantId)
	histories, err := h.service.GetAllHistoryByAssistantId(assistantId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &model.Result{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		})
	}

	if len(histories) == 0 {
		c.JSON(http.StatusOK, &model.Result{
			Code:    http.StatusOK,
			Message: "there is no history",
			Data:    histories,
		})
	} else {
		c.JSON(http.StatusOK, &model.Result{
			Code:    http.StatusOK,
			Message: "select all histories successfully!",
			Data:    histories,
		})
	}
}

func (h *HistoryHandler) DeleteHistoryByHistoryId(c *gin.Context) {
	historyId := c.Param("id")

	err := h.service.DeleteByHistoryId(historyId)
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
		Message: "delete history successfully!",
		Data:    nil,
	})
}

func (h *HistoryHandler) DeleteHistoryByAssistantId(c *gin.Context) {
	assistantId := c.Param("id")
	_, err := Registry.AssistantHandler.service.GetAssistantById(assistantId)
	if err != nil {
		c.JSON(http.StatusBadRequest, &model.Result{
			Code:    http.StatusBadRequest,
			Message: "error param, please checkout your assistant_id, or the histories have been deleted.",
			Data:    nil,
		})
		return
	}
	err = h.service.DeleteByAssistantId(assistantId)
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
		Message: "delete assistant successfully!",
		Data:    nil,
	})
}
