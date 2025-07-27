package repository

import (
	"Voice-Assistant/internal/model"
	"gorm.io/gorm"
)

type HistoryRepo struct {
	db *gorm.DB
}

func NewHistoryRepo(db *gorm.DB) *HistoryRepo {
	return &HistoryRepo{db: db}
}

func (r *HistoryRepo) Save(history *model.History) error {
	return r.db.Create(history).Error
}

func (r *HistoryRepo) SelectByAssistantId(assistantId string) ([]model.History, error) {
	var histories []model.History
	err := r.db.Where("assistant_id = ?", assistantId).Order("created_at DESC").Find(&histories).Error
	return histories, err
}

func (r *HistoryRepo) DeleteByHistoryId(historyId string) error {
	return r.db.Where("id = ?", historyId).Delete(&model.History{}).Error
}

func (r *HistoryRepo) DeleteByAssistantId(assistantId string) error {
	return r.db.Where("assistant_id = ?", assistantId).Delete(&model.History{}).Error
}
