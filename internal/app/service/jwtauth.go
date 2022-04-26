package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"github.com/golang-jwt/jwt"
	"github.com/romm80/gophermart.git/internal/app"
	"github.com/romm80/gophermart.git/internal/app/models"
	"github.com/romm80/gophermart.git/internal/app/server"
	"github.com/romm80/gophermart.git/internal/app/storage"
	"time"
)

type tokenClaims struct {
	jwt.StandardClaims
	Login string `json:"login"`
}

type JWTAuth struct {
	store storage.AuthStore
}

func NewAuth(store storage.AuthStore) *JWTAuth {
	return &JWTAuth{store: store}
}

func (a *JWTAuth) CreateUser(user models.User) error {
	if user.Login == "" || user.Password == "" {
		return app.ErrInvalidLoginOrPassword
	}

	hashPass, err := hashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashPass
	return a.store.CreateUser(user)
}

func (a *JWTAuth) GenerateToken(user models.User) (string, error) {
	hashPass, err := hashPassword(user.Password)
	if err != nil {
		return "", err
	}
	user.Password = hashPass

	if err := a.store.GetUser(user); err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.Login,
	})

	return token.SignedString(server.CFG.Key)
}

func (a *JWTAuth) ParseToken(tokenStr string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return server.CFG.Key, nil
	})
	if err != nil {
		return "", err
	}
	if !token.Valid {
		return "", app.ErrTokenIsNotValid
	}
	claim, ok := token.Claims.(*tokenClaims)
	if !ok {
		return "", app.ErrTokenIsNotValid
	}

	return claim.Login, nil
}

func hashPassword(password string) (string, error) {
	h := hmac.New(sha256.New, server.CFG.Key)
	if _, err := h.Write([]byte(password)); err != nil {
		return "", err
	}
	res := h.Sum(nil)
	return hex.EncodeToString(res), nil
}
