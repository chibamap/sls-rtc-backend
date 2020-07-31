package connection

import (
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

const pkPrefix = "connectionID#"
const roomRecPrefix = "roomID#"
const initialRoom = "-"

// Connection webwocket management
type Connection struct {
	// ConnectionID request.RequestContext.ConnectionID
	PK           string `json:"pk"`
	ConnectionID string `json:"connectionID"`
	RoomID       string `json:"roomID"`
	Username     string `json:"username"`
	Owner        string
}

// New return Connection pointer
func New(connectionID string) *Connection {
	pk := pkPrefix + connectionID
	return &Connection{
		PK:           pk,
		ConnectionID: connectionID,
		RoomID:       initialRoom}
}

// Manager manage connection
type Manager struct {
	table *table
}

// NewManager returns connection manager instance
func NewManager() (*Manager, error) {

	table, err := newTable()
	if err != nil {
		return nil, err
	}

	return &Manager{
		table}, nil
}

// NewConnection make new connection and store to table
func (m *Manager) NewConnection(connectionID string) (*Connection, error) {
	conn := New(connectionID)

	err := m.table.Put(conn)
	return conn, err
}

// Disconnected cleanup records beside connection
func (m *Manager) Disconnected(connectionID string) error {
	conn := New(connectionID)
	return m.table.Delete(conn)
}

// FindConnection find out connection record from table
func (m *Manager) FindConnection(connectionID string) (*Connection, error) {
	return m.table.Get(connectionID)
}

// RetrieveRoomConnections retrieve connections at same room
func (m *Manager) RetrieveRoomConnections(roomID string) ([]*Connection, error) {
	return nil, nil
}

// NewRoom create room
func (m *Manager) NewRoom(roomID string, ownerConn *Connection) (bool, error) {

	room := &Connection{
		PK:           roomRecPrefix + roomID,
		ConnectionID: ownerConn.ConnectionID,
		RoomID:       roomID,
	}

	ownerConn.RoomID = roomID
	ownerKey, _ := dynamodbattribute.MarshalMap(ownerConn)
	roomKey, _ := dynamodbattribute.MarshalMap(room)
	items := [] *dynamodb.TransactWriteItem{
	&dynamodb.TransactWriteItem{
		Update: &dynamodb.Update{
			Key: ownerKey,
			UpdateExpression: "SET roomID = :roomID",
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				":roomID": { S: aws.String(roomID) },
			},
		},
	},
	&dynamodb.TransactWriteItem{
		Put: &dynamodb.Put{
Item: room
			}
		}
	}
}
	err := m.table.TransactPut([]*Connection{
		room,
		ownerConn,
	})

	return false, err
}
