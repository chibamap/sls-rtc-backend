package apigw

import (
	"encoding/json"
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
	ctx    events.APIGatewayWebsocketProxyRequestContext
}

// MessageFrame message frame for transfer
type MessageFrame struct {
	Type string `json:"type"`
	Body string `json:"body"`
}

const (
	// TypeEnter enter room event
	TypeEnter = "enter"
	// TypeMessage message event
	TypeMessage = "message"
)

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
		ctx,
	}, nil
}

// Multicast data to connections
func (a *Apigw) Multicast(message *MessageFrame, conns []*connection.Connection) error {
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}
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

// Respond to current connection
func (a *Apigw) Respond(message *MessageFrame) error {
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}
	postInput := &apigatewaymanagementapi.PostToConnectionInput{
		Data:         data,
		ConnectionId: &a.ctx.ConnectionID,
	}

	_, err = a.client.PostToConnection(postInput)
	return err
}
