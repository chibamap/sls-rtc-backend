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

		err := cm.NewConnection(testConnectionID)
		assert.Nil(t, err)

		fetched, err := tm.GetConn(testConnectionID)
		assert.Nil(t, err)
		assert.Equal(t, testConnectionID, fetched.ConnectionID)
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
	testConnection := NewConnection(testConnectionID)
	if err = table.PutConnection(testConnection); err != nil {
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
	if err = table.PutNewConnection(testConnection); err != nil {
		t.Fatal("Error failed to create test record")
	}

	room := NewRoom(testRoomID)
	if err := table.PutNewRoom(room); err != nil {
		t.Fatal("failed to create test romm record")
	}

	t.Run("Successful delete connection with empty room", func(t *testing.T) {

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

	testRoomID := "test-newroom-id1"

	t.Run("Successful Created Room", func(t *testing.T) {
		success, err := cm.CreateRoom(testRoomID)
		assert.Nil(t, err)
		assert.True(t, success)

		pk := pkPrefixRoom + testRoomID
		roomRec, err := table.consistentRead(pk)
		assert.Equal(t, testRoomID, roomRec.RoomID)
	})

	t.Run("Failed to create dupplicated Room", func(t *testing.T) {
		success, _ := cm.CreateRoom(testRoomID)
		assert.False(t, success)

		pk := pkPrefixRoom + testRoomID
		roomRec, _ := table.consistentRead(pk)
		assert.NotNil(t, roomRec)
	})

	table.DeleteRoom(testRoomID)
}

func TestEnterRoom(t *testing.T) {

	ConnectionTableName = "sls_rtc_connections"
	table := newTable(ConnectionTableName)
	cm, err := NewManager()

	if err != nil {
		t.Fatal("Error failed to initialize connection manager")
	}
	testConnectionID := "test-enterroom-conn1"
	testConnection := NewConnection(testConnectionID)
	if err = table.PutNewConnection(testConnection); err != nil {
		t.Fatal("Error failed to create test record")
	}
	testRoomID := "enter-room-test"
	room := NewRoom(testRoomID)
	if err := table.PutNewRoom(room); err != nil {
		t.Fatal("Error failed to create test record")
	}

	t.Run("Successful enter room", func(t *testing.T) {
		beforeConn, _ := table.GetConn(testConnectionID)
		assert.Equal(t, blankValue, beforeConn.RoomID)

		err := cm.EnterRoom(testConnectionID, testRoomID)
		assert.Nil(t, err)

		afterConn, _ := table.consistentRead(beforeConn.PK)
		assert.Equal(t, testRoomID, afterConn.RoomID)
	})
	testConnectionID2 := "test-enterroom-conn2"
	testConnection2 := NewConnection(testConnectionID2)
	if err = table.PutNewConnection(testConnection2); err != nil {
		t.Fatal("Error failed to create test record")
	}

	t.Run("Fail to enter not found room", func(t *testing.T) {
		fakeRoomID := "enter-room-test-fake"
		beforeConn, _ := table.GetConn(testConnectionID2)
		assert.Equal(t, blankValue, beforeConn.RoomID)

		err := cm.EnterRoom(testConnectionID, fakeRoomID)
		assert.NotNil(t, err)

		afterConn, _ := table.consistentRead(beforeConn.PK)
		assert.Equal(t, blankValue, afterConn.RoomID)
	})
	table.DeleteConnection(testConnectionID)
	table.DeleteConnection(testConnectionID2)
	table.DeleteRoom(testRoomID)

}
