package taskCardComment

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

func (h *Handler) CreateTaskCardComment(c *gin.Context) {
	var taskCardComment TaskCardComment
	if err := c.ShouldBindJSON(&taskCardComment); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.usecase.Create(&taskCardComment); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, taskCardComment)
}

func (h *Handler) GetAllTaskCardComment(c *gin.Context) {
	taskCardComments, err := h.usecase.FindAll()
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, taskCardComments)
}

func (h *Handler) GetTaskCardCommentByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	taskCardComment, err := h.usecase.FindByID(uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error())
		return
	}

	response.Success(c, taskCardComment)
}

func (h *Handler) GetTaskCardCommentByTaskCardID(c *gin.Context) {
	taskCardIDParam := c.Param("task_card_id")
	taskCardID, err := strconv.Atoi(taskCardIDParam)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	taskCardComments, err := h.usecase.FindByTaskCardID(uint(taskCardID))
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, taskCardComments)
}

func (h *Handler) UpdateTaskCardComment(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	var taskCardComment TaskCardComment
	if err := c.ShouldBindJSON(&taskCardComment); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	taskCardComment.ID = id

	if err := h.usecase.Update(&taskCardComment); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, taskCardComment)
}

func (h *Handler) DeleteTaskCardComment(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.usecase.Delete(uint(id)); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, nil)
}
