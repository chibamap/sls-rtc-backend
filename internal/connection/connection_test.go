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
	/*
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

		testRoomID := "test-new-room-id"

		t.Run("Successful Created Room", func(t *testing.T) {
			err := cm.NewRoom(testRoomID, testConnection)
			assert.Nil(t, err)

			fetchedConn, err := table.Get(testConnectionID)
			assert.Nil(t, err)
			assert.Nil(t, fetchedConn)
			assert.Equal(t, fetchedConn.RoomID, testRoomID)
			roomRec, err := table.Get(testConnectionID)

		})

		t.Run("Successful Created Room", func(t *testing.T) {

			err := cm.NewRoom(testRoomID, testConnection)
			assert.Nil(t, err)

			fetched, err := tm.Get(testConnectionID)
			assert.Nil(t, err)
			assert.Nil(t, fetched)
		})
	*/
}
