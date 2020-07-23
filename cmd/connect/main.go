package main

import (
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"sls-rtc-backend/internal/connection"
)

type proxyRequest events.APIGatewayWebsocketProxyRequest
type proxyResponse events.APIGatewayProxyResponse

func handler(request proxyRequest) (proxyResponse, error) {
	log.Println("connected request")

	if _, err := connection.OnConnected(); err != nil {
		return proxyResponse{}, err
	}

	return proxyResponse{
		Body:       "connected",
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
