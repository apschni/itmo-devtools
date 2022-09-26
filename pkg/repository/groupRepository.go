package repository

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"homeworkdeliverysystem/errors"
)

type GroupRepository struct {
	db *sqlx.DB
}

func NewGroupRepository(db *sqlx.DB) *GroupRepository {
	return &GroupRepository{db: db}
}

func (g *GroupRepository) GetSubjectsByGroupNumber(ctx context.Context, number string) ([]string, error) {
	var subjects pq.StringArray

	query := "SELECT subjects FROM groups WHERE number=$1"

	err := g.db.GetContext(ctx, &subjects, query, number)
	if err != nil {
		return nil, errors.NewNotFound("number", number)
	}
	return subjects, nil
}
