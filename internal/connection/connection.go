package connection

import (
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// Connection dynamodb record structure
type Connection struct {
	// ConnectionID request.RequestContext.ConnectionID
	PK           string `json:"pk"`
	ConnectionID string `json:"connectionID"`
	RoomID       string `json:"roomID"`
	Username     string `json:"username"`
}

// Room dynamodb record structure
type Room struct {
	// ConnectionID request.RequestContext.ConnectionID
	PK           string `json:"pk"`
	ConnectionID string `json:"connectionID"`
	RoomID       string `json:"roomID"`
}

// NewConnection new connection record
func NewConnection(connectionID string) *Connection {
	pk := pkPrefixConn + connectionID
	return &Connection{
		PK:           pk,
		ConnectionID: connectionID,
		RoomID:       blankValue,
	}
}

// NewRoom return table record struct
func NewRoom(roomID string) *Room {
	pk := pkPrefixRoom + roomID
	return &Room{
		PK:           pk,
		ConnectionID: blankValue,
		RoomID:       roomID,
	}
}

// Manager manage connection
type Manager struct {
	table *table
}

// ConnectionTableName the name of dynamodb table
var ConnectionTableName string

func init() {
	ConnectionTableName = os.Getenv("TABLE_NAME")
}

// NewManager returns connection manager instance
func NewManager() (*Manager, error) {

	if ConnectionTableName == "" {
		return nil, errors.New("tabne name was not set")
	}

	table := newTable(ConnectionTableName)

	return &Manager{
		table}, nil
}

// NewConnection make new connection and store to table
func (m *Manager) NewConnection(connectionID string) error {
	conn := NewConnection(connectionID)

	return m.table.PutConnection(conn)
}

// Disconnected cleanup records beside connection
func (m *Manager) Disconnected(connectionID string) error {
	return m.table.DeleteConnection(connectionID)
}

// FindConnection find out connection record from table
func (m *Manager) FindConnection(connectionID string) (*Connection, error) {
	return m.table.GetConn(connectionID)
}

// RetrieveRoomConnections retrieve connections at same room
func (m *Manager) RetrieveRoomConnections(roomID string) ([]*Connection, error) {
	return nil, nil
}

// CreateRoom create room
func (m *Manager) CreateRoom(roomID string) (bool, error) {
	room := NewRoom(roomID)
	err := m.table.PutNewRoom(room)
	if nil == err {
		return true, nil
	}
	if aerr, ok := err.(awserr.Error); ok {
		if aerr.Code() == dynamodb.ErrCodeConditionalCheckFailedException {
			// retryable error
			return false, nil
		}
	}
	return false, err
}

// EnterRoom update roomid
func (m *Manager) EnterRoom(connID string, roomID string) error {
	if room, _ := m.table.GetRoom(roomID); room == nil {
		return fmt.Errorf("room[%s] not found", roomID)
	}

	return m.table.UpdateConnectionRoomID(connID, roomID)
}

// FindRoomMates find all room mates
func (m *Manager) FindRoomMates(roomID string) ([]Connection, error) {
	keyCondition := "roomID = :roomID"
	filter := "connectionID <> :blank"
	expressionAttributeValues := map[string]*dynamodb.AttributeValue{
		":roomID": {
			S: aws.String(roomID),
		},
		":blank": {
			S: aws.String(blankValue),
		},
	}

	return m.table.QueryRoomMates(keyCondition, expressionAttributeValues, &filter)
	//return m.table.QueryRoomMates(keyCondition, expressionAttributeValues, nil)
}
