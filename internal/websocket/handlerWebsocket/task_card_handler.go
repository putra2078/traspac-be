package handlerWebsocket

import (
	"context"
	"encoding/json"
	"hrm-app/internal/domain/taskCard"
	"hrm-app/internal/domain/taskCardUsers"
	"hrm-app/internal/domain/taskTab"
)

type TaskCardHandler struct {
	BaseHandler
	taskCardUseCase      taskCard.UseCase
	taskTabUseCase       taskTab.UseCase
	taskCardUsersUseCase taskCardUsers.UseCase
	hub                  Hub
}

func NewTaskCardHandler(taskCardUseCase taskCard.UseCase, taskTabUseCase taskTab.UseCase, taskCardUsersUseCase taskCardUsers.UseCase, hub Hub) *TaskCardHandler {
	return &TaskCardHandler{
		taskCardUseCase:      taskCardUseCase,
		taskTabUseCase:       taskTabUseCase,
		taskCardUsersUseCase: taskCardUsersUseCase,
		hub:                  hub,
	}
}

type UpdateTaskTabIDPayload struct {
	TaskCardID uint `json:"task_card_id"`
	TaskTabID  uint `json:"task_tab_id"`
}

type UpdateTaskCardPayload struct {
	TaskCardID uint   `json:"task_card_id"`
	TaskTabID  uint   `json:"task_tab_id,omitempty"`
	Content    string `json:"content,omitempty"`
	Comment    string `json:"comment,omitempty"`
	Date       string `json:"date,omitempty"`
	Status     *bool  `json:"status,omitempty"`
	Name       string `json:"name,omitempty"`
}

type AssignTaskCardUserPayload struct {
	TaskCardID uint `json:"task_card_id"`
	UserID     uint `json:"user_id"`
}

type UnassignTaskCardUserPayload struct {
	ID uint `json:"id"`
}

type CreateTaskCardPayload struct {
	TaskTabID uint   `json:"task_tab_id"`
	Name      string `json:"name"`
	Content   string `json:"content,omitempty"`
	Date      string `json:"date,omitempty"`
}

func (h *TaskCardHandler) HandleUpdateTaskTabID(client Client, payload json.RawMessage) {
	var msg UpdateTaskTabIDPayload
	if err := json.Unmarshal(payload, &msg); err != nil {
		h.SendError(client, "update_task_tab_id", "Invalid payload")
		return
	}

	taskCardData, err := h.taskCardUseCase.FindByID(context.Background(), msg.TaskCardID)
	if err != nil {
		h.SendError(client, "update_task_tab_id", "Task card not found")
		return
	}

	taskCardData.TaskTabID = msg.TaskTabID

	if err := h.taskCardUseCase.Update(context.Background(), taskCardData); err != nil {
		h.SendError(client, "update_task_tab_id", "Failed to update task card")
		return
	}

	// Fetch fresh data with preloads and updated fields
	freshTaskCard, err := h.taskCardUseCase.FindByID(context.Background(), msg.TaskCardID)
	if err != nil {
		h.SendError(client, "update_task_tab_id", "Failed to refresh task card data")
		return
	}

	taskTab, err := h.taskTabUseCase.FindByID(freshTaskCard.TaskTabID)
	if err != nil {
		h.SendError(client, "update_task_tab_id", "Task tab not found")
		return
	}

	h.SendSuccess(client, "update_task_tab_id", msg, freshTaskCard)
	h.BroadcastSuccess(h.hub, taskTab.BoardID, "update_task_tab_id", msg, freshTaskCard)
}

func (h *TaskCardHandler) HandleUpdateTaskCard(client Client, payload json.RawMessage) {
	var msg UpdateTaskCardPayload
	if err := json.Unmarshal(payload, &msg); err != nil {
		h.SendError(client, "update_task_card", "Invalid payload")
		return
	}

	taskCardData, err := h.taskCardUseCase.FindByID(context.Background(), msg.TaskCardID)
	if err != nil {
		h.SendError(client, "update_task_card", "Task card not found")
		return
	}

	if msg.TaskTabID != 0 {
		taskCardData.TaskTabID = msg.TaskTabID
	}
	if msg.Content != "" {
		taskCardData.Content = msg.Content
	}
	if msg.Date != "" {
		taskCardData.Date = msg.Date
	}
	if msg.Name != "" {
		taskCardData.Name = msg.Name
	}
	if msg.Status != nil {
		taskCardData.Status = *msg.Status
	}

	if err := h.taskCardUseCase.Update(context.Background(), taskCardData); err != nil {
		h.SendError(client, "update_task_card", "Failed to update task card")
		return
	}

	// Fetch fresh data with preloads and updated fields
	freshTaskCard, err := h.taskCardUseCase.FindByID(context.Background(), msg.TaskCardID)
	if err != nil {
		h.SendError(client, "update_task_card", "Failed to refresh task card data")
		return
	}

	taskTab, err := h.taskTabUseCase.FindByID(freshTaskCard.TaskTabID)
	if err != nil {
		h.SendError(client, "update_task_card", "Task tab not found")
		return
	}

	h.SendSuccess(client, "update_task_card", msg, freshTaskCard)
	h.BroadcastSuccess(h.hub, taskTab.BoardID, "update_task_card", msg, freshTaskCard)
}

