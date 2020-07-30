package connection

import (
	"errors"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// table class
type table struct {
	// tableName which use dynamodb table name
	tableName string
	ddb       *dynamodb.DynamoDB
}

var ddbsession *dynamodb.DynamoDB
var once sync.Once

// Newtable instance from table name
func newTable() (*table, error) {
	once.Do(func() {
		sess := session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
		}))
		ddbsession = dynamodb.New(sess)
	})

	tableName := os.Getenv("TABLE_NAME")

	if tableName == "" {
		return nil, errors.New("tabne name was not set")
	}

	table := &table{
		ddb:       ddbsession,
		tableName: tableName,
	}

	return table, nil
}
func (table *table) Get(connectionID string) (*Connection, error) {
	conn := New(connectionID)
	attributeValues, _ := dynamodbattribute.MarshalMap(conn)
	input := &dynamodb.GetItemInput{
		Key:       attributeValues,
		TableName: &table.tableName,
	}
	res, err := table.ddb.GetItem(input)
	if err != nil {
		return nil, err
	}
	connRecord := Connection{}
	err = dynamodbattribute.UnmarshalMap(res.Item, &connRecord)
	return &connRecord, err
}

// Put connection item to dynamo db
func (table *table) Put(conn *Connection) error {
	attributeValues, _ := dynamodbattribute.MarshalMap(conn)

	input := &dynamodb.PutItemInput{
		Item:      attributeValues,
		TableName: aws.String(table.tableName),
	}

	_, err := table.ddb.PutItem(input)
	return err
}

// Delete connection item from dynamo db
func (table *table) Delete(conn *Connection) error {
	attributeValues, _ := dynamodbattribute.MarshalMap(conn)

	input := &dynamodb.DeleteItemInput{
		Key:       attributeValues,
		TableName: aws.String(table.tableName),
	}

	_, err := table.ddb.DeleteItem(input)
	return err
}

// ScanAll from connection table
func (table *table) ScanAll() ([]Connection, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(table.tableName),
	}
	output, err := table.ddb.Scan(input)
	if err != nil {
		return nil, err
	}
	recs := []Connection{}
	dynamodbattribute.UnmarshalListOfMaps(output.Items, &recs)
	return recs, nil
}
