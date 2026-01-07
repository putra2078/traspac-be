package boards

import (
	"hrm-app/internal/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	usecase UseCase
}

func NewHandler(usecase UseCase) Handler {
	return Handler{usecase: usecase}
}

func (h *Handler) CreateBoard(c *gin.Context) {
	var boards Boards
	if err := c.ShouldBindJSON(&boards); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	// Get user_id from context (set by AuthMiddleware)
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	boards.CreatedBy = userID.(uint)
	ctx := c.Request.Context()

	if err := h.usecase.Create(ctx, &boards); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, "Board created successfully")
}

func (h *Handler) GetAllBoard(c *gin.Context) {
	ctx := c.Request.Context()
	boards, err := h.usecase.FindAll(ctx)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, boards)
}

func (h *Handler) GetBoardByID(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil || idInt < 0 {
		response.Error(c, http.StatusBadRequest, "Invalid board ID")
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	ctx := c.Request.Context()
	boards, err := h.usecase.FindByID(ctx, uint(idInt), userID.(uint))
	if err != nil {
		if err.Error() == "unauthorized" {
			response.Error(c, http.StatusUnauthorized, "Unauthorized")
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, boards)
}

func (h *Handler) UpdateBoard(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid board ID")
		return
	}
	var boards Boards
	if err := c.ShouldBindJSON(&boards); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if idInt < 0 {
		response.Error(c, http.StatusBadRequest, "Invalid board ID")
		return
	}

	boards.ID = uint(idInt)

	ctx := c.Request.Context()
	if err := h.usecase.Update(ctx, &boards); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, "Board updated successfully")
}

func (h *Handler) DeleteBoard(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil || idInt < 0 {
		response.Error(c, http.StatusBadRequest, "Invalid board ID")
		return
	}
	ctx := c.Request.Context()
	if err := h.usecase.Delete(ctx, uint(idInt)); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, "Board deleted successfully")
}

func (h *Handler) GetByWorkspaceID(c *gin.Context) {
	workspaceID := c.Param("workspace_id")
	workspaceIDInt, err := strconv.Atoi(workspaceID)
	if err != nil || workspaceIDInt < 0 {
		response.Error(c, http.StatusBadRequest, "Invalid workspace ID")
		return
	}
	ctx := c.Request.Context()
	boards, err := h.usecase.FindByWorkspaceID(ctx, uint(workspaceIDInt))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, boards)
}

func (h *Handler) GetByUserID(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	ctx := c.Request.Context()
	boards, err := h.usecase.FindByUserID(ctx, userID.(uint))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, boards)
}

func (h *Handler) GetBoardTabs(c *gin.Context) {
	boardIDStr := c.Param("id")
	boardID, err := strconv.ParseUint(boardIDStr, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid board ID")
		return
	}

	ctx := c.Request.Context()
	tabs, err := h.usecase.GetTabsByBoardID(ctx, uint(boardID))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, tabs)
}

func (h *Handler) GetTabCards(c *gin.Context) {
	tabIDStr := c.Param("tab_id")
	tabID, err := strconv.ParseUint(tabIDStr, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid task tab ID")
		return
	}

	limitStr := c.DefaultQuery("limit", "50")
	pageStr := c.DefaultQuery("page", "1")

	limit, _ := strconv.Atoi(limitStr)
	page, _ := strconv.Atoi(pageStr)
	if limit <= 0 {
		limit = 50
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	ctx := c.Request.Context()
	cards, err := h.usecase.GetCardsByTaskTabID(ctx, uint(tabID), limit, offset)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, cards)
}
