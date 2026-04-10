package config

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
)

type Config struct {
	BaseURL    string
	ServerAddr string
	DefaultTTL time.Duration
	DB         DBConfig
	Redis      RedisConfig
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	URL      string // takes precedence if set
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
}

func Load() (*Config, error) {
	return &Config{
		BaseURL:    getEnv("BASE_URL", "http://localhost:8080"),
		ServerAddr: getEnv("SERVER_ADDR", ":8080"),
		DefaultTTL: 24 * time.Hour,
		DB: DBConfig{
			URL:      os.Getenv("DATABASE_URL"),
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     requireEnv("DB_USER"),
			Password: requireEnv("DB_PASS"),
			Name:     requireEnv("DB_NAME"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: os.Getenv("REDIS_PASS"),
		},
	}, nil
}

func InitRedis(cfg *Config) (*redis.Client, error) {
	addr := fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port)

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.Redis.Password,
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis ping failed at %s: %w", addr, err)
	}

	return rdb, nil
}

func InitDB(cfg *Config) (*pgx.Conn, error) {
	connStr := cfg.DB.URL
	if connStr == "" {
		connStr = fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s",
			cfg.DB.User, cfg.DB.Password,
			cfg.DB.Host, cfg.DB.Port, cfg.DB.Name,
		)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		return nil, fmt.Errorf("postgres connect failed: %w", err)
	}

	if err := conn.Ping(ctx); err != nil {
		conn.Close(context.Background())
		return nil, fmt.Errorf("postgres ping failed: %w", err)
	}

	return conn, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func requireEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		panic(fmt.Sprintf("required environment variable %q is not set", key))
	}
	return v
}
