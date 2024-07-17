package dsredis

import (
	"context"
	"errors"

	"github.com/go-redis/redis"
	"github.com/izaakdale/service-ids/internal/datastore"
)

func (c *client) List(ctx context.Context, pk string) ([]datastore.Record, error) {
	cmd := c.store.HScan(pk, 0, "*", -1)
	if cmd.Err() != nil {
		if errors.Is(cmd.Err(), redis.Nil) {
			return nil, datastore.ErrNotFound
		}
		return nil, cmd.Err()
	}

	idRecs := make([]datastore.Record, 0)
	it := cmd.Iterator()
	for it.Next() {
		currentSK := it.Val()
		it.Next()

		idRecs = append(idRecs, datastore.Record{
			Keys: datastore.Keys{
				PK: pk,
				SK: currentSK,
			},
			Data: it.Val(),
		})
	}
	if len(idRecs) == 0 {
		return nil, datastore.ErrNotFound
	}
	return idRecs, nil
}
