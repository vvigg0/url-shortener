package handler

import (
	"errors"
	"net/http"

	"github.com/vvigg0/l3/url-shortener/internal/repository"
	"github.com/vvigg0/l3/url-shortener/internal/service"
)

func getErrorCode(err error) int {
	if errors.Is(err, service.ErrBadStatus) ||
		errors.Is(err, service.ErrInvalidURL) ||
		errors.Is(err, service.ErrInvalidCode) ||
		errors.Is(err, service.ErrInvalidCustomCode) {
		return http.StatusBadRequest
	}
	if errors.Is(err, repository.ErrAlreadyExist) {
		return http.StatusConflict
	}
	if errors.Is(err, repository.ErrNotFound) {
		return http.StatusNotFound
	}
	return http.StatusInternalServerError
}
