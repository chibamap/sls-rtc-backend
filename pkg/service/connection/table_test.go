package connection

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/stretchr/testify/assert"
)

// Newtable instance from table name
func TestNewTable(t *testing.T) {
	tableName := "test_table"
	t.Run("Successful new table instance", func(t *testing.T) {

		table := newTable(tableName)

		assert.NotNil(t, table)
		assert.Equal(t, table.tableName, tableName)
	})
}

func TestGetConn(t *testing.T) {
	tableName := "sls_rtc_connections"
	table := newTable(tableName)

	t.Run("Successful retrieve test record", func(t *testing.T) {
		testConnID := "test"

		conn, err := table.GetConn(testConnID)
		assert.Nil(t, err)
		assert.NotNil(t, conn)
		assert.Equal(t, conn.ConnectionID, testConnID)
	})
}

func TestGetRoom(t *testing.T) {
	tableName := "sls_rtc_connections"
	table := newTable(tableName)

	t.Run("Successful retrieve room record", func(t *testing.T) {
		testRoomID := "testroom"

		rec, err := table.GetRoom(testRoomID)
		assert.Nil(t, err)
		assert.NotNil(t, rec)
		assert.Equal(t, rec.RoomID, testRoomID)
	})
}

func TestFind(t *testing.T) {
	tableName := "sls_rtc_connections"
	table := newTable(tableName)

	t.Run("Successful retrieve record by pk", func(t *testing.T) {
		pk := pkPrefixConn + "test"

		item, err := table.find(pk)
		assert.Nil(t, err)
		assert.NotNil(t, item)
		room := Room{}
		err = dynamodbattribute.UnmarshalMap(item, &room)
		assert.Nil(t, err)
		assert.Equal(t, pk, room.PK)
	})
}

// Put connection item to dynamo db
func TestPutNewConnection(t *testing.T) {
	tableName := "sls_rtc_connections"
	table := newTable(tableName)

	t.Run("Successful put new test record", func(t *testing.T) {
		testConnID2 := "test2"
		testConn := NewConnection(testConnID2)

		err := table.PutNewConnection(testConn)

		assert.Nil(t, err)
		fetched, _ := table.GetConn(testConnID2)
		assert.Equal(t, testConnID2, fetched.ConnectionID)
	})
}

// Delete connection item from dynamo db
func TestDelete(t *testing.T) {

	t.Run("Successful delete test record", func(t *testing.T) {
		tableName := "sls_rtc_connections"
		table := newTable(tableName)
		pk := pkPrefixRoom + "test2"

		err := table.delete(pk)
		assert.Nil(t, err)
		fetched, _ := table.find(pk)
		assert.Nil(t, fetched)
	})
}
