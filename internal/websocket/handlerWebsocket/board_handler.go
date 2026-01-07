package handlerWebsocket

import (
	"encoding/json"
	"hrm-app/internal/domain/boards"
	"hrm-app/internal/domain/boardsUsers"
	"log"
)

type BoardHandler struct {
	BaseHandler
	boardsUseCase      boards.UseCase
	boardsUsersUseCase boardsUsers.UseCase
	hub                Hub
}

func NewBoardHandler(boardsUseCase boards.UseCase, boardsUsersUseCase boardsUsers.UseCase, hub Hub) *BoardHandler {
	return &BoardHandler{
		boardsUseCase:      boardsUseCase,
		boardsUsersUseCase: boardsUsersUseCase,
		hub:                hub,
	}
}

type JoinBoardPayload struct {
	BoardID uint `json:"board_id"`
}

type AssignBoardUserPayload struct {
	BoardID uint `json:"board_id"`
	UserID  uint `json:"user_id"`
}

type UnassignBoardUserPayload struct {
	ID uint `json:"id"`
}

func (h *BoardHandler) HandleJoinBoard(client Client, payload json.RawMessage) {
	var msg JoinBoardPayload
	if err := json.Unmarshal(payload, &msg); err != nil {
		h.SendError(client, "join_board", "Invalid payload")
		return
	}

	// Check if user is member of the board
	hasAccess, err := h.boardsUsersUseCase.HasAccess(msg.BoardID, client.GetUserID())
	if err != nil || !hasAccess {
		log.Printf("[WS Auth] Unauthorized board join attempt: UserID=%d, BoardID=%d", client.GetUserID(), msg.BoardID)
		h.SendError(client, "join_board", "Unauthorized: You are not a member of this board")
		return
	}

	h.hub.RegisterClientToBoard(client, msg.BoardID)

	log.Printf("[WS] User %d joined board %d", client.GetUserID(), msg.BoardID)

	h.SendSuccess(client, "join_board", msg, map[string]interface{}{"board_id": msg.BoardID})
}

func (h *BoardHandler) HandleAssignBoardUser(client Client, payload json.RawMessage) {
	var msg AssignBoardUserPayload
	if err := json.Unmarshal(payload, &msg); err != nil {
		h.SendError(client, "assign_board_user", "Invalid payload")
		return
	}

	assignment := &boardsUsers.BoardsUsers{
		BoardID: msg.BoardID,
		UserID:  msg.UserID,
	}

	if err := h.boardsUsersUseCase.Create(assignment, msg.UserID); err != nil {
		h.SendError(client, "assign_board_user", "Failed to assign user to board: "+err.Error())
		return
	}

	fullAssignment, err := h.boardsUsersUseCase.GetByID(assignment.ID)
	if err != nil {
		h.SendError(client, "assign_board_user", "Failed to fetch assignment details")
		return
	}

	h.SendSuccess(client, "assign_board_user", msg, fullAssignment)
	h.BroadcastSuccess(h.hub, msg.BoardID, "assign_board_user", msg, fullAssignment)
}

func (h *BoardHandler) HandleUnassignBoardUser(client Client, payload json.RawMessage) {
	var msg UnassignBoardUserPayload
	if err := json.Unmarshal(payload, &msg); err != nil {
		h.SendError(client, "unassign_board_user", "Invalid payload")
		return
	}

	assignment, err := h.boardsUsersUseCase.GetByID(msg.ID)
	if err != nil {
		h.SendError(client, "unassign_board_user", "Assignment not found")
		return
	}

	if err := h.boardsUsersUseCase.Delete(msg.ID); err != nil {
		h.SendError(client, "unassign_board_user", "Failed to unassign user from board: "+err.Error())
		return
	}

	h.SendSuccess(client, "unassign_board_user", msg, map[string]interface{}{"id": msg.ID})
	h.BroadcastSuccess(h.hub, assignment.BoardID, "unassign_board_user", msg, map[string]interface{}{"id": msg.ID})
}
