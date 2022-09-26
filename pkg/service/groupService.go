package service

import (
	"context"
	"homeworkdeliverysystem/dto"
	"homeworkdeliverysystem/pkg/repository"
)

type GroupService struct {
	groupRepo repository.Group
	userRepo  repository.User
}

func NewGroupService(groupRepo repository.Group, userRepo repository.User) *GroupService {
	return &GroupService{groupRepo: groupRepo, userRepo: userRepo}
}

func (g *GroupService) GetSubjectsByNumber(ctx context.Context, number string) ([]string, error) {
	return g.groupRepo.GetSubjectsByGroupNumber(ctx, number)
}

func (g *GroupService) GetStudentsByNumber(ctx context.Context, number string) ([]dto.GetStudentsResp, error) {
	return g.userRepo.GetByGroupNumber(ctx, number)
}
