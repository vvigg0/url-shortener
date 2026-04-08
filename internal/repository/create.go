package repository

import (
	"time"

	"github.com/wb-go/wbf/zlog"
)

func (r *Repository) InsertShortLink(fullURL, code string, createdAt time.Time) error {
	expiresAt := createdAt.Add(r.urlsTTL)
	query := `INSERT INTO urls (full_url,short_code,created_at,expires_at)
			 VALUES($1,$2,$3,$4) ON CONFLICT (short_code) DO NOTHING`
	res, err := r.data.Master.Exec(query, fullURL, code, createdAt, expiresAt)
	if err != nil {
		zlog.Logger.Error().Msgf("ошибка вставки короткого кода: %v", err)
		return ErrInternal
	}
	aff, err := res.RowsAffected()
	if err != nil {
		zlog.Logger.Error().Msgf("ошибка RowsAffected: %v", err)
		return ErrInternal
	}
	if aff == 0 {
		return ErrAlreadyExist
	}
	return nil
}

func (r *Repository) InsertAnalytics(id int, agent string,
	clickTime time.Time) error {
	queryURLVisits := `INSERT INTO url_visits 
			(url_id,user_agent,visited_at) 
			VALUES($1,$2,$3)`

	_, err := r.data.Master.Exec(queryURLVisits, id, agent, clickTime)
	if err != nil {
		zlog.Logger.Error().Msgf("ошибка вставки аналитики: %v", err)
		return ErrInternal
	}
	return nil
}
