package domain

import "errors"

var (
	ErrUserNotFound = errors.New("user doesn't exists")
	ErrUserExists   = errors.New("user already exists")
)
