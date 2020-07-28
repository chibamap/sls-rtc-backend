package main

import (
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/hogehoge-banana/sls-rtc-backend/pkg/api"
)

type proxyResponse events.APIGatewayProxyResponse

func handler(request events.APIGatewayWebsocketProxyRequest) (proxyResponse, error) {
	log.Println("message handler main")
	if _, err := api.OnMessage(request); err != nil {
		return proxyResponse{}, err
	}

	return proxyResponse{
		Body:       fmt.Sprintf("ok"),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
