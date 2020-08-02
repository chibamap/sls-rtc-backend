package connection

import (
	"errors"
	"log"
	"os"
)

// Connection dynamodb record structure
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
	pk := pkPrefixConn + connectionID
	return &Connection{
		PK:           pk,
		ConnectionID: connectionID,
		RoomID:       blankValue}
}

// NewRoom return table record struct
func newRoom(roomID string) *tableRecord {
	pk := pkPrefixRoom + roomID
	return &tableRecord{
		PK:           pk,
		ConnectionID: blankValue,
		RoomID:       blankValue,
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
func (m *Manager) NewConnection(connectionID string) (*Connection, error) {
	conn := New(connectionID)

	err := m.table.Put(conn)
	return conn, err
}

// Disconnected cleanup records beside connection
func (m *Manager) Disconnected(connectionID string) error {
	conn, err := m.table.GetConn(connectionID)
	if err != nil {
		return err
	}

	if err = m.table.DeleteConnection(connectionID); err != nil {
		log.Println(err)
	}
	if err = m.table.DeleteRoom(conn.RoomID); err != nil {
		log.Println(err)
	}
	return err
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
func (m *Manager) CreateRoom(roomID string, ownerConnectionID string) (bool, error) {

	pk := pkPrefixRoom + roomID
	room := &Connection{
		PK:           pk,
		ConnectionID: blankValue,
		RoomID:       blankValue,
	}

	return m.table.PutNewRoom(room, ownerConnectionID)
}
