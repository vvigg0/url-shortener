package service

import (
	"time"

	"github.com/vvigg0/l3/url-shortener/internal/model"
)

type shortenerRepository interface {
	GetAllLinks() (*model.GetAllResponse, error)
	InsertShortLink(fullURL, code string, createdAt time.Time) error
	GetFullURL(shortCode string) (int, string, error)
	InsertAnalytics(id int, agent string, clickTime time.Time) error
	GetRawAnalytics(shortCode string) (*model.RawAnalyticsResponse, error)
	GetTimeAnalytics(shortCode string, period string) (*model.TimeAnalyticsResponse, error)
	GetUserAgentAnalytics(shortCode string) (*model.UserAgentAnalyticsResponse, error)
}

type shortenerCache interface {
	Get(code string) (int, string, error)
	Set(code string, urlID int, fullURL string) error
}
type Service struct {
	rootAddress string
	repo        shortenerRepository
	cache       shortenerCache
}

func New(rootAddress string, r shortenerRepository, c shortenerCache) *Service {
	return &Service{rootAddress: rootAddress, repo: r, cache: c}
}
