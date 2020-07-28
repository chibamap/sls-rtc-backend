package connection

// Connection webwocket management
type Connection struct {
	// ConnectionID request.RequestContext.ConnectionID
	ConnectionID string `json:"connectionId"`
}

// New return Connection pointer
func New(connectionID string) *Connection {
	return &Connection{connectionID}
}
