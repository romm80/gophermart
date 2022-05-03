package app

import (
	"errors"
	"net/http"
)

var (
	ErrInvalidLoginOrPassword   = errors.New("invalid login or password")
	ErrLoginIsUsed              = errors.New("login is already used")
	ErrTokenIsNotValid          = errors.New("token is not valid")
	ErrInvalidRequestFormat     = errors.New("invalid request format")
	ErrInvalidOrderFormat       = errors.New("invalid order format")
	ErrOrderUploaded            = errors.New("order uploaded")
	ErrOrderUploadedAnotherUser = errors.New("order uploaded another user")
	ErrNotEnoughFunds           = errors.New("not enough funds on the account")
	ErrInvalidUserID            = errors.New("invalid user id")
)

func ErrStatusCode(err error) int {
	switch {
	case errors.Is(err, ErrLoginIsUsed) || errors.Is(err, ErrOrderUploadedAnotherUser):
		return http.StatusConflict
	case errors.Is(err, ErrInvalidRequestFormat):
		return http.StatusBadRequest
	case errors.Is(err, ErrInvalidLoginOrPassword) || errors.Is(err, ErrInvalidUserID) || errors.Is(err, ErrTokenIsNotValid):
		return http.StatusUnauthorized
	case errors.Is(err, ErrInvalidOrderFormat):
		return http.StatusUnprocessableEntity
	case errors.Is(err, ErrNotEnoughFunds):
		return http.StatusPaymentRequired
	default:
		return http.StatusInternalServerError
	}
}
