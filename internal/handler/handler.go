package handler

import (
	"time"

	"github.com/vvigg0/l3/url-shortener/internal/model"
)

type shortenerService interface {
	GetAllLinks() (*model.GetAllResponse, error)
	CreateShortLink(u, customCode string) (string, error)
	GetFullURL(string) (int, string, error)
	WriteAnalytics(int, string, time.Time) error
	GetAnalytics(shortCode string, mode string) (any, error)
}

type Handler struct {
	srvc shortenerService
}

func New(s shortenerService) *Handler {
	return &Handler{srvc: s}
}
