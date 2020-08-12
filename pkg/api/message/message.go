package message

import (
	"errors"
	"os"

	"github.com/hogehoge-banana/sls-rtc-backend/internal/connection"
	"github.com/hogehoge-banana/sls-rtc-backend/internal/socket"

	"github.com/aws/aws-lambda-go/events"
)

// MessageInterface api parameter interface
type MessageInterface struct {
	Dest    string `json:"dest"`
	Type    string `json:"type"`
	Message string `json:"message"`
}

// Handler handle message
type Handler struct {
	ctx events.APIGatewayWebsocketProxyRequestContext
}

var (
	// ApigatewayEndpoint ex. Prod
	apigatewayEndpoint string
)

func init() {
	apigatewayEndpoint = os.Getenv("APIGW_ENDPOINT")
}

// NewHandler return handler
func NewHandler(ctx events.APIGatewayWebsocketProxyRequestContext) *Handler {
	return &Handler{ctx: ctx}
}

// OnMessage send message to room mate
func (h *Handler) OnMessage(param *MessageInterface) (string, error) {

	cm, err := connection.NewManager()
	if err != nil {
		return "couldn't initialize connection manager", err
	}
	src, err := cm.FindConnection(h.ctx.ConnectionID)
	if err != nil || src == nil {
		return "could not find source connection", err
	}

	dest, err := cm.FindConnection(param.Dest)
	if err != nil {
		return "could not find destination connection", err
	}
	if dest == nil || src.RoomID != dest.RoomID {
		msg := "could not find destination connection"
		err = errors.New(msg)
		return msg, err
	}

	message := &socket.MessageFrame{
		Type: param.Type,
		Data: param.Message,
		From: h.ctx.ConnectionID,
	}

	s, err := socket.New(apigatewayEndpoint)
	if err != nil {
		return "could not initialize socket", err
	}
	if err = s.SendMessage(dest.ConnectionID, message); err != nil {
		return "could not send message", err
	}

	return "sent", nil
}
