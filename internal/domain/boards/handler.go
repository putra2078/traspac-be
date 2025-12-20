package boards

import (
	"hrm-app/internal/response"
	"strconv"
	"net/http"

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
	if err := h.usecase.Create(&boards); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, "Board created successfully")
}

func (h *Handler) GetAllBoard(c *gin.Context) {
	boards, err := h.usecase.FindAll()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, boards)
}

func (h *Handler) GetBoardByID(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid board ID")
		return
	}
	boards, err := h.usecase.FindByID(uint(idInt))
	if err != nil {
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
	boards.ID = uint(idInt)
	if err := h.usecase.Update(&boards); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, "Board updated successfully")
}

func (h *Handler) DeleteBoard(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid board ID")
		return
	}
	if err := h.usecase.Delete(uint(idInt)); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, "Board deleted successfully")
}

func (h *Handler) GetByWorkspaceID(c *gin.Context) {
	workspaceID := c.Param("workspace_id")
	workspaceIDInt, err := strconv.Atoi(workspaceID)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid workspace ID")
		return
	}
	boards, err := h.usecase.FindByWorkspaceID(uint(workspaceIDInt))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, boards)
}
