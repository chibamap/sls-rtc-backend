package createroom

import (
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/google/uuid"
	"github.com/hogehoge-banana/sls-rtc-backend/internal/connection"
	"github.com/hogehoge-banana/sls-rtc-backend/internal/socket"
)

const maxTry = 5

// CreateRoom endpoint handler
func CreateRoom(req events.APIGatewayWebsocketProxyRequest) (string, error) {

	cm, err := connection.NewManager()
	if err != nil {
		return "could not initialize connection manager", err
	}

	// loop until unique room has been created. up to 5 times. return error it retry more than 5 times
	try := 1
	var uid string
	var success bool
	for {
		if maxTry < try {
			log.Println("creating uuid achieved to max retry count")
			break
		}
		uid = generateRoomID()
		success, err = cm.CreateRoom(uid)
		if err != nil {
			return "failed to create room", err
		}
		if success {
			break
		}

		try++
	}

	s, err := socket.New(req.RequestContext.DomainName, req.RequestContext.Stage)
	if err != nil {
		return "failed to initialize apigateway client", err
	}

	if err := s.SendRoomCreated(req.RequestContext.ConnectionID, uid); err != nil {
		return "failed to respond", err
	}
	return "ok", nil
}

func generateRoomID() string {
	return uuid.New().String()
}
