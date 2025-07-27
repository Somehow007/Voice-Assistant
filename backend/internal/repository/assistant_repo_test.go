package repository

import (
	"Voice-Assistant/internal/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
)

func setUpTestAssistantDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"))
	assert.NoError(t, err)

	err = db.AutoMigrate(&model.Assistant{})
	assert.NoError(t, err)

	return db
}

func TestAssistantRepository_Save(t *testing.T) {
	db := setUpTestAssistantDB(t)
	repo := NewAssistantRepository(db)

	assistant := &model.Assistant{
		ID:          uuid.New().String(),
		Name:        "Test Assistant",
		Description: "Test Description",
		Prompt:      "Test Prompt",
		CreatedAt:   "2025-07-27",
		UpdatedAt:   "2025-07-27",
	}

	err := repo.Save(assistant)
	assert.NoError(t, err)

	var found model.Assistant
	err = db.First(&found, "id = ?", assistant.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, assistant.Name, found.Name)
}

func TestAssistantRepository_FindByID(t *testing.T) {
	db := setUpTestAssistantDB(t)
	repo := NewAssistantRepository(db)

	assistant := &model.Assistant{
		ID:          uuid.New().String(),
		Name:        "Test Assistant",
		Description: "Test Description",
		Prompt:      "Test Prompt",
		CreatedAt:   "2025-07-27",
		UpdatedAt:   "2025-07-27",
	}

	err := repo.Save(assistant)
	assert.NoError(t, err)

	found, err := repo.FindByID(assistant.ID)
	assert.NoError(t, err)
	assert.Equal(t, found.Name, assistant.Name)

	// test not exist data
	_, err = repo.FindByID("no-existent-id")
	assert.Error(t, err)
}

func TestAssistantRepository_SelectAll(t *testing.T) {
	db := setUpTestAssistantDB(t)
	repo := NewAssistantRepository(db)

	assistant1 := &model.Assistant{
		ID:          uuid.New().String(),
		Name:        "Assistant 1",
		Description: "Description 1",
		Prompt:      "Prompt 1",
		CreatedAt:   "2023-01-01",
		UpdatedAt:   "2023-01-01",
	}
	assistant2 := &model.Assistant{
		ID:          uuid.New().String(),
		Name:        "Assistant 2",
		Description: "Description 2",
		Prompt:      "Prompt 2",
		CreatedAt:   "2023-01-01",
		UpdatedAt:   "2023-01-01",
	}

	err := repo.Save(assistant1)
	assert.NoError(t, err)
	err = repo.Save(assistant2)
	assert.NoError(t, err)

	assistants, err := repo.SelectAll()
	assert.NoError(t, err)
	assert.Len(t, assistants, 2)
}

func TestAssistantRepository_DeleteByID(t *testing.T) {
	db := setUpTestAssistantDB(t)
	repo := NewAssistantRepository(db)

	assistant := &model.Assistant{
		ID:          uuid.New().String(),
		Name:        "Test Assistant",
		Description: "Test Description",
		Prompt:      "Test Prompt",
		CreatedAt:   "2025-07-27",
		UpdatedAt:   "2025-07-27",
	}

	err := repo.Save(assistant)
	assert.NoError(t, err)

	err = repo.DeleteByID(assistant.ID)
	assert.NoError(t, err)

	// test not exist data
	err = repo.DeleteByID("no-existent-data")
	assert.NoError(t, err)

}

func TestAssistantRepository_UpdateByID(t *testing.T) {
	db := setUpTestAssistantDB(t)
	repo := NewAssistantRepository(db)

	assistant := &model.Assistant{
		ID:          uuid.New().String(),
		Name:        "Test Assistant",
		Description: "Test Description",
		Prompt:      "Test Prompt",
		CreatedAt:   "2025-07-27",
		UpdatedAt:   "2025-07-27",
	}

	update := &model.Assistant{
		ID:          uuid.New().String(),
		Name:        "Update Assistant",
		Description: "Update Description",
		Prompt:      "Update Prompt",
	}

	err := repo.Save(assistant)
	assert.NoError(t, err)

	err = repo.UpdateByID(assistant.ID, update)
	assert.NoError(t, err)

	// test not exist data
	err = repo.UpdateByID("no-existent-data", update)
	assert.NoError(t, err)
}

func TestNewAssistantRepository(t *testing.T) {
	db := setUpTestAssistantDB(t)
	repo := NewAssistantRepository(db)
	assert.NotNil(t, repo)
	assert.Equal(t, db, repo.db)
}
