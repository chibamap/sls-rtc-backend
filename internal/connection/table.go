package connection

import (
	"log"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// table class
type table struct {
	// tableName which use dynamodb table name
	tableName *string
	ddb       *dynamodb.DynamoDB
}

// Connection dynamodb record structure
type tableRecord struct {
	// ConnectionID request.RequestContext.ConnectionID
	PK           string `json:"pk"`
	ConnectionID string `json:"connectionID"`
	RoomID       string `json:"roomID"`
	Username     string `json:"username"`
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
	return table.find(pk)
}

func (table *table) GetRoom(roomID string) (*Connection, error) {
	pk := pkPrefixRoom + roomID
	return table.find(pk)
}

func (table *table) find(pk string) (*Connection, error) {
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
	if res.Item == nil {
		return nil, nil
	}
	connRecord := Connection{}
	err = dynamodbattribute.UnmarshalMap(res.Item, &connRecord)
	return &connRecord, err
}

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

// Put connection item to dynamo db
func (table *table) Put(conn *Connection) error {
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
		Item:                attributeValues,
		TableName:           table.tableName,
		ConditionExpression: aws.String(conditionExpressionPKNotExists),
	}

	_, err := table.ddb.PutItem(input)
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

// TransactPut
func (table *table) PutNewRoom(room *Connection, ownerConnectionID string) (bool, error) {

	roomItem, _ := dynamodbattribute.MarshalMap(room)
	ownerPK := pkPrefixConn + ownerConnectionID

	items := []*dynamodb.TransactWriteItem{
		&dynamodb.TransactWriteItem{
			Put: &dynamodb.Put{
				TableName:           table.tableName,
				ConditionExpression: aws.String(conditionExpressionPKNotExists),
				Item:                roomItem,
			},
		},
		&dynamodb.TransactWriteItem{
			Update: &dynamodb.Update{
				TableName: table.tableName,
				Key: map[string]*dynamodb.AttributeValue{
					tablePK: {S: aws.String(ownerPK)},
				},
				UpdateExpression: aws.String("SET roomID = :roomID"),
				ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
					":roomID": {S: aws.String(room.RoomID)},
				},
				ConditionExpression: aws.String(conditionExpressionPKExists),
			},
		},
	}

	return table.transactWrite(items)
}

// TransactPut
func (table *table) transactWrite(transactionItems []*dynamodb.TransactWriteItem) (bool, error) {

	input := &dynamodb.TransactWriteItemsInput{
		TransactItems: transactionItems,
	}
	_, err := table.ddb.TransactWriteItems(input)

	if err != nil {
		switch t := err.(type) {
		case *dynamodb.TransactionCanceledException:
			log.Printf("failed to write items: %s\n%v",
				t.Message(), t.CancellationReasons)
			return false, nil
		default:
			return false, err

		}
	}
	return true, nil
}
