package workspaces

import (
	"net/http"
	"strconv"

	"hrm-app/internal/response"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	usecase UseCase
}

func NewHandler(u UseCase) *Handler {
	return &Handler{usecase: u}
}

func (h *Handler) Create(c *gin.Context) {
	var workspace Workspace
	if err := c.ShouldBindJSON(&workspace); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	// Get user_id from context (set by AuthMiddleware)
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	workspace.CreatedBy = userID.(uint)

	if err := h.usecase.Create(&workspace); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, workspace)
}

func (h *Handler) GetAll(c *gin.Context) {
	workspaces, err := h.usecase.GetAll()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, workspaces)
}

func (h *Handler) GetByID(c *gin.Context) {
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

	workspace, err := h.usecase.GetByID(uint(id), userID.(uint))
	if err != nil {
		response.Error(c, http.StatusForbidden, err.Error())
		return
	}

	response.Success(c, workspace)
}

func (h *Handler) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil || id < 0 {
		response.Error(c, http.StatusBadRequest, "Invalid id parameter")
		return
	}

	var workspace Workspace
	if err := c.ShouldBindJSON(&workspace); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	workspace.ID = uint(id)

	if err := h.usecase.Update(&workspace); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, workspace)
}

func (h *Handler) GetByUserID(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	workspaces, err := h.usecase.GetByUserID(userID.(uint))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, workspaces)
}

func (h *Handler) GetGuestWorkspaces(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	workspaces, err := h.usecase.GetGuestWorkspaces(userID.(uint))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, workspaces)
}

func (h *Handler) Delete(c *gin.Context) {
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

	if err := h.usecase.DeleteByID(uint(id), userID.(uint)); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.DeleteSuccess(c, "Workspace deleted successfully")
}
