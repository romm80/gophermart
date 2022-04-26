package app

import "errors"

var (
	ErrInvalidLoginOrPassword   = errors.New("invalid login or password")
	ErrLoginIsUsed              = errors.New("login is already used")
	ErrTokenIsNotValid          = errors.New("token is not valid")
	ErrInvalidRequestFormat     = errors.New("invalid request format")
	ErrInvalidOrderFormat       = errors.New("invalid order format")
	ErrOrderUploaded            = errors.New("order uploaded")
	ErrOrderUploadedAnotherUser = errors.New("order uploaded another user")
	ErrNotEnoughFunds           = errors.New("not enough funds on the account")
)
