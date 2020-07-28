package main

import (
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/hogehoge-banana/sls-rtc-backend/pkg/socket"
)

type proxyRequest events.APIGatewayWebsocketProxyRequest
type proxyResponse events.APIGatewayProxyResponse

func handler(request proxyRequest) (proxyResponse, error) {
	log.Println("connected request")
	connectionID := request.RequestContext.ConnectionID
	msg, err := socket.OnConnected(connectionID)
	if err != nil {
		log.Println(msg)
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
