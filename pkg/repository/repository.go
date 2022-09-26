package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"homeworkdeliverysystem/dto"
	"homeworkdeliverysystem/model"
	"time"
)

type Task interface {
	Create(ctx context.Context, task model.Task) (string, error)
	GetByUserId(ctx context.Context, id uuid.UUID) ([]dto.GetTaskResp, error)
	UpdateFileNameOnMultipleTasks(ctx context.Context, ids pq.StringArray, fileName string) error
	GetFileNameById(ctx context.Context, id string) (string, error)
	Open(ctx context.Context, id uuid.UUID) error
	Close(ctx context.Context, id uuid.UUID) error
}

type User interface {
	Create(ctx context.Context, user model.User) (string, error)
	FindByUsername(ctx context.Context, username string) (*model.User, error)
	FindById(ctx context.Context, id uuid.UUID) (*model.User, error)
	GetByGroupNumber(ctx context.Context, number string) ([]dto.GetStudentsResp, error)
}

type Token interface {
	SetRefreshToken(ctx context.Context, userID string, tokenID string, expiresIn time.Duration) error
	DeleteRefreshToken(ctx context.Context, userID string, prevTokenID string) error
	DeleteUserRefreshToken(ctx context.Context, userID string) error
}

type Group interface {
	GetSubjectsByGroupNumber(ctx context.Context, number string) ([]string, error)
}

type Repository struct {
	Task
	User
	Token
	Group
}

func NewRepository(dataSources *dataSources) *Repository {
	return &Repository{
		User:  NewUserRepository(dataSources.DB),
		Token: NewTokenRepository(dataSources.RedisClient),
		Task:  NewTaskRepository(dataSources.DB),
		Group: NewGroupRepository(dataSources.DB),
	}
}
