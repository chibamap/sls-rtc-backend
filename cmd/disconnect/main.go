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
	log.Println("disconnected request")

	msg, err := onDisconnected(request.RequestContext.ConnectionID)
	if err != nil {
		return proxyResponse{}, err
	}

	return proxyResponse{
		Body:       msg,
		StatusCode: 200,
	}, nil
}

// onDisconnected event handling on disconnected
func onDisconnected(connectionID string) (string, error) {
	cm, err := connection.NewManager()

	if err != nil {
		return "failed to initialize manager", err
	}

	if err := cm.Disconnected(connectionID); err != nil {
		return "", err
	}
	return "bye", nil
}
