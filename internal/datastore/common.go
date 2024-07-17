package datastore

import "fmt"

var (
	ErrNotFound = fmt.Errorf("not found")
)

type Keys struct {
	PK string `dynamodbav:"PK" json:"pk"`
	SK string `dynamodbav:"SK" json:"sk"`
}

type Record struct {
	Keys
	Cursor string `dynamodbav:"cursor,omitempty" json:"cursor,omitempty"`
	Data   string `dynamodbav:"data" json:"data"`
}
