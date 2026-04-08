package repository

import "errors"

var (
	ErrAlreadyExist = errors.New("такой код уже существует")
	ErrInternal     = errors.New("внутренняя ошибка")
	ErrNotFound     = errors.New("не найдено")
)
