package connection

// Connection webwocket management
type Connection struct {
	// ConnectionID request.RequestContext.ConnectionID
	PK           string `json:"pk"`
	ConnectionID string `json:"connectionID"`
	RoomID       string `json:"roomID"`
	Username     string `json:"username"`
	Owner        string
}

const pkPrefix = "connectionID:"
const roomRecPrefix = "roomID:"

// Manager manage connection
type Manager struct {
	table *table
}

// New return Connection pointer
func New(connectionID string) *Connection {
	pk := pkPrefix + connectionID
	return &Connection{
		PK:           pk,
		ConnectionID: connectionID}
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
