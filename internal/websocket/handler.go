package websocket

import (
	"context"
	"encoding/json"
	"hrm-app/internal/domain/boards"
	"hrm-app/internal/domain/boardsUsers"
	"hrm-app/internal/domain/contact"
	"hrm-app/internal/domain/labels"
	room_chats "hrm-app/internal/domain/roomChats"
	room_messages "hrm-app/internal/domain/roomMessages"
	"hrm-app/internal/domain/roomUsers"
	"hrm-app/internal/domain/taskCard"
	"hrm-app/internal/domain/taskCardComment"
	"hrm-app/internal/domain/taskCardUsers"
	"hrm-app/internal/domain/taskTab"
	"hrm-app/internal/domain/user"
	"hrm-app/internal/domain/workspacesUsers"
	"hrm-app/internal/response"
	"hrm-app/internal/websocket/handlerWebsocket"
	"log"
	"net/http"
	"time"

	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		// Allow development origins
		if origin == "http://localhost:3000" || origin == "http://localhost:8080" {
			return true
		}
		// Allow production origin
		if origin == "http://app.putratek.my.id" || origin == "http://be.putratek.my.id" {
			return true
		}
		// Allow all during competition/dev if specific domain is unknown or dynamic
		return true
	},
}

type Handler struct {
	hub              *Hub
	boardHandler     *handlerWebsocket.BoardHandler
	taskCardHandler  *handlerWebsocket.TaskCardHandler
	taskTabHandler   *handlerWebsocket.TaskTabHandler
	commentHandler   *handlerWebsocket.CommentHandler
	labelHandler     *handlerWebsocket.LabelHandler
	workspaceHandler *handlerWebsocket.WorkspaceHandler
	chatHandler      *handlerWebsocket.ChatHandler
	contactUC        contact.UseCase
	userUC           user.UseCase
}

func NewHandler(hub *Hub, taskCardUC taskCard.UseCase, taskTabUC taskTab.UseCase, commentUC taskCardComment.UseCase, labelsUC labels.UseCase, taskCardUsersUC taskCardUsers.UseCase, boardsUsersUC boardsUsers.UseCase, workspacesUsersUC workspacesUsers.UseCase, boardsUC boards.UseCase, roomMessageUC room_messages.UseCase, roomChatUC room_chats.UseCase, roomUserUC roomUsers.UseCase, contactUC contact.UseCase, userUC user.UseCase) *Handler {
	return &Handler{
		hub:              hub,
		boardHandler:     handlerWebsocket.NewBoardHandler(boardsUC, boardsUsersUC, hub),
		taskCardHandler:  handlerWebsocket.NewTaskCardHandler(taskCardUC, taskTabUC, taskCardUsersUC, hub),
		taskTabHandler:   handlerWebsocket.NewTaskTabHandler(taskTabUC, hub),
		commentHandler:   handlerWebsocket.NewCommentHandler(commentUC, taskCardUC, taskTabUC, hub),
		labelHandler:     handlerWebsocket.NewLabelHandler(labelsUC, taskCardUC, taskTabUC, hub),
		workspaceHandler: handlerWebsocket.NewWorkspaceHandler(workspacesUsersUC, hub),
		chatHandler:      handlerWebsocket.NewChatHandler(roomMessageUC, roomChatUC, roomUserUC, hub),
		contactUC:        contactUC,
		userUC:           userUC,
	}
}

// WSMessage represents the structure of WebSocket messages
type WSMessage struct {
	Action  string          `json:"action"`
	Payload json.RawMessage `json:"payload"`
}

// HandleWebSocket handles WebSocket connections
func (h *Handler) HandleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}

	uID, _ := c.Get("user_id")
	userID := uID.(uint)
	userIDStr := strconv.FormatUint(uint64(userID), 10)

	// Fetch user's name from contact domain and username from user domain
	userName := "User"
	userUsername := "user"
	if contact, err := h.contactUC.GetByUserID(context.Background(), userID, ""); err == nil && contact != nil {
		userName = contact.Name
	}
	if user, err := h.userUC.GetByID(userID); err == nil && user != nil {
		userUsername = user.Username
	}

	// Create context with cancel for the client
	ctx, cancel := context.WithCancel(context.Background())

	client := &Client{
		hub:          h.hub,
		conn:         conn,
		send:         make(chan []byte, 256),
		userID:       userIDStr,
		UserID:       userID,
		UserName:     userName,
		UserUsername: userUsername,
		Ctx:          ctx,
		Cancel:       cancel,
	}

	client.hub.register <- client

	// Setup RabbitMQ for this user
	if err := client.setupRabbitMQ(); err != nil {
		log.Printf("Failed to setup RabbitMQ for %d: %v", userID, err)
		_ = client.conn.Close()
		return
	}

	// Start goroutines for reading and writing
	go client.writePump()
	go h.handleMessages(client)
}

