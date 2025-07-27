package repository

import (
	"Voice-Assistant/internal/model"
	"gorm.io/gorm"
)

type AssistantRepository struct {
	db *gorm.DB
}

func NewAssistantRepository(db *gorm.DB) *AssistantRepository {
	return &AssistantRepository{db: db}
}

func (r *AssistantRepository) Save(assistant *model.Assistant) error {
	return r.db.Create(assistant).Error
}

func (r *AssistantRepository) SelectAll() ([]model.Assistant, error) {
	var assistants []model.Assistant
	err := r.db.Find(&assistants).Error
	return assistants, err
}

func (r *AssistantRepository) FindByID(id string) (*model.Assistant, error) {
	var assistant model.Assistant
	err := r.db.First(&assistant, "id = ?", id).Error
	if err != nil {
		return nil, err // 当出错时返回nil和错误
	}
	return &assistant, nil
}

func (r *AssistantRepository) UpdateByID(id string, assistant *model.Assistant) error {
	return r.db.Model(&model.Assistant{}).Where("id = ?", id).Updates(assistant).Error
}

func (r *AssistantRepository) DeleteByID(id string) error {
	return r.db.Where("id = ?", id).Delete(&model.Assistant{}).Error
}
