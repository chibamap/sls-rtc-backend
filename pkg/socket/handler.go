package socket

import (
	"github.com/hogehoge-banana/sls-rtc-backend/internal/connection"
)

// OnConnected event handling on connected websocket
func OnConnected(connectionID string) (string, error) {

	cm, err := connection.NewManager()
	if err != nil {
		return "failed to initialize manager", err
	}

	if err = cm.NewConnection(connectionID); err != nil {
		return "failed to connnect dynamodb", err
	}

	return "ok", nil
}

// OnDisconnected event handling on disconnected
func OnDisconnected(connectionID string) (string, error) {
	cm, err := connection.NewManager()

	if err != nil {
		return "failed to initialize manager", err
	}

	if err := cm.Disconnected(connectionID); err != nil {
		return "", err
	}
	return "bye", nil
}
