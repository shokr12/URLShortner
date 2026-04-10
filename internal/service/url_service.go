// internal/service/url_service.go
package service

import (
	"context"
	"errors"
	"math/rand"
	"time"
	"urlshortner/internal/model"
	"urlshortner/internal/repo"

	"github.com/jackc/pgx/v5/pgconn"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type UrlService interface {
	ShortenURL(ctx context.Context, originalURL string, customKey string, ttl time.Duration) (string, error)
	GetOriginalURL(ctx context.Context, shortKey string) (string, error)
	UpdateURL(ctx context.Context, shortKey string, newURL string, ttl time.Duration) error
	GenerateRandomKey(length int) string
}

type urlService struct {
	cacheRepo repo.CacheRepo
	repo      repo.UrlRepo
}

func NewUrlService(cacheRepo repo.CacheRepo, pgRepo repo.UrlRepo) UrlService {
	return urlService{cacheRepo: cacheRepo, repo: pgRepo}
}

func (u urlService) GenerateRandomKey(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func (u urlService) ShortenURL(ctx context.Context, originalURL string, customKey string, ttl time.Duration) (string, error) {
	code := customKey
	if code == "" {
		code = u.GenerateRandomKey(8)
	}

	url := model.URL{
		OriginalURL: originalURL,
		ShortURL:    code,
		ClickCount:  0,
	}

	if err := u.repo.CreateURL(&url, ctx); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return "", ErrDuplicateKey
		}
		return "", err
	}

	if err := u.cacheRepo.SetUrl(ctx, code, originalURL, ttl); err != nil {
		return "", err
	}

	return code, nil
}

func (u urlService) GetOriginalURL(ctx context.Context, shortKey string) (string, error) {
	url, err := u.repo.GetURLByShortKey(shortKey, ctx)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return "", ErrNotFound
		}
		return "", err 
	}
	return url, nil
}

func (u urlService) UpdateURL(ctx context.Context, shortKey string, newURL string, ttl time.Duration) error {
	if err := u.repo.UpdateURL(shortKey, newURL, ctx); err != nil {
		return err
	}
	return u.cacheRepo.SetUrl(ctx, shortKey, newURL, ttl)
}
