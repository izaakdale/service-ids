package dsredis

import (
	"context"

	"github.com/izaakdale/service-ids/internal/datastore"
)

func (c *client) Insert(ctx context.Context, rec datastore.Record) error {
	cmd := c.store.HSet(rec.Keys.PK, rec.Keys.SK, rec.Data)
	if cmd.Err() != nil {
		return cmd.Err()
	}
	return nil
}
