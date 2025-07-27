package service

import (
	"Voice-Assistant/internal/model"
	"Voice-Assistant/internal/repository"
	"errors"
	"github.com/google/uuid"
	"log"
	"time"
)

type HistoryService struct {
	repo *repository.HistoryRepo
}

func NewHistoryService(repo *repository.HistoryRepo) *HistoryService {
	return &HistoryService{repo: repo}
}

func (s *HistoryService) SaveHistory(history *model.History) error {
	if history == nil {
		err := errors.New("the new history is nil")
		log.Println(err)
		return err
	}
	history.ID = uuid.New().String()
	history.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
	return s.repo.Save(history)
}

func (s *HistoryService) GetAllHistoryByAssistantId(assistantId string) ([]model.History, error) {
	return s.repo.SelectByAssistantId(assistantId)
}

func (s *HistoryService) DeleteByHistoryId(id string) error {
	return s.repo.DeleteByHistoryId(id)
}

func (s *HistoryService) DeleteByAssistantId(assistantId string) error {
	if assistantId == "" {
		err := errors.New("delete this assistant's histories filed. Cause the assistant_id is empty.")
		return err
	}
	return s.repo.DeleteByAssistantId(assistantId)
}
