package service

import "errors"

var (
	ErrInvalidURL        = errors.New("invalid url")
	ErrInvalidCode       = errors.New("invalid short code")
	ErrInvalidCustomCode = errors.New("invalid custom code")
	ErrBadStatus         = errors.New("bad target status")
	ErrInternal          = errors.New("internal error")
)
