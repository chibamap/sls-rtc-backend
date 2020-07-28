package apigw

import (
	"fmt"
	"log"

	"github.com/hogehoge-banana/sls-rtc-backend/internal/connection"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi"
)

// Apigw apigateway client wrapper
type Apigw struct {
	// TableName which use dynamodb table name
	client *apigatewaymanagementapi.ApiGatewayManagementApi
}

// New make dynamodb session
func New(ctx events.APIGatewayWebsocketProxyRequestContext) (*Apigw, error) {
	session, err := session.NewSession()
	if err != nil {
		return nil, err
	}

	client := apigatewaymanagementapi.New(session)

	client.Endpoint = fmt.Sprintf("https://%s/%s", ctx.DomainName, ctx.Stage)

	return &Apigw{
		client,
	}, nil
}

// Send data to connections
func (a *Apigw) Send(data []byte, conns []connection.Connection) error {
	postInput := &apigatewaymanagementapi.PostToConnectionInput{
		Data: data,
	}

	for _, conn := range conns {
		postInput.ConnectionId = &conn.ConnectionID
		if _, err := a.client.PostToConnection(postInput); err != nil {
			log.Println(err)
		}
	}
	return nil
}
