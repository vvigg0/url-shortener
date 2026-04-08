package service

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/vvigg0/l3/url-shortener/internal/repository"
	"github.com/vvigg0/l3/url-shortener/pkg/generator"
	"github.com/wb-go/wbf/zlog"
)

var customCodeRegexp = regexp.MustCompile(`^[A-Za-z0-9_-]+$`)

func (s *Service) CreateShortLink(u, customCode string) (string, error) {
	if err := checkURL(u); err != nil {
		return "", err
	}
	var link string

	createdAt := time.Now().UTC()

	if customCode == "" {
		var code string
		var e error
		for range 5 {
			var err error
			code, err = generator.GenerateShortCode(8)
			if err != nil {
				zlog.Logger.Error().Msgf("ошибка генерации short code: %v", err)
				return "", ErrInternal
			}

			createdAt := time.Now().UTC()

			if e = s.repo.InsertShortLink(u, code, createdAt); e != nil {
				if errors.Is(e, repository.ErrAlreadyExist) {
					continue
				}
				return "", e
			} else {
				break
			}
		}
		if e != nil {
			return "", ErrInternal
		}

		link = fmt.Sprintf("%s/s/%s", s.rootAddress, code)
	} else {
		if err := checkCustomCode(customCode); err != nil {
			return "", err
		}
		if err := s.repo.InsertShortLink(u, customCode, createdAt); err != nil {
			return "", err
		}
		link = fmt.Sprintf("%s/s/%s", s.rootAddress, customCode)
	}
	return link, nil
}

func checkURL(u string) error {
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		zlog.Logger.Error().Msgf("ошибка создания request: %v", err)
		return ErrInternal
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "ru-RU,ru;q=0.9,en-US;q=0.8,en;q=0.7")
	req.Header.Set("Referer", "https://www.google.com/")
	req.Header.Set("Connection", "keep-alive")

	client := http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		zlog.Logger.Error().Msgf("ошибка выполнения request: %v", err)
		return ErrInvalidURL
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 && resp.StatusCode <= 599 {
		zlog.Logger.Info().Msgf("ссылка вернула: %d", resp.StatusCode)
		return ErrBadStatus
	}

	return nil
}

func (s *Service) WriteAnalytics(id int, agent string,
	clickTime time.Time) error {
	if err := s.repo.InsertAnalytics(id, agent, clickTime); err != nil {
		return err
	}

	return nil
}

func checkCustomCode(code string) error {
	code = strings.TrimSpace(code)

	if code == "" {
		return ErrInvalidCustomCode
	}
	if len(code) < 3 || len(code) > 32 {
		return ErrInvalidCustomCode
	}
	if !customCodeRegexp.MatchString(code) {
		return ErrInvalidCustomCode
	}

	return nil
}
