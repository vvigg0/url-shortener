package myredis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/vvigg0/l3/url-shortener/internal/model"
	"github.com/wb-go/wbf/redis"
)

type Mredis struct {
	Client *redis.Client
	ttl    time.Duration
}

func New(client *redis.Client, urlsTTL time.Duration) *Mredis {
	return &Mredis{client, urlsTTL}
}

func (r *Mredis) Get(code string) (int, string, error) {
	raw, err := r.Client.Get(context.Background(), code)
	if err != nil {
		return 0, "", err
	}
	var cached model.CachedURL
	if err := json.Unmarshal([]byte(raw), &cached); err != nil {
		return 0, "", err
	}
	return cached.ID, cached.URL, nil
}

func (r *Mredis) Set(code string, urlID int, fullURL string) error {
	payload, err := json.Marshal(model.CachedURL{ID: urlID, URL: fullURL})
	if err != nil {
		return err
	}
	if err := r.Client.SetWithExpiration(
		context.Background(),
		code,
		string(payload),
		r.ttl); err != nil {
		return err
	}
	return nil
}
