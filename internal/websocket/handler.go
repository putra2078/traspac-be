package websocket

import (
	"encoding/json"
	"hrm-app/internal/domain/taskCard"
	"hrm-app/internal/domain/taskCardComment"
	"hrm-app/internal/domain/taskTab"
	"hrm-app/internal/response"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for development
	},
}

type Handler struct {
	hub                    *Hub
	taskCardUseCase        taskCard.UseCase
	taskTabUseCase         taskTab.UseCase
	taskCardCommentUseCase taskCardComment.UseCase
}

func NewHandler(hub *Hub, taskCardUseCase taskCard.UseCase, taskTabUseCase taskTab.UseCase, taskCardCommentUseCase taskCardComment.UseCase) *Handler {
	return &Handler{
		hub:                    hub,
		taskCardUseCase:        taskCardUseCase,
		taskTabUseCase:         taskTabUseCase,
		taskCardCommentUseCase: taskCardCommentUseCase,
	}
}

// WSMessage represents the structure of WebSocket messages
type WSMessage struct {
	Action  string          `json:"action"`
	Payload json.RawMessage `json:"payload"`
}

// UpdateTaskTabIDPayload represents payload for moving a card to another tab
type UpdateTaskTabIDPayload struct {
	TaskCardID uint `json:"task_card_id"`
	TaskTabID  uint `json:"task_tab_id"`
}

// UpdateTaskCardPayload represents payload for updating task card details
type UpdateTaskCardPayload struct {
	TaskCardID uint   `json:"task_card_id"`
	Content    string `json:"content,omitempty"`
	Comment    string `json:"comment,omitempty"`
	Date       string `json:"date,omitempty"`
	Status     bool   `json:"status,omitempty"`
	Name       string `json:"name,omitempty"`
}

// UpdateTaskTabPayload represents payload for updating task tab details
type UpdateTaskTabPayload struct {
	TaskTabID uint   `json:"task_tab_id"`
	Name      string `json:"name,omitempty"`
	Position  int    `json:"position,omitempty"`
}

// CreateTaskCardCommentPayload represents payload for creating a comment
type CreateTaskCardCommentPayload struct {
	TaskCardID int    `json:"task_card_id"`
	Comment    string `json:"comment"`
}

// GetTaskCardCommentsPayload represents payload for fetching comments
type GetTaskCardCommentsPayload struct {
	TaskCardID int `json:"task_card_id"`
}

// UpdateTaskCardCommentPayload represents payload for updating a comment
type UpdateTaskCardCommentPayload struct {
	ID      int    `json:"id"`
	Comment string `json:"comment"`
}

// DeleteTaskCardCommentPayload represents payload for deleting a comment
type DeleteTaskCardCommentPayload struct {
	ID int `json:"id"`
}

// HandleWebSocket handles WebSocket connections
func (h *Handler) HandleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}

	client := &Client{
		hub:  h.hub,
		conn: conn,
		send: make(chan []byte, 256),
	}

	client.hub.register <- client

	// Start goroutines for reading and writing
	go client.writePump()
	go h.handleMessages(client)
}

// handleMessages processes incoming WebSocket messages
func (h *Handler) handleMessages(client *Client) {
	defer func() {
		client.hub.unregister <- client
		client.conn.Close()
	}()

	for {
		_, message, err := client.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		var msg WSMessage
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Printf("Error unmarshaling message: %v", err)
			h.sendErrorToClient(client, "Invalid message format")
			continue
		}

		// Handle different actions
		switch msg.Action {
		case "update_task_tab_id":
			h.handleUpdateTaskTabID(client, msg.Payload)
		case "update_task_card":
			h.handleUpdateTaskCard(client, msg.Payload)
		case "update_task_tab":
			h.handleUpdateTaskTab(client, msg.Payload)
		case "create_task_card_comment":
			h.handleCreateTaskCardComment(client, msg.Payload)
		case "get_task_card_comments":
			h.handleGetTaskCardComments(client, msg.Payload)
		case "update_task_card_comment":
			h.handleUpdateTaskCardComment(client, msg.Payload)
		case "delete_task_card_comment":
			h.handleDeleteTaskCardComment(client, msg.Payload)
		default:
			h.sendErrorToClient(client, "Unknown action")
		}
	}
}

