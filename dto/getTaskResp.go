package dto

import (
	"github.com/google/uuid"
	"time"
)

type GetTaskResp struct {
	Id        uuid.UUID `json:"id"`
	Label     string    `json:"label"`
	Text      string    `json:"text"`
	Subject   string    `json:"subject"`
	Teacher   string    `json:"teacher"`
	Keypoint  bool      `json:"keypoint"`
	Points    int       `json:"points"`
	Completed bool      `json:"completed"`
	Deadline  time.Time `json:"deadline"`
}
