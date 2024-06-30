package dsredis

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/redis"
)

type RedisAPI interface {
	Get(key string) *redis.StringCmd
	Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Scan(cursor uint64, match string, count int64) *redis.ScanCmd
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

var sep = "-:-"

func createCompositeKey(pk, sk string) (string, error) {
	if strings.Contains(sk, sep) {
		return "", fmt.Errorf("sk cannot contain %s", sep)
	}
	if strings.Contains(pk, sep) {
		return "", fmt.Errorf("pk cannot contain %s", sep)
	}
	return strings.ToLower(fmt.Sprintf("%s%s%s", pk, sep, sk)), nil
}
