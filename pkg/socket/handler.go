package socket

import (
	"github.com/hogehoge-banana/sls-rtc-backend/internal/connection"
)

// OnConnected event handling on connected websocket
func OnConnected(connectionID string) (string, error) {
	cm, err := connection.NewConnectionManager(connectionID)

	if err := conn.StoreConnection(); err != nil {
		return "failed to connnect dynamodb", err
	}

	if err := cm.Store(); err != nil {
		return "failed to put item", err
	}
	return "ok", nil
}

// OnDisconnected event handling on disconnected
func OnDisconnected(connectionID string) (string, error) {
	cm, err := connection.NewConnectionManager(connectionID)

	if err != nil {
		return "", err
	}
	if err := cm.Delete(); err != nil {
		return "", err
	}
	return "bye", nil
}
