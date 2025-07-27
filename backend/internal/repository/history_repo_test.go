package repository

import (
	"Voice-Assistant/internal/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
)

func setUpTestHistoryDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"))
	assert.NoError(t, err)

	err = db.AutoMigrate(&model.History{})
	assert.NoError(t, err)

	return db
}

func TestHistoryRepo_Save(t *testing.T) {
	db := setUpTestHistoryDB(t)
	repo := NewHistoryRepo(db)

	history := &model.History{
		ID:          uuid.New().String(),
		AssistantID: "Test assistantId",
		Input: []model.Message{
			{
				Role:    "system",
				Content: "test",
			},
			{
				Role:    "user",
				Content: "test",
			},
		},
		Output: model.Output{
			Content:      "yes test",
			FinishReason: "stop",
		},
		Usage: model.Usage{
			InputTokens:  1,
			OutputTokens: 2,
			TotalTokens:  3,
		},
		CreatedAt: "2027-07-27",
	}

	err := repo.Save(history)
	assert.NoError(t, err)

	var found model.History
	err = db.First(&found, "id = ?", history.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, &found, history)
}

func TestHistoryRepo_SelectByAssistantId(t *testing.T) {
	db := setUpTestHistoryDB(t)
	repo := NewHistoryRepo(db)

	history1 := &model.History{
		ID:          uuid.New().String(),
		AssistantID: "Test assistantId",
		Input: []model.Message{
			{
				Role:    "system",
				Content: "test1",
			},
			{
				Role:    "user",
				Content: "test1",
			},
		},
		Output: model.Output{
			Content:      "yes test",
			FinishReason: "stop",
		},
		Usage: model.Usage{
			InputTokens:  1,
			OutputTokens: 2,
			TotalTokens:  3,
		},
		CreatedAt: "2027-07-27",
	}

	history2 := &model.History{
		ID:          uuid.New().String(),
		AssistantID: "Test assistantId",
		Input: []model.Message{
			{
				Role:    "system",
				Content: "test1",
			},
			{
				Role:    "user",
				Content: "test1",
			},
		},
		Output: model.Output{
			Content:      "yes test",
			FinishReason: "stop",
		},
		Usage: model.Usage{
			InputTokens:  1,
			OutputTokens: 2,
			TotalTokens:  3,
		},
		CreatedAt: "2027-07-27",
	}

	err := repo.Save(history1)
	assert.NoError(t, err)
	err = repo.Save(history2)
	assert.NoError(t, err)

	histories, err := repo.SelectByAssistantId("Test assistantId")
	assert.NoError(t, err)
	assert.Len(t, histories, 2)
}

func TestHistoryRepo_DeleteByHistoryId(t *testing.T) {
	db := setUpTestHistoryDB(t)
	repo := NewHistoryRepo(db)

	history := &model.History{
		ID:          uuid.New().String(),
		AssistantID: "Test assistantId",
		Input: []model.Message{
			{
				Role:    "system",
				Content: "test",
			},
			{
				Role:    "user",
				Content: "test",
			},
		},
		Output: model.Output{
			Content:      "yes test",
			FinishReason: "stop",
		},
		Usage: model.Usage{
			InputTokens:  1,
			OutputTokens: 2,
			TotalTokens:  3,
		},
		CreatedAt: "2027-07-27",
	}

	err := repo.Save(history)
	assert.NoError(t, err)

	err = repo.DeleteByHistoryId(history.ID)
	assert.NoError(t, err)

	// test not exist data
	err = repo.DeleteByHistoryId("no-existent-data")
	assert.NoError(t, err)
}

func TestHistoryRepo_DeleteByAssistantId(t *testing.T) {
	db := setUpTestHistoryDB(t)
	repo := NewHistoryRepo(db)

	history1 := &model.History{
		ID:          uuid.New().String(),
		AssistantID: "Test assistantId",
		Input: []model.Message{
			{
				Role:    "system",
				Content: "test1",
			},
			{
				Role:    "user",
				Content: "test1",
			},
		},
		Output: model.Output{
			Content:      "yes test",
			FinishReason: "stop",
		},
		Usage: model.Usage{
			InputTokens:  1,
			OutputTokens: 2,
			TotalTokens:  3,
		},
		CreatedAt: "2027-07-27",
	}

	history2 := &model.History{
		ID:          uuid.New().String(),
		AssistantID: "Test assistantId",
		Input: []model.Message{
			{
				Role:    "system",
				Content: "test1",
			},
			{
				Role:    "user",
				Content: "test1",
			},
		},
		Output: model.Output{
			Content:      "yes test",
			FinishReason: "stop",
		},
		Usage: model.Usage{
			InputTokens:  1,
			OutputTokens: 2,
			TotalTokens:  3,
		},
		CreatedAt: "2027-07-27",
	}

	err := repo.Save(history1)
	assert.NoError(t, err)
	err = repo.Save(history2)
	assert.NoError(t, err)

	histories, err := repo.SelectByAssistantId("Test assistantId")
	assert.NoError(t, err)
	assert.Len(t, histories, 2)

	err = repo.DeleteByAssistantId("Test assistantId")
	assert.NoError(t, err)

	histories, _ = repo.SelectByAssistantId("Test assistantId")
	assert.Len(t, histories, 0)
}

func TestNewHistoryRepo(t *testing.T) {
	db := setUpTestHistoryDB(t)
	repo := NewHistoryRepo(db)
	assert.NotNil(t, repo)
	assert.Equal(t, db, repo.db)
}
