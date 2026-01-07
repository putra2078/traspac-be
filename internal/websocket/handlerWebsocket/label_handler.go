package handlerWebsocket

import (
	"context"
	"encoding/json"
	"hrm-app/internal/domain/labels"
	"hrm-app/internal/domain/taskCard"
	"hrm-app/internal/domain/taskTab"
)

type LabelHandler struct {
	BaseHandler
	labelsUseCase   labels.UseCase
	taskCardUseCase taskCard.UseCase
	taskTabUseCase  taskTab.UseCase
	hub             Hub
}

func NewLabelHandler(labelsUseCase labels.UseCase, taskCardUseCase taskCard.UseCase, taskTabUseCase taskTab.UseCase, hub Hub) *LabelHandler {
	return &LabelHandler{
		labelsUseCase:   labelsUseCase,
		taskCardUseCase: taskCardUseCase,
		taskTabUseCase:  taskTabUseCase,
		hub:             hub,
	}
}

type CreateLabelPayload struct {
	TaskCardID uint   `json:"task_card_id"`
	Title      string `json:"title"`
	Color      string `json:"color"`
}

type UpdateLabelPayload struct {
	ID    uint   `json:"id"`
	Title string `json:"title,omitempty"`
	Color string `json:"color,omitempty"`
}

type DeleteLabelPayload struct {
	ID uint `json:"id"`
}

func (h *LabelHandler) HandleCreateLabel(client Client, payload json.RawMessage) {
	var msg CreateLabelPayload
	if err := json.Unmarshal(payload, &msg); err != nil {
		h.SendError(client, "create_label", "Invalid payload")
		return
	}

	label := &labels.TaskCardLabel{
		TaskCardID: msg.TaskCardID,
		Title:      msg.Title,
		Color:      msg.Color,
	}

	if err := h.labelsUseCase.Create(label); err != nil {
		h.SendError(client, "create_label", "Failed to create label: "+err.Error())
		return
	}

	taskCard, err := h.taskCardUseCase.FindByID(context.Background(), msg.TaskCardID)
	if err != nil {
		h.SendError(client, "create_label", "Task card not found")
		return
	}
	taskTab, err := h.taskTabUseCase.FindByID(taskCard.TaskTabID)
	if err != nil {
		h.SendError(client, "create_label", "Task tab not found")
		return
	}

	h.SendSuccess(client, "create_label", msg, label)
	h.BroadcastSuccess(h.hub, taskTab.BoardID, "create_label", msg, label)
}

func (h *LabelHandler) HandleUpdateLabel(client Client, payload json.RawMessage) {
	var msg UpdateLabelPayload
	if err := json.Unmarshal(payload, &msg); err != nil {
		h.SendError(client, "update_label", "Invalid payload")
		return
	}

	label, err := h.labelsUseCase.FindByID(msg.ID)
	if err != nil {
		h.SendError(client, "update_label", "Label not found")
		return
	}

	if msg.Title != "" {
		label.Title = msg.Title
	}
	if msg.Color != "" {
		label.Color = msg.Color
	}

	if err := h.labelsUseCase.Update(label); err != nil {
		h.SendError(client, "update_label", "Failed to update label: "+err.Error())
		return
	}

	taskCard, err := h.taskCardUseCase.FindByID(context.Background(), label.TaskCardID)
	if err != nil {
		h.SendError(client, "update_label", "Task card not found")
		return
	}
	taskTab, err := h.taskTabUseCase.FindByID(taskCard.TaskTabID)
	if err != nil {
		h.SendError(client, "update_label", "Task tab not found")
		return
	}

	h.SendSuccess(client, "update_label", msg, label)
	h.BroadcastSuccess(h.hub, taskTab.BoardID, "update_label", msg, label)
}

func (h *LabelHandler) HandleDeleteLabel(client Client, payload json.RawMessage) {
	var msg DeleteLabelPayload
	if err := json.Unmarshal(payload, &msg); err != nil {
		h.SendError(client, "delete_label", "Invalid payload")
		return
	}

	label, err := h.labelsUseCase.FindByID(msg.ID)
	if err != nil {
		h.SendError(client, "delete_label", "Label not found")
		return
	}
	taskCard, err := h.taskCardUseCase.FindByID(context.Background(), label.TaskCardID)
	if err != nil {
		h.SendError(client, "delete_label", "Task card not found")
		return
	}
	taskTab, err := h.taskTabUseCase.FindByID(taskCard.TaskTabID)
	if err != nil {
		h.SendError(client, "delete_label", "Task tab not found")
		return
	}

	if err := h.labelsUseCase.Delete(msg.ID); err != nil {
		h.SendError(client, "delete_label", "Failed to delete label: "+err.Error())
		return
	}

	h.SendSuccess(client, "delete_label", msg, map[string]interface{}{"id": msg.ID})
	h.BroadcastSuccess(h.hub, taskTab.BoardID, "delete_label", msg, map[string]interface{}{"id": msg.ID})
}
