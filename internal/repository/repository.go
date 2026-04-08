package repository

import (
	"time"

	"github.com/wb-go/wbf/dbpg"
)

type Repository struct {
	data    *dbpg.DB
	urlsTTL time.Duration
}

func New(db *dbpg.DB, ttl time.Duration) *Repository {
	return &Repository{data: db, urlsTTL: ttl}
}
