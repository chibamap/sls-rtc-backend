package connection

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConnection(t *testing.T) {

	ConnectionTableName = "sls_rtc_connections"
	tm := newTable(ConnectionTableName)
	cm, err := NewManager()

	if err != nil {
		t.Fatal("Error failed to initialize connection manager")
	}
	testConnectionID := "test-new-conn1"

	t.Run("Successful Connected", func(t *testing.T) {

		created, err := cm.NewConnection(testConnectionID)
		assert.Nil(t, err)

		fetched, err := tm.GetConn(testConnectionID)
		assert.Nil(t, err)
		assert.Equal(t, created.PK, fetched.PK)
	})

	tm.delete(pkPrefixRoom + testConnectionID)
}

func TestDisconnected(t *testing.T) {

	ConnectionTableName = "sls_rtc_connections"
	table := newTable(ConnectionTableName)
	cm, err := NewManager()

	if err != nil {
		t.Fatal("Error failed to initialize connection manager")
	}
	testConnectionID := "test-disconnected-conn1"
	testConnection := New(testConnectionID)
	if err = table.Put(testConnection); err != nil {
		t.Fatal("Error failed to create test record")
	}

	t.Run("Successful Disconnected", func(t *testing.T) {

		err := cm.Disconnected(testConnectionID)
		assert.Nil(t, err)

		pk := pkPrefixConn + testConnectionID
		fetched, err := table.consistentRead(pk)
		assert.Nil(t, err)
		assert.Nil(t, fetched)
	})

	testRoomID := "testroom1"
	testConnection.RoomID = testRoomID
	if err = table.Put(testConnection); err != nil {
		t.Fatal("Error failed to create test record")
	}

	room := newRoom(testRoomID)
	if err := table.PutRecord(room); err != nil {
		t.Fatal("failed to create test romm record")
	}

	t.Run("Successful delete connection with room", func(t *testing.T) {

		err := cm.Disconnected(testConnectionID)
		assert.Nil(t, err)

		pk := pkPrefixConn + testConnectionID
		fetchedConn, _ := table.consistentRead(pk)
		assert.Nil(t, fetchedConn)
		pk = pkPrefixRoom + testRoomID
		fetchedRoom, _ := table.GetRoom(pk)
		assert.Nil(t, fetchedRoom)
	})

}

func TestCreateRoom(t *testing.T) {

	ConnectionTableName = "sls_rtc_connections"
	table := newTable(ConnectionTableName)
	cm, err := NewManager()

	if err != nil {
		t.Fatal("Error failed to initialize connection manager")
	}
	testConnectionID := "test-newroom-conn1"
	testConnection := New(testConnectionID)
	if err = table.Put(testConnection); err != nil {
		t.Fatal("Error failed to create test record")
	}

	testRoomID := "test-newroom-id1"

	t.Run("Successful Created Room", func(t *testing.T) {
		success, err := cm.CreateRoom(testRoomID, testConnectionID)
		assert.Nil(t, err)
		assert.True(t, success)

		pk := pkPrefixConn + testConnectionID
		fetchedConn, err := table.consistentRead(pk)
		assert.Equal(t, testRoomID, fetchedConn.RoomID)
		pk = pkPrefixRoom + testRoomID
		roomRec, err := table.consistentRead(pk)
		assert.Equal(t, testRoomID, roomRec.RoomID)
	})

	testConnectionID2 := "test-newroom-conn2"
	testConnection2 := New(testConnectionID2)
	if err = table.Put(testConnection2); err != nil {
		t.Fatal("Error failed to create test record")
	}

	t.Run("Failed to create dupplicated Room", func(t *testing.T) {
		success, _ := cm.CreateRoom(testRoomID, testConnectionID2)
		assert.False(t, success)

		pk := pkPrefixConn + testConnectionID2
		fetchedConn, _ := table.consistentRead(pk)
		assert.NotEqual(t, testRoomID, fetchedConn.RoomID)
		pk = pkPrefixRoom + testRoomID
		roomRec, _ := table.consistentRead(pk)
		assert.NotNil(t, roomRec)
	})

	table.DeleteRoom(testRoomID)
	table.DeleteConnection(testConnectionID)
	table.DeleteConnection(testConnectionID2)

}
