package repo

import(
	"context"
	"time"
	"github.com/redis/go-redis/v9"
)

type CacheRepo interface{
	SetUrl(ctx context.Context,key string, value string, ttl time.Duration) error
	GetUrl(ctx context.Context,key string) (string, error)
	GetAndDeleteUrl(ctx context.Context,key string) error
	IncrementUrl(ctx context.Context,key string) error
}

type cacheRepo struct{
	rdb *redis.Client
}

func NewCacheRepo(rdb *redis.Client) CacheRepo {
	return &cacheRepo{rdb: rdb}
}

func (c *cacheRepo) SetUrl(ctx context.Context,key string, value string, ttl time.Duration) error {
	ctx =context.Background()
	return c.rdb.Set(ctx,"URL:"+key,value,ttl).Err()
}

func (c *cacheRepo) GetUrl(ctx context.Context,key string) (string, error) {
	ctx =context.Background()
	return c.rdb.Get(ctx, "URL:"+key).Result()
}

func (c *cacheRepo) GetAndDeleteUrl(ctx context.Context,key string) error {
	ctx =context.Background()
	return c.rdb.GetDel(ctx, "URL:"+key).Err()
}

func (c *cacheRepo) IncrementUrl(ctx context.Context,key string) error {
	ctx =context.Background()
	return c.rdb.Incr(ctx, "URL:"+key).Err()
}