func (h *TaskCardHandler) HandleAssignTaskCardUser(client Client, payload json.RawMessage) {
	var msg AssignTaskCardUserPayload
	if err := json.Unmarshal(payload, &msg); err != nil {
		h.SendError(client, "assign_task_card_user", "Invalid payload")
		return
	}

	assignment := &taskCardUsers.TaskCardUsers{
		TaskCardID: msg.TaskCardID,
		UserID:     msg.UserID,
	}

	if err := h.taskCardUsersUseCase.Create(assignment); err != nil {
		h.SendError(client, "assign_task_card_user", "Failed to assign user: "+err.Error())
		return
	}

	fullAssignment, err := h.taskCardUsersUseCase.GetByID(assignment.ID)
	if err != nil {
		h.SendError(client, "assign_task_card_user", "Failed to get assignment details")
		return
	}

	taskCard, err := h.taskCardUseCase.FindByID(context.Background(), msg.TaskCardID)
	if err != nil {
		h.SendError(client, "assign_task_card_user", "Task card not found")
		return
	}
	taskTab, err := h.taskTabUseCase.FindByID(taskCard.TaskTabID)
	if err != nil {
		h.SendError(client, "assign_task_card_user", "Task tab not found")
		return
	}

	h.SendSuccess(client, "assign_task_card_user", msg, fullAssignment)
	h.BroadcastSuccess(h.hub, taskTab.BoardID, "assign_task_card_user", msg, fullAssignment)
}

func (h *TaskCardHandler) HandleUnassignTaskCardUser(client Client, payload json.RawMessage) {
	var msg UnassignTaskCardUserPayload
	if err := json.Unmarshal(payload, &msg); err != nil {
		h.SendError(client, "unassign_task_card_user", "Invalid payload")
		return
	}

	assignment, err := h.taskCardUsersUseCase.GetByID(msg.ID)
	if err != nil {
		h.SendError(client, "unassign_task_card_user", "Assignment not found")
		return
	}
	taskCard, err := h.taskCardUseCase.FindByID(context.Background(), assignment.TaskCardID)
	if err != nil {
		h.SendError(client, "unassign_task_card_user", "Task card not found")
		return
	}
	taskTab, err := h.taskTabUseCase.FindByID(taskCard.TaskTabID)
	if err != nil {
		h.SendError(client, "unassign_task_card_user", "Task tab not found")
		return
	}

	if err := h.taskCardUsersUseCase.Delete(msg.ID); err != nil {
		h.SendError(client, "unassign_task_card_user", "Failed to unassign user: "+err.Error())
		return
	}

	h.SendSuccess(client, "unassign_task_card_user", msg, map[string]interface{}{"id": msg.ID})
	h.BroadcastSuccess(h.hub, taskTab.BoardID, "unassign_task_card_user", msg, map[string]interface{}{"id": msg.ID})
}

func (h *TaskCardHandler) HandleCreateTaskCard(client Client, payload json.RawMessage) {
	var msg CreateTaskCardPayload
	if err := json.Unmarshal(payload, &msg); err != nil {
		h.SendError(client, "create_task_card", "Invalid payload")
		return
	}

	taskCardData := &taskCard.TaskCard{
		TaskTabID: msg.TaskTabID,
		Name:      msg.Name,
		Content:   msg.Content,
		Date:      msg.Date,
		Status:    false, // Default status
	}

	if err := h.taskCardUseCase.Create(context.Background(), taskCardData); err != nil {
		h.SendError(client, "create_task_card", "Failed to create task card: "+err.Error())
		return
	}

	// Fetch fresh data with preloads
	freshTaskCard, err := h.taskCardUseCase.FindByID(context.Background(), taskCardData.ID)
	if err != nil {
		h.SendError(client, "create_task_card", "Failed to refresh task card data")
		return
	}

	taskTab, err := h.taskTabUseCase.FindByID(freshTaskCard.TaskTabID)
	if err != nil {
		h.SendError(client, "create_task_card", "Task tab not found")
		return
	}

	h.SendSuccess(client, "create_task_card", msg, freshTaskCard)
	h.BroadcastSuccess(h.hub, taskTab.BoardID, "create_task_card", msg, freshTaskCard)
}
