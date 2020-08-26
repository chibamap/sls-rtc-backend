package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/hogehoge-banana/sls-rtc-backend/internal/socket"
	"github.com/hogehoge-banana/sls-rtc-backend/pkg/api/roommate"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

const (
	blankValue = "-"
)

var (
	// ApigatewayEndpoint ex. Prod
	apigatewayEndpoint string
	// socket
	websocket *socket.Socket
)

func init() {
	apigatewayEndpoint = os.Getenv("APIGW_ENDPOINT")
}

func main() {
	lambda.Start(handleRequest)
}

// Connection dynamodb record structure
type tableRecord struct {
	// ConnectionID request.RequestContext.ConnectionID
	PK           string `json:"pk"`
	ConnectionID string `json:"connectionID"`
	RoomID       string `json:"roomID"`
}

func handleRequest(ctx context.Context, e events.DynamoDBEvent) {
	websocket, _ = socket.New(apigatewayEndpoint)

	for _, record := range e.Records {
		fmt.Printf("Processing request data for event ID %s, type %s.\n", record.EventID, record.EventName)

		switch record.EventName {
		case "INSERT":
			handleInsertStream(record.Change)
		case "MODIFY":
			handleModifyStream(record.Change)
		case "REMOVE":
			handleRemoveStrem(record.Change)
		}
	}
}

func handleInsertStream(streamRecord events.DynamoDBStreamRecord) {
	fmt.Printf("new image: %v\n", streamRecord.NewImage)
	v, ok := streamRecord.NewImage["connectionID"]
	if !ok {
		log.Println("connection id not set")
		return
	}
	connectionID := v.String()
	if blankValue == connectionID {
		log.Println("connection id was blank")
		return
	}

	message := &socket.MessageFrame{
		Type: socket.TypeConnected,
		Data: connectionID,
	}
	if err := websocket.SendMessage(connectionID, message); err != nil {
		fmt.Printf("failed to send message %v \n", err)
	}
}

func handleModifyStream(streamRecord events.DynamoDBStreamRecord) {
	v, ok := streamRecord.NewImage["connectionID"]
	if !ok {
		return
	}
	connectionID := v.String()
	if blankValue == connectionID {
		return
	}

	roomID := streamRecord.NewImage["roomID"].String()
	oldRoomID := streamRecord.OldImage["roomID"].String()

	if roomID == oldRoomID {
		return
	}

	if roomID != blankValue {
		fmt.Printf("new room mate. connection ID %s, room ID %s \n", connectionID, roomID)
		if err := roommate.NewRoomMate(websocket, connectionID, roomID); err != nil {
			fmt.Printf("failed to handle new room mate %v\n", err)
		}
	}

	if oldRoomID != blankValue {
		fmt.Printf("left room mate. connection ID %s, room ID %s \n", connectionID, roomID)
		if err := roommate.LeaveRoomMate(websocket, connectionID, roomID); err != nil {
			fmt.Printf("failed to send message %v \n", err)
		}
	}
}

func handleRemoveStrem(streamRecord events.DynamoDBStreamRecord) {
	connectionID := streamRecord.OldImage["connectionID"].String()
	if blankValue == connectionID {
		return
	}
	roomID := streamRecord.OldImage["roomID"].String()

	if blankValue == roomID {
		return
	}

	fmt.Printf("left room mate. connection ID %s, room ID %s \n", connectionID, roomID)
	roommate.LeaveRoomMate(websocket, connectionID, roomID)
}
