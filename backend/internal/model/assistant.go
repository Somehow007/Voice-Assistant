package model

type Assistant struct {
	ID          string `gorm:"primaryKey"` // primary key
	Name        string // the name of assistant
	Description string // describe the assistant
	Prompt      string // give LLM scene
	CreatedAt   string // time when assistant is created
	UpdatedAt   string // time when assistant is updated
}