// handleMessages processes incoming WebSocket messages
func (h *Handler) handleMessages(client *Client) {
	defer func() {
		if client.Cancel != nil {
			client.Cancel()
		}
		client.hub.unregister <- client
		_ = client.conn.Close()
	}()

	// Configure WebSocket connection safety
	client.conn.SetReadLimit(512) // Hardcoded 512 bytes limit for text messages
	_ = client.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	client.conn.SetPongHandler(func(string) error { _ = client.conn.SetReadDeadline(time.Now().Add(60 * time.Second)); return nil })

	for {
		_, message, err := client.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		// Rate Limit Check
		rateLimiter := h.hub.GetRateLimiter()
		result := rateLimiter.CheckLimit(client.userID, 10, 60) // 10 msg per 60 sec

		if !result.Allowed {
			log.Printf("Rate limit exceeded for user %s", client.userID)
			h.sendErrorToClient(client, "Rate limit exceeded. Please slow down.")
			continue
		}

		// Normalize message (trim newline/space) similar to client.go implementation
		// Note: "bytes" package is not imported in handler.go, need to add it or skip if not critical.
		// JSON unmarshal usually handles whitespace, so strictly trimming might not be needed for JSON,
		// but let's stick to safety. I'll add the checks first.

		var msg WSMessage
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Printf("Error unmarshaling message: %v", err)
			h.sendErrorToClient(client, "Invalid message format")
			continue
		}

		// Handle different actions using domain handlers
		switch msg.Action {
		// Board Actions
		case "join_board":
			h.boardHandler.HandleJoinBoard(client, msg.Payload)
		case "assign_board_user":
			h.boardHandler.HandleAssignBoardUser(client, msg.Payload)
		case "unassign_board_user":
			h.boardHandler.HandleUnassignBoardUser(client, msg.Payload)

		// Task Card Actions
		case "create_task_card":
			h.taskCardHandler.HandleCreateTaskCard(client, msg.Payload)
		case "update_task_tab_id":
			h.taskCardHandler.HandleUpdateTaskTabID(client, msg.Payload)
		case "update_task_card":
			h.taskCardHandler.HandleUpdateTaskCard(client, msg.Payload)
		case "assign_task_card_user":
			h.taskCardHandler.HandleAssignTaskCardUser(client, msg.Payload)
		case "unassign_task_card_user":
			h.taskCardHandler.HandleUnassignTaskCardUser(client, msg.Payload)

		// Task Tab Actions
		case "update_task_tab":
			h.taskTabHandler.HandleUpdateTaskTab(client, msg.Payload)

		// Comment Actions
		case "create_task_card_comment":
			h.commentHandler.HandleCreateTaskCardComment(client, msg.Payload)
		case "update_task_card_comment":
			h.commentHandler.HandleUpdateTaskCardComment(client, msg.Payload)
		case "delete_task_card_comment":
			h.commentHandler.HandleDeleteTaskCardComment(client, msg.Payload)

		// Label Actions
		case "create_label":
			h.labelHandler.HandleCreateLabel(client, msg.Payload)
		case "update_label":
			h.labelHandler.HandleUpdateLabel(client, msg.Payload)
		case "delete_label":
			h.labelHandler.HandleDeleteLabel(client, msg.Payload)

		// Workspace Actions
		case "assign_workspace_user":
			h.workspaceHandler.HandleAssignWorkspaceUser(client, msg.Payload)
		case "unassign_workspace_user":
			h.workspaceHandler.HandleUnassignWorkspaceUser(client, msg.Payload)

		// Chat Room Actions
		case "join_room_chat":
			h.chatHandler.HandleJoinRoomChat(client, msg.Payload)
		case "send_room_chat_message":
			h.chatHandler.HandleSendRoomChatMessage(client, msg.Payload)
		case "edit_room_chat_message":
			h.chatHandler.HandleEditRoomChatMessage(client, msg.Payload)
		case "delete_room_chat_message":
			h.chatHandler.HandleDeleteRoomChatMessage(client, msg.Payload)
		case "typing_indicator":
			h.chatHandler.HandleTypingIndicator(client, msg.Payload)

		default:
			h.sendErrorToClient(client, "Unknown action")
		}
	}
}

// sendErrorToClient sends an error message to a specific client
func (h *Handler) sendErrorToClient(client *Client, errorMsg string) {
	response := map[string]interface{}{
		"status": "error",
		"error":  errorMsg,
	}
	responseJSON, _ := json.Marshal(response)
	client.Send(responseJSON)
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
