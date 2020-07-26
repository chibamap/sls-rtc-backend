package ddb

import (
	"errors"
	"log"
	"os"
	"sls-rtc-backend/pkg/connection"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// ConnectionTable class
type ConnectionTable struct {
	// TableName which use dynamodb table name
	TableName string
	ddb       *dynamodb.DynamoDB
}

// NewConnectionTable instance from table name
func NewConnectionTable() (*ConnectionTable, error) {

	ddbSession, err := NewDynamoDBSession()
	if err != nil {
		return nil, err
	}

	tableName := os.Getenv("TABLE_NAME")

	if tableName == "" {
		return nil, errors.New("tabne name was not set")
	}

	conn := &ConnectionTable{
		ddb:       ddbSession,
		TableName: tableName,
	}

	return conn, nil
}

// Put connection item to dynamo db
func (table *ConnectionTable) Put(conn *connection.Connection) error {
	attributeValues, _ := dynamodbattribute.MarshalMap(conn)

	input := &dynamodb.PutItemInput{
		Item:      attributeValues,
		TableName: aws.String(table.TableName),
	}

	_, err := table.ddb.PutItem(input)
	return err
}

// Delete connection item from dynamo db
func (table *ConnectionTable) Delete(conn *connection.Connection) error {
	attributeValues, _ := dynamodbattribute.MarshalMap(conn)
	log.Println(attributeValues)
	input := &dynamodb.DeleteItemInput{
		Key:       attributeValues,
		TableName: aws.String(table.TableName),
	}

	_, err := table.ddb.DeleteItem(input)
	return err
}
