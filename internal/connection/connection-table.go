package connection

import (
	"errors"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/hogehoge-banana/sls-rtc-backend/internal/connection/ddb"
)

// Table class
type Table struct {
	// TableName which use dynamodb table name
	TableName string
	ddb       *dynamodb.DynamoDB
}

// NewTable instance from table name
func NewTable() (*Table, error) {

	ddbSession, err := ddb.NewDynamoDBSession()
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
func (table *Table) Put(conn *Connection) error {
	attributeValues, _ := dynamodbattribute.MarshalMap(conn)

	input := &dynamodb.PutItemInput{
		Item:      attributeValues,
		TableName: aws.String(table.TableName),
	}

	_, err := table.ddb.PutItem(input)
	return err
}

// Delete connection item from dynamo db
func (table *Table) Delete(conn *Connection) error {
	attributeValues, _ := dynamodbattribute.MarshalMap(conn)

	input := &dynamodb.DeleteItemInput{
		Key:       attributeValues,
		TableName: aws.String(table.TableName),
	}

	_, err := table.ddb.DeleteItem(input)
	return err
}

// ScanAll from connection table
func (table *Table) ScanAll() ([]Connection, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(table.TableName),
	}
	output, err := table.ddb.Scan(input)
	if err != nil {
		return nil, err
	}
	recs := []Connection{}
	dynamodbattribute.UnmarshalListOfMaps(output.Items, &recs)
	return recs, nil
}
