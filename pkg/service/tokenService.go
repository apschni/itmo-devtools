package service

import (
	"context"
	"crypto/rsa"
	"github.com/google/uuid"
	apperrors "homeworkdeliverysystem/errors"
	"homeworkdeliverysystem/model"
	"homeworkdeliverysystem/pkg/repository"
	"log"
)

type TokenService struct {
	repo                  repository.Token
	PrivKey               *rsa.PrivateKey
	PubKey                *rsa.PublicKey
	RefreshSecret         string
	IDExpirationSecs      int64
	RefreshExpirationSecs int64
}

func NewTokenService(repo repository.Token, privKey *rsa.PrivateKey, pubKey *rsa.PublicKey, refreshSecret string, IDExpirationSecs int64, refreshExpirationSecs int64) *TokenService {
	return &TokenService{repo: repo, PrivKey: privKey, PubKey: pubKey, RefreshSecret: refreshSecret, IDExpirationSecs: IDExpirationSecs, RefreshExpirationSecs: refreshExpirationSecs}
}

func (s *TokenService) ValidateIdToken(tokenString string) (*model.User, error) {
	claims, err := validateIdToken(tokenString, s.PubKey)

	if err != nil {
		log.Printf("Unable to validate or parse idToken - Error: %v\n", err)
		return nil, apperrors.NewAuthorization("Unable to verify user from idToken")
	}

	return claims.User, nil
}

func (s *TokenService) NewPairFromUser(ctx context.Context, user *model.User, prevTokenID string) (*model.TokenPair, error) {
	if prevTokenID != "" {
		if err := s.repo.DeleteRefreshToken(ctx, user.Id.String(), prevTokenID); err != nil {
			log.Printf("Could not delete previous refreshToken for uid: %v, tokenID: %v\n", user.Id, prevTokenID)
			return nil, err
		}
	}

	idToken, err := generateToken(user, s.PrivKey, s.IDExpirationSecs)

	if err != nil {
		log.Printf("Error generating idToken for id: %v. Error: %v\n", user.Id, err.Error())
		return nil, apperrors.NewInternal()
	}

	refreshToken, err := generateRefreshToken(user.Id, s.RefreshSecret, s.RefreshExpirationSecs)

	if err != nil {
		log.Printf("Error generating refreshToken for id: %v. Error: %v\n", user.Id, err.Error())
		return nil, apperrors.NewInternal()
	}

	if err := s.repo.SetRefreshToken(ctx, user.Id.String(), refreshToken.Id.String(), refreshToken.ExpiresIn); err != nil {
		log.Printf("Error storing tokenID for uid: %v. Error: %v\n", user.Id, err.Error())
		return nil, apperrors.NewInternal()
	}

	return &model.TokenPair{
		IDToken:      model.IDToken{SS: idToken},
		RefreshToken: model.RefreshToken{SS: refreshToken.SS, ID: refreshToken.Id, UID: user.Id},
	}, nil
}

func (s *TokenService) ValidateRefreshToken(refreshTokenString string) (*model.RefreshToken, error) {
	claims, err := validateRefreshToken(refreshTokenString, s.RefreshSecret)

	if err != nil {
		log.Printf("Unable to validate or parse refreshToken for token string: %s\n%v\n", refreshTokenString, err)
		return nil, apperrors.NewAuthorization("Unable to verify user from refresh token")
	}

	tokenUUID, err := uuid.Parse(claims.Id)

	if err != nil {
		log.Printf("Claims ID could not be parsed as UUID: %s\n%v\n", claims.Id, err)
		return nil, apperrors.NewAuthorization("Unable to verify user from refresh token")
	}

	return &model.RefreshToken{
		SS:  refreshTokenString,
		ID:  tokenUUID,
		UID: claims.UID,
	}, nil
}
