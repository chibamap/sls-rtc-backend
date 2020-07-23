package main

import (
	"log"
	"sls-rtc-backend/internal/connection"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type proxyRequest events.APIGatewayWebsocketProxyRequest
type proxyResponse events.APIGatewayProxyResponse

func handler(request proxyRequest) (proxyResponse, error) {
	log.Println("disconnected request")

	if _, err := connection.OnDisconnected(); err != nil {
		return proxyResponse{}, err
	}

	return proxyResponse{
		Body:       "disconnected",
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
