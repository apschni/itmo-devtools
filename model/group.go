package model

import "github.com/lib/pq"

type Group struct {
	Label    string         `json:"label" db:"label"`
	Subjects pq.StringArray `json:"subjects" db:"subjects"`
}