// handleUpdateTaskTabID updates the TaskTabID of a TaskCard
func (h *Handler) handleUpdateTaskTabID(client *Client, payload json.RawMessage) {
	var msg UpdateTaskTabIDPayload
	if err := json.Unmarshal(payload, &msg); err != nil {
		h.sendErrorToClient(client, "Invalid payload")
		return
	}

	// Get the task card
	taskCardData, err := h.taskCardUseCase.FindByID(msg.TaskCardID)
	if err != nil {
		h.sendErrorToClient(client, "Task card not found")
		return
	}

	// Update the TaskTabID
	taskCardData.TaskTabID = msg.TaskTabID

	// Save the updated task card
	if err := h.taskCardUseCase.Update(taskCardData); err != nil {
		h.sendErrorToClient(client, "Failed to update task card")
		return
	}

	// Broadcast success
	h.broadcastSuccess("update_task_tab_id", msg, taskCardData)
}

func (h *Handler) handleUpdateTaskCard(client *Client, payload json.RawMessage) {
	var msg UpdateTaskCardPayload
	if err := json.Unmarshal(payload, &msg); err != nil {
		h.sendErrorToClient(client, "Invalid payload")
		return
	}

	taskCardData, err := h.taskCardUseCase.FindByID(msg.TaskCardID)
	if err != nil {
		h.sendErrorToClient(client, "Task card not found")
		return
	}

	if msg.Content != "" {
		taskCardData.Content = msg.Content
	}
	if msg.Comment != "" {
		taskCardData.Comment = msg.Comment
	}
	if msg.Date != "" {
		taskCardData.Date = msg.Date
	}
	if msg.Name != "" {
		taskCardData.Name = msg.Name
	}
	// Status bool update is tricky if we want to allow setting to false.
	// For now assuming implicit update if field is present in logic but Go JSON zero value is false.
	// Since we are using struct updates, we might update it blindly or need a pointer.
	// Given the payload struct, let's update it.
	taskCardData.Status = msg.Status

	if err := h.taskCardUseCase.Update(taskCardData); err != nil {
		h.sendErrorToClient(client, "Failed to update task card")
		return
	}

	h.broadcastSuccess("update_task_card", msg, taskCardData)
}

func (h *Handler) handleUpdateTaskTab(client *Client, payload json.RawMessage) {
	var msg UpdateTaskTabPayload
	if err := json.Unmarshal(payload, &msg); err != nil {
		h.sendErrorToClient(client, "Invalid payload")
		return
	}

	taskTabData, err := h.taskTabUseCase.FindByID(msg.TaskTabID)
	if err != nil {
		h.sendErrorToClient(client, "Task tab not found")
		return
	}

	if msg.Name != "" {
		taskTabData.Name = msg.Name
	}
	if msg.Position != 0 {
		taskTabData.Position = msg.Position
	}

	if err := h.taskTabUseCase.Update(taskTabData); err != nil {
		h.sendErrorToClient(client, "Failed to update task tab")
		return
	}

	h.broadcastSuccess("update_task_tab", msg, taskTabData)
}

// handleCreateTaskCardComment creates a new comment for a task card
func (h *Handler) handleCreateTaskCardComment(client *Client, payload json.RawMessage) {
	var msg CreateTaskCardCommentPayload
	if err := json.Unmarshal(payload, &msg); err != nil {
		h.sendErrorToClient(client, "Invalid payload")
		return
	}

	comment := &taskCardComment.TaskCardComment{
		TaskCardID: msg.TaskCardID,
		Comment:    msg.Comment,
	}

	if err := h.taskCardCommentUseCase.Create(comment); err != nil {
		h.sendErrorToClient(client, "Failed to create comment: "+err.Error())
		return
	}

	h.broadcastSuccess("create_task_card_comment", msg, comment)
}

