package dsdynamo

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/izaakdale/service-ids/internal/datastore"
)

func (c *client) List(ctx context.Context, pk string) ([]datastore.Record, uint64, error) {
	keyCond := expression.Key("PK").Equal(expression.Value(pk))

	expr, err := expression.NewBuilder().WithKeyCondition(keyCond).Build()
	if err != nil {
		return nil, 0, err
	}

	out, err := c.store.Query(ctx, &dynamodb.QueryInput{
		TableName:                 aws.String(c.table),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
	})
	if err != nil {
		return nil, 0, err
	}
	if len(out.Items) == 0 {
		return nil, 0, datastore.ErrNotFound
	}

	recs := make([]datastore.Record, len(out.Items))
	for idx, item := range out.Items {
		var idRec datastore.Record
		if err = attributevalue.UnmarshalMap(item, &idRec); err != nil {
			return nil, 0, fmt.Errorf("failed to unmarshal item: %w", err)
		}
		recs[idx] = idRec
	}

	return recs, 0, nil
}
