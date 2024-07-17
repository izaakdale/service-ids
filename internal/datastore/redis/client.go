package dsredis

import (
	"github.com/go-redis/redis"
)

type RedisAPI interface {
	HGet(key, field string) *redis.StringCmd
	HSet(key, field string, value interface{}) *redis.BoolCmd
	Scan(cursor uint64, match string, count int64) *redis.ScanCmd
	HScan(key string, cursor uint64, match string, count int64) *redis.ScanCmd
}

func New(cli RedisAPI, table string) *client {
	return &client{
		store: cli,
		table: table,
	}
}

type client struct {
	store RedisAPI
	table string
}
