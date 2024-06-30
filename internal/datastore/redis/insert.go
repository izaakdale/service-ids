package dsredis

import (
	"context"

	"github.com/izaakdale/service-ids/internal/datastore"
)

func (c *client) Insert(ctx context.Context, rec datastore.IDRecord) error {
	key, err := createCompositeKey(rec.Keys.PK, rec.Keys.SK)
	if err != nil {
		return err
	}
	cmd := c.store.Set(key, rec.ID, 0)
	if cmd.Err() != nil {
		return cmd.Err()
	}
	return nil
}
