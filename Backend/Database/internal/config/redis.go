package config

// import (
// 	"context"
// 	"time"

// 	"github.com/redis/go-redis/v9"
// )

// type RedisCondig struct {
// 	Addr        string
// 	Password    string
// 	User        string
// 	DB          int
// 	MaxRetries  int
// 	DialTimeout time.Duration
// 	Timeout     time.Duration
// }

// func NewRedisConfig(ctx context.Context, cfg RedisCondig) (*redis.Client, error) {
// 	db := redis.NewClient(&redis.Options{
// 		Addr: cfg.Addr,

// 		DB: cfg.DB,

// 		MaxRetries:   cfg.MaxRetries,
// 		DialTimeout:  cfg.DialTimeout,
// 		ReadTimeout:  cfg.Timeout,
// 		WriteTimeout: cfg.Timeout,
// 	})

// 	if err := db.Ping(ctx).Err(); err != nil {
// 		return nil, err
// 	}

// 	return db, nil
// }
