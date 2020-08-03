package connection

import (
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// table class
type table struct {
	tableName *string
	ddb       *dynamodb.DynamoDB
}

// Connection dynamodb record structure
type tableRecord struct {
	PK           string `json:"pk"`
	ConnectionID string `json:"connectionID"`
	RoomID       string `json:"roomID"`
}

const (
	tablePK                        = "pk"
	pkPrefixConn                   = "connectionID#"
	pkPrefixRoom                   = "roomID#"
	blankValue                     = "-"
	conditionExpressionPKExists    = "attribute_exists(pk)"
	conditionExpressionPKNotExists = "attribute_not_exists(pk)"
)

var ddbsession *dynamodb.DynamoDB
var once sync.Once

// Newtable instance from table name
func newTable(tableName string) *table {
	once.Do(func() {
		sess := session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
		}))
		ddbsession = dynamodb.New(sess)
	})

	table := &table{
		ddb:       ddbsession,
		tableName: aws.String(tableName),
	}

	return table
}

func (table *table) GetConn(connectionID string) (*Connection, error) {
	pk := pkPrefixConn + connectionID
	item, err := table.find(pk)
	if err != nil || item == nil {
		return nil, err
	}
	conn := Connection{}
	err = dynamodbattribute.UnmarshalMap(item, &conn)
	return &conn, err
}

func (table *table) GetRoom(roomID string) (*Room, error) {
	pk := pkPrefixRoom + roomID
	item, err := table.find(pk)
	if err != nil || item == nil {
		return nil, err
	}
	room := Room{}
	err = dynamodbattribute.UnmarshalMap(item, &room)
	return &room, err
}

func (table *table) find(pk string) (map[string]*dynamodb.AttributeValue, error) {
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			tablePK: {S: aws.String(pk)},
		},
		TableName: table.tableName,
	}
	res, err := table.ddb.GetItem(input)
	if err != nil {
		return nil, err
	}
	return res.Item, nil

}

// Put connection item to dynamo db
func (table *table) PutConnection(conn *Connection) error {
	attributeValues, _ := dynamodbattribute.MarshalMap(conn)

	input := &dynamodb.PutItemInput{
		Item:      attributeValues,
		TableName: table.tableName,
	}

	_, err := table.ddb.PutItem(input)
	return err
}

// Put connection item to dynamo db
func (table *table) PutNewConnection(conn *Connection) error {
	attributeValues, _ := dynamodbattribute.MarshalMap(conn)

	input := &dynamodb.PutItemInput{
		Item:                attributeValues,
		TableName:           table.tableName,
		ConditionExpression: aws.String(conditionExpressionPKNotExists),
	}

	_, err := table.ddb.PutItem(input)
	return err
}

// Put connection item to dynamo db
func (table *table) PutRecord(record *tableRecord) error {
	attributeValues, _ := dynamodbattribute.MarshalMap(record)

	input := &dynamodb.PutItemInput{
		Item:      attributeValues,
		TableName: table.tableName,
	}

	_, err := table.ddb.PutItem(input)
	return err
}

// Put connection item to dynamo db
func (table *table) PutNewRoom(room *Room) error {
	attributeValues, _ := dynamodbattribute.MarshalMap(room)

	input := &dynamodb.PutItemInput{
		Item:                attributeValues,
		TableName:           table.tableName,
		ConditionExpression: aws.String(conditionExpressionPKNotExists),
	}

	_, err := table.ddb.PutItem(input)
	return err
}

func (table *table) UpdateConnectionRoomID(connID string, roomID string) error {
	pk := pkPrefixConn + connID

	input := &dynamodb.UpdateItemInput{
		TableName: table.tableName,
		Key: map[string]*dynamodb.AttributeValue{
			tablePK: {S: aws.String(pk)}},
		UpdateExpression: aws.String("SET roomID = :roomID"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":roomID": {S: aws.String(roomID)},
		},
		ConditionExpression: aws.String(conditionExpressionPKExists),
	}

	_, err := table.ddb.UpdateItem(input)
	return err
}

func (table *table) DeleteConnection(connectionID string) error {
	pk := pkPrefixConn + connectionID
	return table.delete(pk)
}
func (table *table) DeleteRoom(roomID string) error {
	pk := pkPrefixRoom + roomID
	return table.delete(pk)
}

// Delete connection item from dynamo db
func (table *table) delete(pk string) error {

	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			tablePK: {S: aws.String(pk)},
		},
		TableName: table.tableName,
	}

	_, err := table.ddb.DeleteItem(input)
	return err
}

// ScanAll from connection table
func (table *table) ScanAll() ([]Connection, error) {
	input := &dynamodb.ScanInput{
		TableName: table.tableName,
	}
	output, err := table.ddb.Scan(input)
	if err != nil {
		return nil, err
	}
	recs := []Connection{}
	dynamodbattribute.UnmarshalListOfMaps(output.Items, &recs)
	return recs, nil
}

// for debug use
func (table *table) consistentRead(pk string) (*tableRecord, error) {
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			tablePK: {S: aws.String(pk)},
		},
		TableName:      table.tableName,
		ConsistentRead: aws.Bool(true),
	}
	res, err := table.ddb.GetItem(input)
	if err != nil {
		return nil, err
	}
	if res.Item == nil {
		return nil, nil
	}
	connRecord := tableRecord{}
	err = dynamodbattribute.UnmarshalMap(res.Item, &connRecord)
	return &connRecord, err
}
