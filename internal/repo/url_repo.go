package repo

import (
	"context"
	"urlshortner/internal/model"

	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
)

type UrlRepo struct {
	pgconn *pgx.Conn
	rdb    *redis.Client
}

func NewUrlRepo(pgconn *pgx.Conn, rdb *redis.Client) *UrlRepo {
	return &UrlRepo{pgconn: pgconn, rdb: rdb}
}

func (r *UrlRepo) CreateURL(url *model.URL, ctx context.Context) error {
	query := "INSERT INTO url (url, short_url, click_count) VALUES ($1, $2, $3)"
	_, err := r.pgconn.Exec(ctx, query, url.OriginalURL, url.ShortURL, url.ClickCount)
	return err
}

func (r *UrlRepo) GetURLByShortKey(shortKey string, ctx context.Context) (string, error) {
	// 1. Try to get from Redis (Cache Hit)
	if r.rdb != nil {
		val, err := r.rdb.Get(ctx, "URL:"+shortKey).Result()
		if err == nil {
			return val, nil
		}
	}

	// 2. If not in Redis, get from Database (Cache Miss)
	var originalURL string
	query := "SELECT url FROM url WHERE short_url = $1"
	err := r.pgconn.QueryRow(ctx, query, shortKey).Scan(&originalURL)
	if err != nil {
		return "", err
	}

	return originalURL, nil
}

func (r *UrlRepo) IncrementClickCount(shortKey string, ctx context.Context) error {
	query := "UPDATE url SET click_count = click_count + 1 WHERE short_url = $1"
	_, err := r.pgconn.Exec(ctx, query, shortKey)
	if err != nil {
		return err
	}

	// Increment the click count in Redis
	if r.rdb != nil {
		return r.rdb.Incr(ctx, "clicks:"+shortKey).Err()
	}
	return nil
}


func (r *UrlRepo) UpdateURL(shortKey string, newURL string, ctx context.Context) error {
	query := "UPDATE url SET url = $1 WHERE short_url = $2"
	_, err := r.pgconn.Exec(ctx, query, newURL, shortKey)
	if err != nil {
		return err
	}
	return nil
}

func (r *UrlRepo) GetAllURLs(ctx context.Context) ([]model.URL, error) {
	query := "SELECT url, short_url, click_count FROM url ORDER BY created_at DESC"
	rows, err := r.pgconn.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var urls []model.URL
	for rows.Next() {
		var u model.URL
		if err := rows.Scan(&u.OriginalURL, &u.ShortURL, &u.ClickCount); err != nil {
			return nil, err
		}
		urls = append(urls, u)
	}
	return urls, nil
}

func (r *UrlRepo) GetStats(ctx context.Context) (map[string]interface{}, error) {
	var totalURLs int
	var totalClicks int
	err := r.pgconn.QueryRow(ctx, "SELECT COUNT(*), COALESCE(SUM(click_count), 0) FROM url").Scan(&totalURLs, &totalClicks)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total_urls":   totalURLs,
		"total_clicks": totalClicks,
	}, nil
}

func (r *UrlRepo) DeleteURL(shortKey string, ctx context.Context) error {
	query := "DELETE FROM url WHERE short_url = $1"
	_, err := r.pgconn.Exec(ctx, query, shortKey)
	if err != nil {
		return err
	}
	if r.rdb != nil {
		return r.rdb.Del(ctx, "URL:"+shortKey, "clicks:"+shortKey).Err()
	}
	return nil
}
