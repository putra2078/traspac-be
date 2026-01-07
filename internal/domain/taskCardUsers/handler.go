package taskCardUsers

import (
	"hrm-app/internal/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	usecase UseCase
}

func NewHandler(usecase UseCase) *Handler {
	return &Handler{usecase: usecase}
}

func (h *Handler) CreateTaskCardUser(c *gin.Context) {
	var payload TaskCardUsers

	if err := c.ShouldBindJSON(&payload); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.usecase.Create(&payload); err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to create task card user")
		return
	}

	response.Success(c, payload)

}

func (h *Handler) GetTaskCardUserByTaskCardID(c *gin.Context) {
	taskCardIDParam := c.Param("task_card_id")
	taskCardID, err := strconv.Atoi(taskCardIDParam)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid task card id")
		return
	}

	taskCardUsers, err := h.usecase.GetByTaskCardID(uint(taskCardID))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to get task card users")
		return
	}

	response.Success(c, taskCardUsers)
}

func (h *Handler) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	var taskCardUsers TaskCardUsers
	if err := c.ShouldBindJSON(&taskCardUsers); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	taskCardUsers.ID = uint(id)

	if err := h.usecase.Update(&taskCardUsers); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, taskCardUsers)
}

func (h *Handler) Delete(c *gin.Context) {
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
