package model

import (
	"errors"
)

var (
	ErrNotFound      = errors.New("not found")
	ErrInvalidInput  = errors.New("invalid input")
	ErrFailedToHash  = errors.New("failed to hash")
	ErrWrongPassword = errors.New("wrong password")
	ErrCartIsEmpty   = errors.New("cart is empty")
	ErrAlreadyExists = errors.New("already exists")
	ErrNoAccess      = errors.New("no access")
)

type ValidationError struct {
	msg string
}

func NewValidationError(msg string) *ValidationError {
	return &ValidationError{
		msg: msg,
	}
}

func (ve *ValidationError) Error() string {
	return ve.msg
}
