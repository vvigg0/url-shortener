package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/vvigg0/l3/url-shortener/internal/model"
	"github.com/wb-go/wbf/zlog"
)

func (r *Repository) GetAllLinks() (*model.GetAllResponse, error) {
	resp := model.GetAllResponse{}
	query := `SELECT full_url,short_code FROM urls`
	rows, err := r.data.Master.Query(query)
	if err != nil {
		zlog.Logger.Error().Msgf("ошибка запроса всех ссылок к БД: %v", err)
		return nil, ErrInternal
	}
	defer rows.Close()
	for rows.Next() {
		link := model.Link{}
		if err := rows.Scan(&link.FullURL, &link.ShortCode); err != nil {
			zlog.Logger.Error().Msgf("ошибка скана строк в model.Link: %v", err)
			return nil, ErrInternal
		}
		resp.Links = append(resp.Links, link)
	}
	if err := rows.Err(); err != nil {
		zlog.Logger.Error().Msgf("ошибка строк БД: %v", err)
		return nil, ErrInternal
	}
	return &resp, nil
}

func (r *Repository) GetFullURL(code string) (int, string, error) {
	now := time.Now().UTC()
	query := `SELECT id,full_url FROM urls WHERE short_code=$1 AND expires_at > $2`

	id, url := 0, ""
	if err := r.data.Master.QueryRow(query, code, now).
		Scan(&id, &url); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, "", ErrNotFound
		}
		zlog.Logger.Error().Msgf("ошибка не получилось достать full_url: %v", err)
		return 0, "", ErrInternal
	}
	return id, url, nil
}

func (r *Repository) GetRawAnalytics(shortCode string) (*model.RawAnalyticsResponse, error) {
	analytics := model.RawAnalyticsResponse{ShortCode: shortCode}

	queryURLs := `SELECT id,full_url
			  FROM urls WHERE short_code= $1 
			  GROUP BY id,full_url`

	var urlID int
	if err := r.data.Master.QueryRow(queryURLs, shortCode).Scan(
		&urlID,
		&analytics.FullURL,
	); err != nil {
		zlog.Logger.Error().Msgf("ошибка запроса аналитики из urls: %v", err)
		return nil, ErrInternal
	}

	if urlID == 0 {
		return nil, ErrNotFound
	}

	queryURLVisits := `SELECT visited_at,user_agent
					   FROM url_visits WHERE url_id=$1`

	rows, err := r.data.Master.Query(queryURLVisits, urlID)
	if err != nil {
		zlog.Logger.Error().Msgf("ошибка запроса аналитики из url_visits: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var visitedAt time.Time
		var userAgent string
		if err := rows.Scan(&visitedAt, &userAgent); err != nil {
			zlog.Logger.Error().Msgf("ошибка скана строк аналитики: %v", err)
			return nil, ErrInternal
		}
		analytics.Visits = append(analytics.Visits, model.Visit{
			VisitedAt: visitedAt,
			UserAgent: userAgent})
		analytics.Clicks++
	}
	if err := rows.Err(); err != nil {
		zlog.Logger.Error().Msgf("ошибка строк аналитики БД: %v", err)
		return nil, ErrInternal
	}

	return &analytics, nil
}

func (r *Repository) GetTimeAnalytics(shortCode string, period string) (*model.TimeAnalyticsResponse, error) {
	analytics := model.TimeAnalyticsResponse{ShortCode: shortCode}

	queryURLs := `SELECT id,full_url FROM urls WHERE short_code=$1`

	var urlID int
	if err := r.data.Master.QueryRow(queryURLs, shortCode).Scan(&urlID, &analytics.FullURL); err != nil {
		zlog.Logger.Error().Msgf("ошибка скана значений из БД в структуру TimeAnalyticsResponse: %v", err)
		return nil, ErrInternal
	}

	queryURLVisits := queryByPeriod(period)

	rows, err := r.data.Master.Query(queryURLVisits, urlID)
	if err != nil {
		zlog.Logger.Error().Msgf("ошибка запроса аналитики к БД: %v", err)
		return nil, ErrInternal
	}
	for rows.Next() {
		var timeStat model.TimeStat
		if err := rows.Scan(&timeStat.Period, &timeStat.Clicks); err != nil {
			zlog.Logger.Error().Msgf("ошибка скана строк аналитики: %v", err)
			return nil, ErrInternal
		}
		analytics.Data = append(analytics.Data, timeStat)
	}
	if err := rows.Err(); err != nil {
		zlog.Logger.Error().Msgf("ошибка строк аналитики БД: %v", err)
		return nil, ErrInternal
	}

	return &analytics, nil
}

func (r *Repository) GetUserAgentAnalytics(shortCode string) (*model.UserAgentAnalyticsResponse, error) {
	analytics := model.UserAgentAnalyticsResponse{ShortCode: shortCode}

	queryURLs := `SELECT id,full_url FROM urls WHERE short_code=$1`

	var urlID int

	if err := r.data.Master.QueryRow(queryURLs, shortCode).Scan(&urlID, &analytics.FullURL); err != nil {
		zlog.Logger.Error().Msgf("ошибка скана значений из БД в структуру UserAgentAnalyticsResponse: %v", err)
		return nil, ErrInternal
	}

	queryURLVisits := `SELECT user_agent,COUNT(*) AS count
						FROM url_visits
						WHERE url_id=$1
						GROUP BY user_agent`

	rows, err := r.data.Master.Query(queryURLVisits, urlID)
	if err != nil {
		zlog.Logger.Error().Msgf("ошибка запроса аналитики к БД: %v", err)
		return nil, ErrInternal
	}
	for rows.Next() {
		var userStat model.UserAgentStat
		if err := rows.Scan(&userStat.UserAgent, &userStat.Clicks); err != nil {
			zlog.Logger.Error().Msgf("ошибка скана строк аналитики: %v", err)
			return nil, ErrInternal
		}
		analytics.Data = append(analytics.Data, userStat)
	}
	if err := rows.Err(); err != nil {
		zlog.Logger.Error().Msgf("ошибка строк аналитики БД: %v", err)
		return nil, ErrInternal
	}

	return &analytics, nil
}
func queryByPeriod(period string) string {
	switch period {
	case "day":
		return `SELECT TO_CHAR(visited_at, 'YYYY-MM-DD') AS period, COUNT(*) AS clicks
				FROM url_visits
				WHERE url_id = $1
				GROUP BY period
				ORDER BY period`
	case "month":
		return `SELECT TO_CHAR(visited_at, 'YYYY-MM') AS period, COUNT(*) AS clicks
				FROM url_visits
				WHERE url_id = $1
				GROUP BY period
				ORDER BY period`
	default:
		return ""
	}
}
