package dto

import (
	"github.com/google/uuid"
	"time"
)

type CreateTaskReq struct {
	Label      string    `json:"label" binding:"required"`
	Subject    string    `json:"subject" binding:"required"`
	Text       string    `json:"text"`
	Deadline   time.Time `json:"deadline" binding:"required"`
	Points     int       `json:"points"`
	IsKeyPoint bool      `json:"is_key_point" binding:"required"`
	StudentId  uuid.UUID `json:"student_id" binding:"required"`
}
