package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hogehoge-banana/sls-rtc-backend/pkg/service/message"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type proxyResponse events.APIGatewayProxyResponse

type messageInterface struct {
	Dest    string `json:"dest"`
	Type    string `json:"type"`
	Message string `json:"message"`
}

func main() {
	lambda.Start(handler)
}

func handler(request events.APIGatewayWebsocketProxyRequest) (proxyResponse, error) {
	if msg, err := onMessage(request); err != nil {
		log.Println(msg)
		return proxyResponse{}, err
	}

	return proxyResponse{
		Body:       fmt.Sprintf("ok"),
		StatusCode: 200,
	}, nil
}

// OnMessage send message to room mate
func onMessage(req events.APIGatewayWebsocketProxyRequest) (string, error) {

	var param message.MessageInterface
	if err := json.Unmarshal([]byte(req.Body), &param); err != nil {
		return fmt.Sprintf("invalid message. %s", req.Body), err
	}

	h := message.NewHandler(req.RequestContext)
	return h.OnMessage(&param)
}
