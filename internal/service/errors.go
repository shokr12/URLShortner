// internal/service/errors.go
package service

import "errors"

var (
	ErrNotFound     = errors.New("url not found")
	ErrExpired      = errors.New("url has expired")
	ErrDuplicateKey = errors.New("custom key already exists")
	ErrInvalidURL   = errors.New("invalid url")
)
