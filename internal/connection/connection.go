package connection

// Connection webwocket management
type Connection struct {
	// ConnectionID request.RequestContext.ConnectionID
	PK           string `json:"pk"`
	ConnectionID string `json:"connectionID"`
	RoomID       string `json:"roomID"`
	Owner        string
}

// ConnectionManager manage connection
type ConnectionManager struct {
	Conn  *Connection
	table *Table
}

// New return Connection pointer
func New(connectionID string) *Connection {
	return &Connection{
		ConnectionID: connectionID}
}

// NewConnectionManager returns connection manager instance
func NewConnectionManager(connectionID string) (*ConnectionManager, error) {

	conn := New(connectionID)
	table, err := NewTable()
	if err != nil {
		return nil, err
	}

	return &ConnectionManager{
		conn, table}, nil
}

// Store connection
func (cm *ConnectionManager) Store() error {
	cm.Conn.PK = cm.Conn.ConnectionID
	return cm.table.Put(cm.Conn)
}

// Delete connection
func (cm *ConnectionManager) Delete() error {
	cm.Conn.PK = cm.Conn.ConnectionID
	return cm.table.Delete(cm.Conn)
}
