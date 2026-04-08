package repocleaner

import (
	"context"
	"time"

	"github.com/wb-go/wbf/dbpg"
	"github.com/wb-go/wbf/zlog"
)

type RepoCleaner struct {
	repo *dbpg.DB
}

func New(db *dbpg.DB) *RepoCleaner {
	return &RepoCleaner{repo: db}
}

func (c *RepoCleaner) Start(ctx context.Context) error {
	ticker := time.NewTicker(30 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			if err := c.deleteExpired(); err != nil {
				return err
			}
		}
	}
}

func (c *RepoCleaner) deleteExpired() error {
	now := time.Now().UTC()
	query := `DELETE FROM urls WHERE expires_at<=$1`
	res, err := c.repo.Master.Exec(query, now)
	if err != nil {
		return err
	}
	aff, err := res.RowsAffected()
	if err != nil {
		return err
	}
	zlog.Logger.Info().Msgf("удалено %v записей", aff)
	return nil
}
