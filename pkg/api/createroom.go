package api

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/hogehoge-banana/sls-rtc-backend/internal/connection"
)

// CreateRoom api class
type CreateRoom struct {
	Connection *connection.Connection
}

// OnCreateMessage endpoint handler
func OnCreateMessage(req events.APIGatewayWebsocketProxyRequest) (string, error) {
	conn := connection.New(req.RequestContext.ConnectionID)

	return "ok", nil
}
