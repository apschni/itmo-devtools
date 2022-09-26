package service

import (
	"context"
	"github.com/google/uuid"
	"homeworkdeliverysystem/model"
	"homeworkdeliverysystem/pkg/repository"
)

type UserService struct {
	repo repository.User
}

func NewUserService(repo repository.User) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Get(ctx context.Context, uid uuid.UUID) (*model.User, error) {
	u, err := s.repo.FindById(ctx, uid)
	return u, err
}
