package service

import (
	"Voice-Assistant/internal/model"
	"Voice-Assistant/internal/repository"
	"github.com/google/uuid"
	"log"
	"time"
)

type AssistantService struct {
	repo *repository.AssistantRepository
}

func NewAssistantService(repo *repository.AssistantRepository) *AssistantService {
	return &AssistantService{repo: repo}
}

func (s *AssistantService) SaveAssistant(name, prompt, description string) (*model.Assistant, error) {
	assistant := &model.Assistant{
		ID:          uuid.New().String(),
		Name:        name,
		Prompt:      prompt,
		Description: description,
		CreatedAt:   time.Now().Format("2006-01-02 15:04:05"),
		UpdatedAt:   time.Now().Format("2006-01-02 15:04:05"),
	}
	err := s.repo.Save(assistant)
	if err != nil {
		log.Printf("Failer to save assistant: %v", err)
		return nil, err
	}
	return assistant, nil
}

// todo 校验有待加强
func (s *AssistantService) GetAllAssistant() ([]model.Assistant, error) {
	result, err := s.repo.SelectAll()
	if err != nil {
		log.Printf("select assistants failed: %v", err)
	}
	return result, err
}

func (s *AssistantService) UpdateAssistantById(id, name, prompt, description string) (*model.Assistant, error) {
	assistant := &model.Assistant{
		Name:        name,
		Prompt:      prompt,
		Description: description,
		UpdatedAt:   time.Now().String(),
	}

	err := s.repo.UpdateByID(id, assistant)
	if err != nil {
		log.Printf("Failed to update assistant: %v", err)
		return nil, err
	}

	updatedAssistant, err := s.repo.FindByID(id)
	if err != nil {
		log.Printf("Failed to select the updated assistant: %v", err)
		return nil, err
	}

	return updatedAssistant, nil
}

func (s *AssistantService) GetAssistantById(id string) (*model.Assistant, error) {
	assistant, err := s.repo.FindByID(id)
	if err != nil {
		log.Printf("Failed to select assistant by id: %v", err)
		return nil, err
	}
	return assistant, nil
}

func (s *AssistantService) DeleteAssistantById(id string) error {
	err := s.repo.DeleteByID(id)
	if err != nil {
		log.Printf("Failed to delete assistant by id: %v", err)
		return err
	}
	return nil
}
