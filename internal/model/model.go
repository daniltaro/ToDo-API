package model

import "time"

type Task struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	Title       string    `json:"title" gorm:"not null"`
	Description string    `json:"description"`
	IsDone      bool    `json:"isDone" gorm:"default:false"`
	Deadline    time.Time `json:"deadline"`
	CreatedAt   time.Time `json:"createdAt" gorm:"default:CURRENT_TIMESTAMP"`
}
