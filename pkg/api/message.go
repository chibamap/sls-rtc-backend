package api

import (
	"encoding/json"

	"github.com/hogehoge-banana/sls-rtc-backend/internal/apigw"
	"github.com/hogehoge-banana/sls-rtc-backend/internal/ddb"

	"github.com/aws/aws-lambda-go/events"
)

type messageFrame struct {
	Message string `json:"message"`
}

// OnMessage send message to room mate
func OnMessage(req events.APIGatewayWebsocketProxyRequest) (string, error) {
	table, err := ddb.NewConnectionTable()
	if err != nil {
		return "", err
	}
	connections, err := table.ScanAll()
	if err != nil {
		return "", err
	}
	apigw, err := apigw.New(req.RequestContext)
	if err != nil {
		return "", err
	}
	var messageFrame messageFrame
	if err := json.Unmarshal([]byte(req.Body), &messageFrame); err != nil {
		return "", err
	}

	data := []byte(messageFrame.Message)
	if err := apigw.Send(data, connections); err != nil {
		return "", err
	}
	return "sent", nil
}
