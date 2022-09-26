package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"homeworkdeliverysystem/dto"
	apperrors "homeworkdeliverysystem/errors"
	"homeworkdeliverysystem/model"
	"log"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(DB *sqlx.DB) *UserRepository {
	return &UserRepository{db: DB}
}

func (u *UserRepository) FindById(ctx context.Context, id uuid.UUID) (*model.User, error) {
	user := &model.User{}
	query := "SELECT * FROM users WHERE id=$1"
	err := u.db.GetContext(ctx, user, query, id)
	if err != nil {
		return user, apperrors.NewNotFound("id", id.String())
	}
	return user, err
}

func (u *UserRepository) Create(ctx context.Context, user model.User) (string, error) {
	var id uuid.UUID
	newUUID, err := uuid.NewUUID()
	if err != nil {
		return "", apperrors.NewInternal()
	}
	query := "INSERT INTO users (id, full_name, group_number, username, password_hash, role) " +
		"VALUES ($1, $2, $3, $4, $5, $6) RETURNING id"
	err = u.db.GetContext(ctx, &id, query, newUUID, user.FullName, user.GroupNumber, user.Username, user.Password, user.Role)
	if err != nil {
		if err, ok := err.(*pq.Error); ok && err.Code.Name() == "unique_violation" {
			log.Printf("Could not create a user with username: %v. Reason: %v\n", user.Username, err.Code.Name())
			return "", apperrors.NewConflict("username", user.Username)
		}

		log.Printf("Could not create a user with username: %v. Reason: %v\n", user.Username, err)
		return "", apperrors.NewInternal()
	}

	return id.String(), nil
}

func (u *UserRepository) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	user := &model.User{}
	query := "SELECT * FROM users WHERE username=$1"
	err := u.db.GetContext(ctx, user, query, username)
	if err != nil {
		log.Printf("Unable to get user with username: %v. Err: %v\n", username, err)
		return user, apperrors.NewNotFound("username", username)
	}
	return user, nil
}

func (u *UserRepository) GetByGroupNumber(ctx context.Context, number string) ([]dto.GetStudentsResp, error) {
	var students []dto.GetStudentsResp

	query := "SELECT id, full_name, username FROM users WHERE group_number=$1 ORDER BY full_name"

	err := u.db.SelectContext(ctx, &students, query, number)
	if err != nil {
		log.Printf("Unable to get user with group number: %v. Err: %v\n", number, err)
		return students, apperrors.NewNotFound("group_number", number)
	}

	return students, nil
}
