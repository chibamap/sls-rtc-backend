package connection

type Connection struct {}

func New() *Connection {
	return &Connection{}
}

func OnConnected() (*Connection, error) {
	return nil, nil
}

func OnDisconnected() (*Connection, error) {
	return nil, nil
}
