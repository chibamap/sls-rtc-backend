package connection

import (
	"testing"

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
		pk := pkPrefixRoom + "test"

		conn, err := table.find(pk)
		assert.Nil(t, err)
		assert.NotNil(t, conn)
		assert.Equal(t, conn.PK, pk)
	})
}

// Put connection item to dynamo db
func TestPut(t *testing.T) {
	tableName := "sls_rtc_connections"
	table := newTable(tableName)

	t.Run("Successful put new test record", func(t *testing.T) {
		testConnID2 := "test2"
		testConn := New(testConnID2)

		err := table.Put(testConn)

		assert.Nil(t, err)
		fetched, _ := table.GetConn(testConnID2)
		assert.Equal(t, testConnID2, fetched.PK)
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

// TransactPut
func TestPutNewRoom(t *testing.T) {
	tableName := "sls_rtc_connections"
	testRoomID := "testroom2"
	testConnID := "roomTestConnID1"
	testConnID2 := "roomTestConnID2"

	t.Run("Successful create room", func(t *testing.T) {
		table := newTable(tableName)

		testConn := New(testConnID)
		if err := table.Put(testConn); err != nil {
			t.Fatal("could not create test record")
		}

		room := newRoom(testRoomID)

		success, err := table.PutNewRoom(room, testConnID)
		assert.True(t, success)
		assert.Nil(t, err)
		createdRoom, err := table.GetRoom(testRoomID)

		assert.Nil(t, err)
		assert.NotNil(t, createdRoom)
		assert.Equal(t, room.PK, createdRoom.PK)
	})

	t.Run("Faile to create duplicated room", func(t *testing.T) {
		table := newTable(tableName)
		testConn2 := New(testConnID2)
		if err := table.Put(testConn2); err != nil {
			t.Fatal("could not create test record2")
		}

		room := newRoom(testRoomID)

		success, err := table.PutNewRoom(room, testConnID2)
		assert.False(t, success)
		assert.Nil(t, err)
	})

	table := newTable(tableName)

	table.delete(pkPrefixConn + testConnID)
	table.delete(pkPrefixConn + testConnID2)
	table.delete(pkPrefixRoom + testRoomID)
}
