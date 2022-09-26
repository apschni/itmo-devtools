package service

import (
	"context"
	"github.com/google/uuid"
	apperrors "homeworkdeliverysystem/errors"
	"homeworkdeliverysystem/model"
	"homeworkdeliverysystem/pkg/repository"
	"log"
)

type AuthService struct {
	userRepo   repository.User
	tokensRepo repository.Token
}

func NewAuthService(userRepo repository.User, tokensRepo repository.Token) *AuthService {
	return &AuthService{userRepo: userRepo, tokensRepo: tokensRepo}
}

func (a *AuthService) SignUp(ctx context.Context, user *model.User) (string, error) {
	pw, err := hashPassword(user.Password)
	if err != nil {
		log.Printf("Unable to signup user for username: %v\n", user.Username)
		return "", apperrors.NewInternal()
	}
	user.Password = pw

	id, err := a.userRepo.Create(ctx, *user)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (a *AuthService) SignIn(ctx context.Context, user *model.User) error {
	uFetched, err := a.userRepo.FindByUsername(ctx, user.Username)
	if err != nil {
		return apperrors.NewAuthorization("Invalid email and password combination")
	}

	match, err := comparePasswords(uFetched.Password, user.Password)

	if err != nil {
		return apperrors.NewInternal()
	}

	if !match {
		return apperrors.NewAuthorization("Invalid email and password combination")
	}

	*user = *uFetched
	return nil
}

func (a *AuthService) SignOut(ctx context.Context, id uuid.UUID) error {
	return a.tokensRepo.DeleteUserRefreshToken(ctx, id.String())
}
