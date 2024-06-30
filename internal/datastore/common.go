package datastore

import "fmt"

var (
	ErrNotFound = fmt.Errorf("not found")
)

type Keys struct {
	PK string `dynamodbav:"PK" json:"pk"`
	SK string `dynamodbav:"SK" json:"sk"`
}

type IDRecord struct {
	Keys
	ID string `dynamodbav:"data" json:"data"`
}
