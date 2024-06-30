package dsredis

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/go-redis/redis"
	"github.com/izaakdale/service-ids/internal/datastore"
)

func pkMatcher(pk string) string {
	return fmt.Sprintf("%s%s*", pk, sep)
}

func (c *client) List(ctx context.Context, pk string) ([]datastore.IDRecord, error) {
	cmd := c.store.Scan(0, pkMatcher(pk), -1)
	if cmd.Err() != nil {
		if errors.Is(cmd.Err(), redis.Nil) {
			return nil, datastore.ErrNotFound
		}
		return nil, cmd.Err()
	}
	keys, _ := cmd.Val()
	idRecs := make([]datastore.IDRecord, len(keys))
	for idx, key := range keys {
		compositeKey := strings.Split(key, sep)
		if len(compositeKey) != 2 {
			continue
		}
		val := c.store.Get(key)
		if val.Err() != nil {
			if errors.Is(val.Err(), redis.Nil) {
				continue
			}
			return nil, val.Err()
		}
		idRecs[idx] = datastore.IDRecord{
			Keys: datastore.Keys{
				PK: compositeKey[0],
				SK: compositeKey[1],
			},
			ID: val.Val(),
		}
	}

	return idRecs, nil
}
