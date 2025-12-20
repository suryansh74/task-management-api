package models

import "time"

type Task struct {
	ID        string    `json:"task_id"`
	Title     string    `json:"title" validate:"required,min=3,max=50"`
	Content   string    `json:"content" validate:"max=500"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
