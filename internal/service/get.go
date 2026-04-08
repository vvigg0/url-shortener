package service

import (
	"errors"
	"time"

	"github.com/vvigg0/l3/url-shortener/internal/model"
	"github.com/wb-go/wbf/redis"
	"github.com/wb-go/wbf/zlog"
)

func (s *Service) GetAllLinks() (*model.GetAllResponse, error) {
	return s.repo.GetAllLinks()
}
func (s *Service) GetFullURL(shortCode string) (int, string, error) {
	now := time.Now()
	urlID, fullURL, err := s.cache.Get(shortCode)
	if err == nil {
		zlog.Logger.Info().Msgf("достали из кэша за %v", time.Since(now))
		return urlID, fullURL, nil
	}
	if !errors.Is(err, redis.NoMatches) {
		zlog.Logger.Error().Msgf("ошибка get from redis: %v", err)
	}
	now = time.Now()
	urlID, fullURL, err = s.repo.GetFullURL(shortCode)
	if err != nil {
		return 0, "", err
	}
	zlog.Logger.Info().Msgf("достали из БД за %v", time.Since(now))
	if err := s.cache.Set(shortCode, urlID, fullURL); err != nil {
		zlog.Logger.Error().Msgf("ошибка set to redis: %v", err)
	}
	return urlID, fullURL, nil
}

func (s *Service) GetAnalytics(shortCode string, mode string) (any, error) {
	switch mode {
	case "day":
		return s.repo.GetTimeAnalytics(shortCode, "day")
	case "month":
		return s.repo.GetTimeAnalytics(shortCode, "month")
	case "user-agent":
		return s.repo.GetUserAgentAnalytics(shortCode)
	default:
		return s.repo.GetRawAnalytics(shortCode)
	}
}
