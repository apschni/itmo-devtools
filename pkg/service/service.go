package service

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"homeworkdeliverysystem/dto"
	"homeworkdeliverysystem/model"
	"homeworkdeliverysystem/pkg/repository"
	"io/ioutil"
	"os"
	"strconv"
)

type Authorization interface {
	SignUp(ctx context.Context, user *model.User) (string, error)
	SignIn(ctx context.Context, user *model.User) error
	SignOut(ctx context.Context, id uuid.UUID) error
}

type Task interface {
	Create(ctx context.Context, task *model.Task) (string, error)
	GetByUserId(ctx context.Context, id uuid.UUID) ([]dto.GetTaskResp, error)
	UpdateMultipleWithFile(ctx context.Context, req *dto.UpdateMultipleWithFileReq) error
	GetFileNameById(ctx context.Context, id string) (string, error)
	Open(ctx context.Context, id uuid.UUID) error
	Close(ctx context.Context, id uuid.UUID) error
}

type User interface {
	Get(ctx context.Context, uid uuid.UUID) (*model.User, error)
}

type Token interface {
	NewPairFromUser(ctx context.Context, user *model.User, prevTokenID string) (*model.TokenPair, error)
	ValidateIdToken(tokenString string) (*model.User, error)
	ValidateRefreshToken(refreshTokenString string) (*model.RefreshToken, error)
}

type Group interface {
	GetSubjectsByNumber(ctx context.Context, number string) ([]string, error)
	GetStudentsByNumber(ctx context.Context, number string) ([]dto.GetStudentsResp, error)
}

type Service struct {
	Authorization
	Task
	User
	Token
	Group
}

func NewService(repository *repository.Repository) *Service {
	priv, _ := ioutil.ReadFile(os.Getenv("PRIV_KEY_FILE"))
	privKey, _ := jwt.ParseRSAPrivateKeyFromPEM(priv)
	pub, _ := ioutil.ReadFile(os.Getenv("PUB_KEY_FILE"))
	pubKey, _ := jwt.ParseRSAPublicKeyFromPEM(pub)
	secret := os.Getenv("SECRET_KEY")
	idTokenExp := os.Getenv("ID_TOKEN_EXP")
	refreshTokenExp := os.Getenv("REFRESH_TOKEN_EXP")
	idExp, _ := strconv.ParseInt(idTokenExp, 0, 64)
	refreshExp, _ := strconv.ParseInt(refreshTokenExp, 0, 64)

	return &Service{
		Authorization: NewAuthService(repository.User, repository.Token),
		Token:         NewTokenService(repository.Token, privKey, pubKey, secret, idExp, refreshExp),
		User:          NewUserService(repository.User),
		Task:          NewTaskService(repository.Task, repository.User),
		Group:         NewGroupService(repository.Group, repository.User),
	}
}
