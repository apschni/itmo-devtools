package model

import (
	"github.com/google/uuid"
	"time"
)

type Task struct {
	Id         uuid.UUID `json:"id" db:"id"`
	Subject    string    `json:"subject" db:"subject"`
	Label      string    `json:"label" db:"label"`
	Text       string    `json:"text" db:"text"`
	Deadline   time.Time `json:"deadline" db:"deadline"`
	Points     int       `json:"points" db:"points"`
	IsKeyPoint bool      `json:"is_key_point" db:"is_key_point"`
	Closed     bool      `json:"closed" db:"closed"`
	TeacherId  uuid.UUID `json:"teacher_id" db:"teacher_id"`
	StudentId  uuid.UUID `json:"student_id" db:"student_id"`
	FileName   string    `json:"file_name" db:"file_name"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}
