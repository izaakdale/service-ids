package dsdynamo

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/izaakdale/service-ids/internal/datastore"
)

func (c *client) Fetch(ctx context.Context, keys datastore.Keys) (*datastore.Record, error) {
	in, err := attributevalue.MarshalMap(keys)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal keys: %w", err)
	}
	out, err := c.store.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(c.table),
		Key:       in,
	})

	if err != nil {
		nfe := &types.ResourceNotFoundException{}
		if errors.As(err, &nfe) {
			return nil, datastore.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get item: %w", err)
	}
	idRec := datastore.Record{}
	if err = attributevalue.UnmarshalMap(out.Item, &idRec); err != nil {
		return nil, fmt.Errorf("failed to unmarshal item: %w", err)
	}
	if idRec.Data == "" {
		return nil, datastore.ErrNotFound
	}
	return &idRec, nil
}
