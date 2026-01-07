package handlerWebsocket

import (
	"encoding/json"
	"hrm-app/internal/domain/workspacesUsers"
)

type WorkspaceHandler struct {
	BaseHandler
	workspacesUsersUseCase workspacesUsers.UseCase
	hub                    Hub
}

func NewWorkspaceHandler(workspacesUsersUseCase workspacesUsers.UseCase, hub Hub) *WorkspaceHandler {
	return &WorkspaceHandler{
		workspacesUsersUseCase: workspacesUsersUseCase,
		hub:                    hub,
	}
}

type AssignWorkspaceUserPayload struct {
	WorkspaceID uint `json:"workspace_id"`
	UserID      uint `json:"user_id"`
}

type UnassignWorkspaceUserPayload struct {
	ID uint `json:"id"`
}

func (h *WorkspaceHandler) HandleAssignWorkspaceUser(client Client, payload json.RawMessage) {
	var msg AssignWorkspaceUserPayload
	if err := json.Unmarshal(payload, &msg); err != nil {
		h.SendError(client, "assign_workspace_user", "Invalid payload")
		return
	}

	assignment := &workspacesUsers.WorkspacesUsers{
		WorkspaceID: msg.WorkspaceID,
		UserID:      msg.UserID,
	}

	if err := h.workspacesUsersUseCase.Create(assignment, client.GetUserID()); err != nil {
		h.SendError(client, "assign_workspace_user", "Failed to assign user to workspace: "+err.Error())
		return
	}

	fullAssignment, err := h.workspacesUsersUseCase.GetByID(assignment.ID)
	if err != nil {
		h.SendError(client, "assign_workspace_user", "Failed to fetch assignment details")
		return
	}

	h.SendSuccess(client, "assign_workspace_user", msg, fullAssignment)
	h.BroadcastGlobalSuccess(h.hub, "assign_workspace_user", msg, fullAssignment)
}

func (h *WorkspaceHandler) HandleUnassignWorkspaceUser(client Client, payload json.RawMessage) {
	var msg UnassignWorkspaceUserPayload
	if err := json.Unmarshal(payload, &msg); err != nil {
		h.SendError(client, "unassign_workspace_user", "Invalid payload")
		return
	}

	if err := h.workspacesUsersUseCase.Delete(msg.ID); err != nil {
		h.SendError(client, "unassign_workspace_user", "Failed to unassign user from workspace: "+err.Error())
		return
	}

	h.SendSuccess(client, "unassign_workspace_user", msg, map[string]interface{}{"id": msg.ID})
	h.BroadcastGlobalSuccess(h.hub, "unassign_workspace_user", msg, map[string]interface{}{"id": msg.ID})
}
