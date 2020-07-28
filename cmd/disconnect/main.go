package main

import (
	"log"

	"github.com/hogehoge-banana/sls-rtc-backend/pkg/socket"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type proxyRequest events.APIGatewayWebsocketProxyRequest
type proxyResponse events.APIGatewayProxyResponse

func handler(request proxyRequest) (proxyResponse, error) {
	log.Println("disconnected request")

	msg, err := socket.OnDisconnected(request.RequestContext.ConnectionID)
	if err != nil {
		return proxyResponse{}, err
	}

	return proxyResponse{
		Body:       msg,
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
