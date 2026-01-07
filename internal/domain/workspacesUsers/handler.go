package workspacesUsers

import (
	"hrm-app/internal/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	usecase UseCase
}

func NewHandler(u UseCase) *Handler {
	return &Handler{usecase: u}
}

func (h *Handler) Create(c *gin.Context) {
	var workspacesUsers WorkspacesUsers
	if err := c.ShouldBindJSON(&workspacesUsers); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	// Get user_id from context (set by AuthMiddleware)
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	if err := h.usecase.Create(&workspacesUsers, userID.(uint)); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, workspacesUsers)
}

func (h *Handler) GetByWorkspaceID(c *gin.Context) {
	workspaceID := c.Param("workspace_id")
	workspaceIDInt, err := strconv.Atoi(workspaceID)
	if err != nil || workspaceIDInt < 0 {
		response.Error(c, http.StatusBadRequest, "Invalid workspace_id parameter")
		return
	}

	workspacesUsers, err := h.usecase.GetByWorkspaceID(uint(workspaceIDInt))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, workspacesUsers)
}

func (h *Handler) GetByUserID(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	workspacesUsers, err := h.usecase.GetByUserID(userID.(uint))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, workspacesUsers)
}

func (h *Handler) GetByID(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil || idInt < 0 {
		response.Error(c, http.StatusBadRequest, "Invalid id parameter")
		return
	}

	workspacesUsers, err := h.usecase.GetByID(uint(idInt))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, workspacesUsers)
}

func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil || idInt < 0 {
		response.Error(c, http.StatusBadRequest, "Invalid id parameter")
		return
	}

	if err := h.usecase.Delete(uint(idInt)); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to delete workspacesUsers")
		return
	}

	response.DeleteSuccess(c, "WorkspacesUsers deleted successfully")
}

func (h *Handler) Update(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil || idInt < 0 {
		response.Error(c, http.StatusBadRequest, "Invalid id parameter")
		return
	}

	var workspacesUsers WorkspacesUsers
	if err := c.ShouldBindJSON(&workspacesUsers); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	workspacesUsers.ID = uint(idInt)

	if err := h.usecase.Update(&workspacesUsers); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, workspacesUsers)
}

func (h *Handler) Join(c *gin.Context) {
	var req struct {
		Token string `json:"token"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	// Get user_id from context (set by AuthMiddleware)
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	if err := h.usecase.Join(userID.(uint), req.Token); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "Joined workspace successfully"})
}

func (h *Handler) GenerateJoinToken(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil || id < 0 {
		response.Error(c, http.StatusBadRequest, "Invalid id parameter")
		return
	}

	// Get user_id from context (set by AuthMiddleware)
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	token, err := h.usecase.GenerateJoinToken(uint(id), userID.(uint))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, gin.H{"token": token})
}
