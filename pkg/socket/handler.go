package socket

import (
	"github.com/hogehoge-banana/sls-rtc-backend/internal/ddb"

	"github.com/hogehoge-banana/sls-rtc-backend/internal/connection"
)

// OnConnected event handling on connected websocket
func OnConnected(connectionID string) (string, error) {
	conn := connection.New(connectionID)
	table, err := ddb.NewConnectionTable()
	if err != nil {
		return "failed to connnect dynamodb", err
	}

	if err := table.Put(conn); err != nil {
		return "failed to put item", err
	}
	return "ok", nil
}

// OnDisconnected event handling on disconnected
func OnDisconnected(connectionID string) (string, error) {
	conn := connection.New(connectionID)
	table, err := ddb.NewConnectionTable()
	if err != nil {
		return "", err
	}
	if err := table.Delete(conn); err != nil {
		return "", err
	}
	return "bye", nil
}
