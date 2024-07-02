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
	Data string `dynamodbav:"data" json:"data"`
}
