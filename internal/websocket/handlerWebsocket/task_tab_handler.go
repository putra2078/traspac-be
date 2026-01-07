package handlerWebsocket

import (
	"encoding/json"
	"hrm-app/internal/domain/taskTab"
)

type TaskTabHandler struct {
	BaseHandler
	taskTabUseCase taskTab.UseCase
	hub            Hub
}

func NewTaskTabHandler(taskTabUseCase taskTab.UseCase, hub Hub) *TaskTabHandler {
	return &TaskTabHandler{
		taskTabUseCase: taskTabUseCase,
		hub:            hub,
	}
}

type UpdateTaskTabPayload struct {
	TaskTabID uint   `json:"task_tab_id"`
	Name      string `json:"name,omitempty"`
	Position  int    `json:"position,omitempty"`
}

func (h *TaskTabHandler) HandleUpdateTaskTab(client Client, payload json.RawMessage) {
	var msg UpdateTaskTabPayload
	if err := json.Unmarshal(payload, &msg); err != nil {
		h.SendError(client, "update_task_tab", "Invalid payload")
		return
	}

	taskTabData, err := h.taskTabUseCase.FindByID(msg.TaskTabID)
	if err != nil {
		h.SendError(client, "update_task_tab", "Task tab not found")
		return
	}

	if msg.Name != "" {
		taskTabData.Name = msg.Name
	}
	if msg.Position != 0 {
		taskTabData.Position = msg.Position
	}

	if err := h.taskTabUseCase.Update(taskTabData); err != nil {
		h.SendError(client, "update_task_tab", "Failed to update task tab")
		return
	}

	h.SendSuccess(client, "update_task_tab", msg, taskTabData)
	h.BroadcastSuccess(h.hub, taskTabData.BoardID, "update_task_tab", msg, taskTabData)
}
