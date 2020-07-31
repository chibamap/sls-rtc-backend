package message

import (
	"encoding/json"

	"github.com/hogehoge-banana/sls-rtc-backend/internal/apigw"
	"github.com/hogehoge-banana/sls-rtc-backend/internal/connection"

	"github.com/aws/aws-lambda-go/events"
)

type messageInterface struct {
	Message string `json:"message"`
}

// OnMessage send message to room mate
func OnMessage(req events.APIGatewayWebsocketProxyRequest) (string, error) {
	cm, err := connection.NewManager()
	if err != nil {
		return "", err
	}
	conn, err := cm.FindConnection(req.RequestContext.ConnectionID)

	connections, err := cm.RetrieveRoomConnections(conn.RoomID)
	if err != nil {
		return "", err
	}

	var messageFrame messageInterface
	if err := json.Unmarshal([]byte(req.Body), &messageFrame); err != nil {
		return "", err
	}
	transferMessage := &apigw.MessageFrame{
		Body: messageFrame.Message,
	}

	apigw, err := apigw.New(req.RequestContext)
	if err != nil {
		return "", err
	}

	if err := apigw.Multicast(transferMessage, connections); err != nil {
		return "", err
	}
	return "sent", nil
}
