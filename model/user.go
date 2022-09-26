package model

import (
	"github.com/google/uuid"
)

type User struct {
	Id          uuid.UUID `json:"id" db:"id"`
	FullName    string    `json:"full_name"  db:"full_name"`
	GroupNumber string    `json:"group_number" db:"group_number"`
	Username    string    `json:"username"  db:"username"`
	Password    string    `json:"-" db:"password_hash"`
	Role        string    `json:"role" db:"role"`
}
