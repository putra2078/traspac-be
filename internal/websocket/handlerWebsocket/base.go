package handlerWebsocket

import (
	"context"
	"encoding/json"
)

// Client interface defines the methods needed by handlers to interact with a client
type Client interface {
	GetUserID() uint
	GetUserName() string
	GetUserUsername() string
	Send(message []byte)
	Close()
	GetContext() context.Context
}

// Hub interface defines the methods needed by handlers to interact with the hub
type Hub interface {
	RegisterClientToBoard(client Client, boardID uint)
	RegisterClientToChatRoom(client Client, roomID uint)
	BroadcastToBoard(boardID uint, message []byte)
	BroadcastToChatRoom(roomID uint, message []byte)
	BroadcastToChatRoomLocal(roomID uint, message []byte)
	BroadcastMessage(message []byte)
}

// BaseHandler provides common utility methods for all WebSocket handlers
type BaseHandler struct{}

// SendError sends an error message to a client
func (bh *BaseHandler) SendError(client Client, action string, errorMsg string) {
	response := map[string]interface{}{
		"action": action,
		"status": "error",
		"error":  errorMsg,
	}
	responseJSON, _ := json.Marshal(response)
	client.Send(responseJSON)
}

// SendSuccess sends a success message to a client
func (bh *BaseHandler) SendSuccess(client Client, action string, payload interface{}, data interface{}) {
	response := map[string]interface{}{
		"action":  action,
		"status":  "success",
		"payload": payload,
		"data":    data,
	}
	responseJSON, _ := json.Marshal(response)
	client.Send(responseJSON)
}

// BroadcastSuccess broadcasts a success message to a board
func (bh *BaseHandler) BroadcastSuccess(hub Hub, boardID uint, action string, payload interface{}, data interface{}) {
	response := map[string]interface{}{
		"action":  action,
		"status":  "success",
		"payload": payload,
		"data":    data,
	}
	responseJSON, _ := json.Marshal(response)
	hub.BroadcastToBoard(boardID, responseJSON)
}

// BroadcastGlobalSuccess broadcasts a success message to all connected clients
func (bh *BaseHandler) BroadcastGlobalSuccess(hub Hub, action string, payload interface{}, data interface{}) {
	response := map[string]interface{}{
		"action":  action,
		"status":  "success",
		"payload": payload,
		"data":    data,
	}
	responseJSON, _ := json.Marshal(response)
	hub.BroadcastMessage(responseJSON)
}
