package dsredis

import (
	"context"
	"errors"

	"github.com/go-redis/redis"
	"github.com/izaakdale/service-ids/internal/datastore"
)

func (c *client) List(ctx context.Context, pk string) ([]datastore.Record, uint64, error) {
	cmd := c.store.HScan(pk, 0, "*", -1)
	if cmd.Err() != nil {
		if errors.Is(cmd.Err(), redis.Nil) {
			return nil, 0, datastore.ErrNotFound
		}
		return nil, 0, cmd.Err()
	}

	vals, curs := cmd.Val()
	idRecs := make([]datastore.Record, 0)

	for i := 0; i < len(vals); i += 2 {
		idRecs = append(idRecs, datastore.Record{
			Keys: datastore.Keys{
				PK: pk,
				SK: vals[i],
			},
			Data: vals[i+1],
		})
	}

	if len(idRecs) == 0 {
		return nil, 0, datastore.ErrNotFound
	}
	return idRecs, curs, nil
}
