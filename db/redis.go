package db

import (
	"github.com/redis/go-redis/v9"
	"golang.org/x/net/context"
	"time"
)

var RedisClient *redis.Client

func InitRedisCluster(ctx context.Context, addrs []string) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:         "redis-15783.crce178.ap-east-1-1.ec2.redns.redis-cloud.com:15783",
		PoolSize:     20,
		MinIdleConns: 5,
		DialTimeout:  5 * time.Second,
		Password:     "XvUdVJaROg81pp4vdKtPr6cqJaLtMwLd",
	})
	RedisClient = redisClient
}
