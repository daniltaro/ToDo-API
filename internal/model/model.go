package model

import (
	"time"

	"github.com/google/uuid" 
)

type Task struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	IsDone      bool      `json:"is_done"`
	Deadline    time.Time `json:"deadline"`
	CreatedAt   time.Time `json:"created_at"`
}

type TaskCondittion struct {
	IsDone bool `json:"is_done"`
}
