package socket

import (
	"encoding/json"
	"log"

	"github.com/hogehoge-banana/sls-rtc-backend/internal/connection"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi"
)

const (
	// TypeConnected connected evetns
	TypeConnected = "connected"
	// TypeRoomCreated room created event. body is containing created room id
	TypeRoomCreated = "room-created"
	// TypeEnter enter room event.
	TypeEnter = "enter"
	// TypeLeave leave room event.
	TypeLeave = "leave"
	// TypeNewRoomMate new room mate joining event
	TypeNewRoomMate = "new-room-mate"
)

// Socket apigateway client wrapper
type Socket struct {
	// TableName which use dynamodb table name
	client *apigatewaymanagementapi.ApiGatewayManagementApi
}

// MessageFrame message frame for transfer
type MessageFrame struct {
	Type string `json:"type"`
	Data string `json:"body"`
	From string `json:"from"`
}

// EnterRoomMessageFrame ...
type EnterRoomMessageFrame struct {
	// should be enter
	Type         string `json:"type"`
	RoomID       string `json:"roomID"`
	ConnectionID string `json:"connectionID"`
}

// New make dynamodb session
func New(endpoint string) (*Socket, error) {
	session, err := session.NewSession()
	if err != nil {
		return nil, err
	}

	client := apigatewaymanagementapi.New(session)

	client.Endpoint = endpoint
	log.Printf("api gateway endpiont: %s", client.Endpoint)

	return &Socket{
		client,
	}, nil
}

// Multicast data to connections
func (s *Socket) Multicast(message *MessageFrame, conns []*connection.Connection) error {
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}
	postInput := &apigatewaymanagementapi.PostToConnectionInput{
		Data: data,
	}

	for _, conn := range conns {
		postInput.ConnectionId = &conn.ConnectionID
		if _, err := s.client.PostToConnection(postInput); err != nil {
			log.Println(err)
		}
	}
	return nil
}

// SendRoomCreated notify room created event
func (s *Socket) SendRoomCreated(connectionID, roomID string) error {
	// notify success
	message := &MessageFrame{
		Type: TypeRoomCreated,
		Data: roomID,
	}

	data, err := json.Marshal(message)
	if err != nil {
		return err
	}
	return s.send(connectionID, data)
}

// SendRoomEntered notify room created event
func (s *Socket) SendRoomEntered(destID, newConnID, roomID string) error {
	// notify success
	message := &EnterRoomMessageFrame{
		Type:         TypeEnter,
		RoomID:       roomID,
		ConnectionID: newConnID,
	}

	data, err := json.Marshal(message)
	if err != nil {
		return err
	}
	return s.send(destID, data)
}

// SendRoomLeave notify room leave
func (s *Socket) SendRoomLeave(destID, leftConnID, roomID string) error {
	// notify success
	message := &EnterRoomMessageFrame{
		Type:         TypeLeave,
		RoomID:       roomID,
		ConnectionID: leftConnID,
	}

	data, err := json.Marshal(message)
	if err != nil {
		return err
	}
	return s.send(destID, data)
}

// SendMessage send a message
func (s *Socket) SendMessage(connID string, message *MessageFrame) error {
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return s.send(connID, data)
}

func (s *Socket) send(connID string, data []byte) error {
	postInput := &apigatewaymanagementapi.PostToConnectionInput{
		Data:         data,
		ConnectionId: &connID,
	}

	_, err := s.client.PostToConnection(postInput)
	return err
}
