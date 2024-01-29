package main

import (
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/hogehoge-banana/sls-rtc-backend/pkg/service/connection"
)

type proxyResponse events.APIGatewayProxyResponse

// CreateRoomIF api class
type enterRoomIF struct {
	RoomID string `json:"roomID"`
}

func main() {
	lambda.Start(handler)
}

func handler(request events.APIGatewayWebsocketProxyRequest) (proxyResponse, error) {

	msg, err := enterRoom(request)
	if err != nil {
		log.Println(msg)
		return proxyResponse{}, err
	}

	return proxyResponse{
		Body:       msg,
		StatusCode: 200,
	}, nil
}

// EnterRoom endpoint handler
func enterRoom(req events.APIGatewayWebsocketProxyRequest) (string, error) {
	params := enterRoomIF{}
	if err := json.Unmarshal([]byte(req.Body), &params); err != nil {
		return "un expected parameter given", err
	}

	m, err := connection.NewManager()
	if err != nil {
		return "could not initialize connection manager", err
	}

	// enter room
	if err = m.EnterRoom(req.RequestContext.ConnectionID, params.RoomID); err != nil {
		return "failed to enter room", err
	}

	return "ok", nil
}
