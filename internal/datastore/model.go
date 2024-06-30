package datastore

type Keys struct {
	PK string `dynamodbav:"PK" json:"id"`
	SK string `dynamodbav:"SK" json:"type"`
}

type IDRecord struct {
	Keys
	ID string `dynamodbav:"id" json:"type_id,omitempty"`
}
