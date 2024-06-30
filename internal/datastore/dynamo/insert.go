package dsdynamo

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/izaakdale/service-ids/internal/datastore"
)

func (c *client) Insert(ctx context.Context, rec datastore.IDRecord) error {
	item, err := attributevalue.MarshalMap(rec)
	if err != nil {
		return err
	}
	_, err = c.store.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(c.table),
		Item:      item,
	})
	return err
}
