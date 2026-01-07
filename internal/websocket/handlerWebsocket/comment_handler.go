package handlerWebsocket

import (
	"context"
	"encoding/json"
	"hrm-app/internal/domain/taskCard"
	"hrm-app/internal/domain/taskCardComment"
	"hrm-app/internal/domain/taskTab"
)

type CommentHandler struct {
	BaseHandler
	taskCardCommentUseCase taskCardComment.UseCase
	taskCardUseCase        taskCard.UseCase
	taskTabUseCase         taskTab.UseCase
	hub                    Hub
}

func NewCommentHandler(taskCardCommentUseCase taskCardComment.UseCase, taskCardUseCase taskCard.UseCase, taskTabUseCase taskTab.UseCase, hub Hub) *CommentHandler {
	return &CommentHandler{
		taskCardCommentUseCase: taskCardCommentUseCase,
		taskCardUseCase:        taskCardUseCase,
		taskTabUseCase:         taskTabUseCase,
		hub:                    hub,
	}
}

type CreateTaskCardCommentPayload struct {
	TaskCardID int    `json:"task_card_id"`
	Comment    string `json:"comment"`
}

type UpdateTaskCardCommentPayload struct {
	ID      int    `json:"id"`
	Comment string `json:"comment"`
}

type DeleteTaskCardCommentPayload struct {
	ID int `json:"id"`
}

func (h *CommentHandler) HandleCreateTaskCardComment(client Client, payload json.RawMessage) {
	var msg CreateTaskCardCommentPayload
	if err := json.Unmarshal(payload, &msg); err != nil {
		h.SendError(client, "create_task_card_comment", "Invalid payload")
		return
	}

	comment := &taskCardComment.TaskCardComment{
		TaskCardID: msg.TaskCardID,
		Comment:    msg.Comment,
		UserID:     client.GetUserID(),
	}

	if err := h.taskCardCommentUseCase.Create(comment); err != nil {
		h.SendError(client, "create_task_card_comment", "Failed to create comment: "+err.Error())
		return
	}

	if comment.ID < 0 {
		h.SendError(client, "create_task_card_comment", "Invalid comment ID")
		return
	}
	fullComment, err := h.taskCardCommentUseCase.FindByID(uint(comment.ID))
	if err != nil {
		h.SendError(client, "create_task_card_comment", "Failed to fetch created comment")
		return
	}

	if msg.TaskCardID < 0 {
		h.SendError(client, "create_task_card_comment", "Invalid task card ID")
		return
	}
	taskCard, err := h.taskCardUseCase.FindByID(context.Background(), uint(msg.TaskCardID))
	if err != nil {
		h.SendError(client, "create_task_card_comment", "Task card not found")
		return
	}
	taskTab, err := h.taskTabUseCase.FindByID(taskCard.TaskTabID)
	if err != nil {
		h.SendError(client, "create_task_card_comment", "Task tab not found")
		return
	}

	h.SendSuccess(client, "create_task_card_comment", msg, fullComment)
	h.BroadcastSuccess(h.hub, taskTab.BoardID, "create_task_card_comment", msg, fullComment)
}

func (h *CommentHandler) HandleUpdateTaskCardComment(client Client, payload json.RawMessage) {
	var msg UpdateTaskCardCommentPayload
	if err := json.Unmarshal(payload, &msg); err != nil {
		h.SendError(client, "update_task_card_comment", "Invalid payload")
		return
	}

	if msg.ID < 0 {
		h.SendError(client, "update_task_card_comment", "Invalid comment ID")
		return
	}
	comment, err := h.taskCardCommentUseCase.FindByID(uint(msg.ID))
	if err != nil {
		h.SendError(client, "update_task_card_comment", "Comment not found")
		return
	}

	comment.Comment = msg.Comment

	if err := h.taskCardCommentUseCase.Update(comment); err != nil {
		h.SendError(client, "update_task_card_comment", "Failed to update comment: "+err.Error())
		return
	}

	if comment.TaskCardID < 0 {
		h.SendError(client, "update_task_card_comment", "Invalid task card ID")
		return
	}
	taskCard, err := h.taskCardUseCase.FindByID(context.Background(), uint(comment.TaskCardID))
	if err != nil {
		h.SendError(client, "update_task_card_comment", "Task card not found")
		return
	}
	taskTab, err := h.taskTabUseCase.FindByID(taskCard.TaskTabID)
	if err != nil {
		h.SendError(client, "update_task_card_comment", "Task tab not found")
		return
	}

	h.SendSuccess(client, "update_task_card_comment", msg, comment)
	h.BroadcastSuccess(h.hub, taskTab.BoardID, "update_task_card_comment", msg, comment)
}

func (h *CommentHandler) HandleDeleteTaskCardComment(client Client, payload json.RawMessage) {
	var msg DeleteTaskCardCommentPayload
	if err := json.Unmarshal(payload, &msg); err != nil {
		h.SendError(client, "delete_task_card_comment", "Invalid payload")
		return
	}

	if msg.ID < 0 {
		h.SendError(client, "delete_task_card_comment", "Invalid comment ID")
		return
	}
	comment, err := h.taskCardCommentUseCase.FindByID(uint(msg.ID))
	if err != nil {
		h.SendError(client, "delete_task_card_comment", "Comment not found")
		return
	}

	if comment.TaskCardID < 0 {
		h.SendError(client, "delete_task_card_comment", "Invalid task card ID")
		return
	}
	taskCard, err := h.taskCardUseCase.FindByID(context.Background(), uint(comment.TaskCardID))
	if err != nil {
		h.SendError(client, "delete_task_card_comment", "Task card not found")
		return
	}
	taskTab, err := h.taskTabUseCase.FindByID(taskCard.TaskTabID)
	if err != nil {
		h.SendError(client, "delete_task_card_comment", "Task tab not found")
		return
	}

	if err := h.taskCardCommentUseCase.Delete(uint(msg.ID)); err != nil {
		h.SendError(client, "delete_task_card_comment", "Failed to delete comment: "+err.Error())
		return
	}

	h.SendSuccess(client, "delete_task_card_comment", msg, map[string]interface{}{"id": msg.ID})
	h.BroadcastSuccess(h.hub, taskTab.BoardID, "delete_task_card_comment", msg, map[string]interface{}{"id": msg.ID})
}
