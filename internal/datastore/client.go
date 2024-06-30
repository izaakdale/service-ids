package datastore

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var (
	ErrNotFound = fmt.Errorf("not found")
)

//go:generate mockgen -destination=mock_test.go -package=datastore_test github.com/izaakdale/service-ids/datastore DynamoAPI

type DynamoAPI interface {
	PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
}

type client struct {
	store DynamoAPI
	table string
}

func New(store DynamoAPI, table string) *client {
	return &client{
		store: store,
		table: table,
	}
}
