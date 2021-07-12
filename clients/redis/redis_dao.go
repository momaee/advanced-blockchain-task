package redis

import (
	"backend_task/conf"
	"context"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

var (
	logger  *zap.Logger
	once    sync.Once
	Storage Redis = &rds{id: 0}
)

// store interface is interface for store things into redis
type Redis interface {
	Connect(*conf.AppConfig) error
	GetDB() (Redis, error)
	Set(ctx context.Context, key string, value interface{}, duration time.Duration) error
	Get(ctx context.Context, key string, dest interface{}) error
	Del(ctx context.Context, key ...string) error
	setClient(client *redis.Client)
}

// rds struct for redis client
type rds struct {
	db *redis.Client
	id int
}

func New(db int) Redis {
	return &rds{id: db}
}
