package ddb

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// DDB class
type DDB struct {
	DdbSession *dynamodb.DynamoDB
	TableName  string
}

// New instance from table name
func New(tableName string) (*DDB, error) {
	session, err := session.NewSession()
	if err != nil {
		return nil, err
	}

	ddb := &DDB{
		DdbSession: dynamodb.New(session),
		TableName:  tableName,
	}
	return ddb, nil
}

// NewDynamoDBSession make dynamodb session
func NewDynamoDBSession() (*dynamodb.DynamoDB, error) {
	session, err := session.NewSession()
	if err != nil {
		return nil, err
	}

	return dynamodb.New(session), nil
}
