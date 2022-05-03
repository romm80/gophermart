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
	UserID int `json:"user_id"`
}

type JWTAuth struct {
	store storage.AuthStore
}

func NewAuth(store storage.AuthStore) *JWTAuth {
	return &JWTAuth{store: store}
}

func (a *JWTAuth) CreateUser(user models.User) (string, error) {
	if user.Login == "" || user.Password == "" {
		return "", app.ErrInvalidLoginOrPassword
	}

	hashPass, err := hashPassword(user.Password)
	if err != nil {
		return "", err
	}
	user.Password = hashPass
	userID, err := a.store.CreateUser(user)
	if err != nil {
		return "", err
	}

	return generateToken(userID)
}

func (a *JWTAuth) LoginUser(user models.User) (string, error) {
	hashPass, err := hashPassword(user.Password)
	if err != nil {
		return "", err
	}
	user.Password = hashPass

	userID, err := a.store.GetUserID(user)
	if err != nil {
		return "", err
	}

	return generateToken(userID)
}

func (a *JWTAuth) ParseToken(tokenStr string) (int, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return server.CFG.Key, nil
	})
	if err != nil {
		return 0, err
	}
	if !token.Valid {
		return 0, app.ErrTokenIsNotValid
	}
	claim, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, app.ErrTokenIsNotValid
	}
	if err := a.ValidUserID(claim.UserID); err != nil {
		return 0, err
	}

	return claim.UserID, nil
}

func (a *JWTAuth) ValidUserID(userID int) error {
	return a.store.ValidUserID(userID)
}

func hashPassword(password string) (string, error) {
	h := hmac.New(sha256.New, server.CFG.Key)
	if _, err := h.Write([]byte(password)); err != nil {
		return "", err
	}
	res := h.Sum(nil)
	return hex.EncodeToString(res), nil
}

func generateToken(userID int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		userID,
	})

	return token.SignedString(server.CFG.Key)
}