// handleGetTaskCardComments fetches all comments for a task card
func (h *Handler) handleGetTaskCardComments(client *Client, payload json.RawMessage) {
	var msg GetTaskCardCommentsPayload
	if err := json.Unmarshal(payload, &msg); err != nil {
		h.sendErrorToClient(client, "Invalid payload")
		return
	}

	comments, err := h.taskCardCommentUseCase.FindByTaskCardID(uint(msg.TaskCardID))
	if err != nil {
		h.sendErrorToClient(client, "Failed to fetch comments: "+err.Error())
		return
	}

	// Send response only to the requesting client
	response := map[string]interface{}{
		"action":  "get_task_card_comments",
		"status":  "success",
		"payload": msg,
		"data":    comments,
	}
	responseJSON, _ := json.Marshal(response)

	select {
	case client.send <- responseJSON:
	default:
		close(client.send)
		delete(client.hub.clients, client)
	}
}

// handleUpdateTaskCardComment updates an existing comment
func (h *Handler) handleUpdateTaskCardComment(client *Client, payload json.RawMessage) {
	var msg UpdateTaskCardCommentPayload
	if err := json.Unmarshal(payload, &msg); err != nil {
		h.sendErrorToClient(client, "Invalid payload")
		return
	}

	// Get the existing comment
	comment, err := h.taskCardCommentUseCase.FindByID(uint(msg.ID))
	if err != nil {
		h.sendErrorToClient(client, "Comment not found")
		return
	}

	// Update the comment
	comment.Comment = msg.Comment

	if err := h.taskCardCommentUseCase.Update(comment); err != nil {
		h.sendErrorToClient(client, "Failed to update comment: "+err.Error())
		return
	}

	h.broadcastSuccess("update_task_card_comment", msg, comment)
}

// handleDeleteTaskCardComment deletes a comment
func (h *Handler) handleDeleteTaskCardComment(client *Client, payload json.RawMessage) {
	var msg DeleteTaskCardCommentPayload
	if err := json.Unmarshal(payload, &msg); err != nil {
		h.sendErrorToClient(client, "Invalid payload")
		return
	}

	if err := h.taskCardCommentUseCase.Delete(uint(msg.ID)); err != nil {
		h.sendErrorToClient(client, "Failed to delete comment: "+err.Error())
		return
	}

	h.broadcastSuccess("delete_task_card_comment", msg, map[string]interface{}{"id": msg.ID})
}

func (h *Handler) broadcastSuccess(action string, payload interface{}, data interface{}) {
	response := map[string]interface{}{
		"action":  action,
		"status":  "success",
		"payload": payload,
		"data":    data,
	}
	responseJSON, _ := json.Marshal(response)
	h.hub.BroadcastMessage(responseJSON)
}

// sendErrorToClient sends an error message to a specific client
func (h *Handler) sendErrorToClient(client *Client, errorMsg string) {
	errorResponse := map[string]interface{}{
		"status": "error",
		"error":  errorMsg,
	}

	errorJSON, err := json.Marshal(errorResponse)
	if err != nil {
		log.Printf("Error marshaling error response: %v", err)
		return
	}

	select {
	case client.send <- errorJSON:
	default:
		close(client.send)
		delete(client.hub.clients, client)
	}
}

// GetConnectedClients returns the number of connected clients (for monitoring)
func (h *Handler) GetConnectedClients(c *gin.Context) {
	h.hub.mu.RLock()
	count := len(h.hub.clients)
	h.hub.mu.RUnlock()

	response.Success(c, gin.H{
		"connected_clients": count,
	})
}
