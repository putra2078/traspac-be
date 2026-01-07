package handlerWebsocket

import (
	"encoding/json"
	room_chats "hrm-app/internal/domain/roomChats"
	room_messages "hrm-app/internal/domain/roomMessages"
	"hrm-app/internal/domain/roomUsers"
)

type ChatHandler struct {
	BaseHandler
	roomMessageUC room_messages.UseCase
	roomChatUC    room_chats.UseCase
	roomUserUC    roomUsers.UseCase
	hub           Hub
}

func NewChatHandler(roomMessageUC room_messages.UseCase, roomChatUC room_chats.UseCase, roomUserUC roomUsers.UseCase, hub Hub) *ChatHandler {
	return &ChatHandler{
		roomMessageUC: roomMessageUC,
		roomChatUC:    roomChatUC,
		roomUserUC:    roomUserUC,
		hub:           hub,
	}
}

func (h *ChatHandler) HandleJoinRoomChat(client Client, payload json.RawMessage) {
	var data struct {
		RoomID uint `json:"room_id"`
	}
	if err := json.Unmarshal(payload, &data); err != nil {
		h.SendError(client, "join_room_chat", "Invalid payload")
		return
	}

	// Verify room exists
	_, err := h.roomChatUC.GetByID(data.RoomID)
	if err != nil {
		h.SendError(client, "join_room_chat", "Room not found")
		return
	}

	// Add user to room_users
	userID := client.GetUserID()
	if err := h.roomUserUC.Join(data.RoomID, userID); err != nil {
		h.SendError(client, "join_room_chat", "Failed to join room: "+err.Error())
		return
	}

	h.hub.RegisterClientToChatRoom(client, data.RoomID)

	// Fetch chat history
	history, err := h.roomMessageUC.GetChatHistory(client.GetContext(), data.RoomID)
	if err != nil {
		h.SendError(client, "join_room_chat", "Failed to fetch chat history")
		return
	}

	h.SendSuccess(client, "join_room_chat", "Successfully joined chat room", history)
}

func (h *ChatHandler) HandleSendRoomChatMessage(client Client, payload json.RawMessage) {
	var data struct {
		RoomID         uint   `json:"room_id"`
		MessageText    string `json:"message_text"`
		MessageContent string `json:"message_content"`
	}
	if err := json.Unmarshal(payload, &data); err != nil {
		h.SendError(client, "send_room_chat_message", "Invalid payload")
		return
	}

	userID := client.GetUserID()

	// Verify user is in room
	inRoom, err := h.roomUserUC.IsUserInRoom(data.RoomID, userID)
	if err != nil || !inRoom {
		h.SendError(client, "send_room_chat_message", "Unauthorized: you are not a member of this room")
		return
	}

	message := room_messages.RoomMessage{
		RoomID:         data.RoomID,
		UserID:         &userID,
		MessageText:    data.MessageText,
		MessageContent: data.MessageContent,
	}

	// Save to DB
	// Save to DB
	savedMessage, err := h.roomMessageUC.SendMessage(client.GetContext(), message)
	if err != nil {
		h.SendError(client, "send_room_chat_message", "Failed to save message")
		return
	}

	// Prepare message for broadcast
	msgJSON, _ := json.Marshal(map[string]interface{}{
		"action":          "new_room_chat_message",
		"status":          "success",
		"data":            savedMessage,
		"sender_name":     client.GetUserName(),
		"sender_username": client.GetUserUsername(),
	})

	// Broadcast via Hub (Local + Kafka)
	h.hub.BroadcastToChatRoom(data.RoomID, msgJSON)
}

func (h *ChatHandler) HandleTypingIndicator(client Client, payload json.RawMessage) {
	var data struct {
		RoomID   uint `json:"room_id"`
		IsTyping bool `json:"is_typing"`
	}
	if err := json.Unmarshal(payload, &data); err != nil {
		return
	}

	userID := client.GetUserID()

	// Broadcast status locally only (transient, no Kafka needed for now to keep it light,
	// but can be changed to BroadcastToChatRoom if multi-instance typing is required)
	msgJSON, _ := json.Marshal(map[string]interface{}{
		"action": "typing_indicator",
		"status": "success",
		"data": map[string]interface{}{
			"room_id":   data.RoomID,
			"user_id":   userID,
			"user_name": client.GetUserName(),
			"is_typing": data.IsTyping,
		},
	})

	h.hub.BroadcastToChatRoomLocal(data.RoomID, msgJSON)
}

func (h *ChatHandler) HandleEditRoomChatMessage(client Client, payload json.RawMessage) {
	var data struct {
		MessageID   uint   `json:"message_id"`
		RoomID      uint   `json:"room_id"` // Needed for broadcast routing
		MessageText string `json:"message_text"`
	}
	if err := json.Unmarshal(payload, &data); err != nil {
		h.SendError(client, "edit_room_chat_message", "Invalid payload")
		return
	}

	userID := client.GetUserID()

	updatedMsg, err := h.roomMessageUC.EditMessage(client.GetContext(), userID, data.MessageID, data.MessageText)
	if err != nil {
		h.SendError(client, "edit_room_chat_message", "Failed to edit message: "+err.Error())
		return
	}

	// Prepare data for broadcast
	msgJSON, _ := json.Marshal(map[string]interface{}{
		"action":          "edit_room_chat_message",
		"status":          "success",
		"data":            updatedMsg,
		"sender_name":     client.GetUserName(),
		"sender_username": client.GetUserUsername(),
	})

	h.hub.BroadcastToChatRoom(data.RoomID, msgJSON)
}

func (h *ChatHandler) HandleDeleteRoomChatMessage(client Client, payload json.RawMessage) {
	var data struct {
		MessageID uint `json:"message_id"`
		RoomID    uint `json:"room_id"` // Needed for broadcast routing
	}
	if err := json.Unmarshal(payload, &data); err != nil {
		h.SendError(client, "delete_room_chat_message", "Invalid payload")
		return
	}

	userID := client.GetUserID()

	if err := h.roomMessageUC.DeleteMessage(client.GetContext(), userID, data.MessageID); err != nil {
		h.SendError(client, "delete_room_chat_message", "Failed to delete message: "+err.Error())
		return
	}

	// Prepare data for broadcast
	msgJSON, _ := json.Marshal(map[string]interface{}{
		"action": "delete_room_chat_message",
		"status": "success",
		"data": map[string]interface{}{
			"message_id": data.MessageID,
			"room_id":    data.RoomID,
		},
	})

	h.hub.BroadcastToChatRoom(data.RoomID, msgJSON)
}
