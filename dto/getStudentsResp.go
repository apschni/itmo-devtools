package dto

import "github.com/google/uuid"

type GetStudentsResp struct {
	Id       uuid.UUID `json:"id"`
	FullName string    `json:"full_name" db:"full_name"`
	Username string    `json:"username"`
}
