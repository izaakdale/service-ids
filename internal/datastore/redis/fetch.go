package dsredis

import (
	"context"
	"errors"

	"github.com/go-redis/redis"
	"github.com/izaakdale/service-ids/internal/datastore"
)

func (c *client) Fetch(ctx context.Context, keys datastore.Keys) (*datastore.Record, error) {
	cmd := c.store.HGet(keys.PK, keys.SK)
	if cmd.Err() != nil {
		if errors.Is(cmd.Err(), redis.Nil) {
			return nil, datastore.ErrNotFound
		}
		return nil, cmd.Err()
	}
	idRec := datastore.Record{
		Keys: keys,
		Data: cmd.Val(),
	}
	return &idRec, nil
}
