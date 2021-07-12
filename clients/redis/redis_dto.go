package redis

import (
	"backend_task/conf"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/go-redis/redis/v8"
)

// Connect, method for connect to redis
func (r *rds) Connect(confs *conf.AppConfig) error {
	var err error

	once.Do(func() {

		r.db = redis.NewClient(&redis.Options{
			DB:       r.id,
			Addr:     net.JoinHostPort(confs.Client.Redis.Host, confs.Client.Redis.Port),
			Username: confs.Client.Redis.User,
			Password: confs.Client.Redis.Password,
		})

		if err = r.db.Ping(context.Background()).Err(); err != nil {
		}
	})

	return err
}

func (r *rds) setClient(client *redis.Client) {
	r.db = client
}

func (r *rds) GetDB() (Redis, error) {
	if r.db == nil {
		return nil, fmt.Errorf("db is not initialized")
	}
	return r, nil
}

// Set meth a new key,value
func (r *rds) Set(ctx context.Context, key string, value interface{}, duration time.Duration) error {
	p, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.db.Set(ctx, key, p, duration).Err()
}

// Get meth, get value with key
func (r *rds) Get(ctx context.Context, key string, dest interface{}) error {
	p, err := r.db.Get(ctx, key).Result()

	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(p), &dest)
}

// Del for delete keys in redis
func (r *rds) Del(ctx context.Context, key ...string) error {
	_, err := r.db.Del(ctx, key...).Result()
	if err != nil {
		return err
	}
	return nil

}
