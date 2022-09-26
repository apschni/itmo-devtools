package service

import (
	"crypto/rsa"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"homeworkdeliverysystem/model"
	"log"
	"time"
)

type idTokenCustomClaims struct {
	User *model.User `json:"user"`
	jwt.StandardClaims
}

func generateToken(user *model.User, key *rsa.PrivateKey, exp int64) (string, error) {
	unixTime := time.Now().Unix()
	tokenExp := unixTime + exp

	claims := idTokenCustomClaims{
		User: user,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  unixTime,
			ExpiresAt: tokenExp,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	ss, err := token.SignedString(key)
	if err != nil {
		log.Println("Failed to sign id token string")
		return "", err
	}

	return ss, nil
}

func validateIdToken(tokenString string, key *rsa.PublicKey) (*idTokenCustomClaims, error) {
	claims := &idTokenCustomClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("Id token is invalid")
	}

	claims, ok := token.Claims.(*idTokenCustomClaims)

	if !ok {
		return nil, fmt.Errorf("Id token valid but couldn't parse claims")
	}

	return claims, nil
}

type refreshTokenData struct {
	SS        string
	Id        uuid.UUID
	ExpiresIn time.Duration
}

type refreshTokenCustomClaims struct {
	UID uuid.UUID `json:"uid"`
	jwt.StandardClaims
}

func generateRefreshToken(uid uuid.UUID, key string, exp int64) (*refreshTokenData, error) {
	currentTime := time.Now()
	tokenExp := currentTime.Add(time.Duration(exp) * time.Second)
	tokenID, err := uuid.NewRandom() // v4

	if err != nil {
		log.Println("Failed to generate refresh token Id")
		return nil, err
	}

	claims := refreshTokenCustomClaims{
		UID: uid,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  currentTime.Unix(),
			ExpiresAt: tokenExp.Unix(),
			Id:        tokenID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(key))

	if err != nil {
		log.Println("Failed to sign refresh token string")
		return nil, err
	}

	return &refreshTokenData{
		SS:        ss,
		Id:        tokenID,
		ExpiresIn: tokenExp.Sub(currentTime),
	}, nil
}

func validateRefreshToken(tokenString string, key string) (*refreshTokenCustomClaims, error) {
	claims := &refreshTokenCustomClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("Refresh token is invalid")
	}

	claims, ok := token.Claims.(*refreshTokenCustomClaims)

	if !ok {
		return nil, fmt.Errorf("Refresh token valid but couldn't parse claims")
	}

	return claims, nil
}
