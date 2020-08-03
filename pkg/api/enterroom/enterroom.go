package enterroom

import (
	"encoding/json"

	"github.com/hogehoge-banana/sls-rtc-backend/internal/apigw"

	"github.com/aws/aws-lambda-go/events"
	"github.com/google/uuid"
	"github.com/hogehoge-banana/sls-rtc-backend/internal/connection"
)

// CreateRoomIF api class
type enterRoomIF struct {
	RoomID string `json:"roomID"`
}

const maxTry = 5

// EnterRoom endpoint handler
func EnterRoom(req events.APIGatewayWebsocketProxyRequest) (string, error) {
	params := enterRoomIF{}
	if err := json.Unmarshal([]byte(req.Body), &params); err != nil {
		return "un expected parameter given", err
	}

	m, err := connection.NewManager()
	if err != nil {
		return "could not initialize connection manager", err
	}

	// enter room
	if err = m.EnterRoom(req.RequestContext.ConnectionID, params.RoomID); err != nil {
		// TODO: notify failed to client
		return "failed to enter room", err
	}

	// respond to
	gw, err := apigw.New(req.RequestContext)
	if err != nil {
		return "failed initialize apigateway client", err
	}
	gw.RespondRoomEntered(params.RoomID)

	return "ok", nil
}

func generateRoomID() string {
	return uuid.New().String()
}
