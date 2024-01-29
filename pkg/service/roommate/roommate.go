package roommate

import (
	"fmt"

	"github.com/hogehoge-banana/sls-rtc-backend/pkg/service/connection"
	"github.com/hogehoge-banana/sls-rtc-backend/pkg/service/socket"
)

// NewRoomMate ...
func NewRoomMate(socket *socket.Socket, connectionID string, roomID string) error {
	cm, err := connection.NewManager()
	if err != nil {
		return err
	}

	roomMates, err := cm.FindRoomMates(roomID)
	for _, roommate := range roomMates {
		if senderr := socket.SendRoomEntered(roommate.ConnectionID, connectionID, roomID); senderr != nil {
			fmt.Printf("failed to send event. dest[%s] from[%s] room[%s]", roommate.ConnectionID, connectionID, roomID)
		}
	}

	return nil
}

// LeaveRoomMate ...
func LeaveRoomMate(socket *socket.Socket, connectionID string, roomID string) error {
	cm, err := connection.NewManager()
	if err != nil {
		return err
	}

	roomMates, err := cm.FindRoomMates(roomID)
	for _, roommate := range roomMates {
		socket.SendRoomLeave(roommate.ConnectionID, connectionID, roomID)
	}

	return nil
}
