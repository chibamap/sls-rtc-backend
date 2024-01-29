package main

import (
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/hogehoge-banana/sls-rtc-backend/pkg/service/connection"
)

type proxyRequest events.APIGatewayWebsocketProxyRequest
type proxyResponse events.APIGatewayProxyResponse

func main() {
	lambda.Start(handler)
}

func handler(request proxyRequest) (proxyResponse, error) {
	connectionID := request.RequestContext.ConnectionID
	msg, err := onConnected(connectionID)
	if err != nil {
		log.Println(msg)
		return proxyResponse{}, err
	}

	return proxyResponse{
		Body:       msg,
		StatusCode: 200,
	}, nil
}

// onConnected event handling on connected websocket
func onConnected(connectionID string) (string, error) {

	cm, err := connection.NewManager()
	if err != nil {
		return "failed to initialize manager", err
	}

	if err = cm.NewConnection(connectionID); err != nil {
		return "failed to connnect dynamodb", err
	}

	return "ok", nil
}
