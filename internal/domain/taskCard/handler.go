package taskCard

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
	var taskCard TaskCard
	if err := c.ShouldBindJSON(&taskCard); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.usecase.Create(&taskCard); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, taskCard)
}

func (h *Handler) GetAll(c *gin.Context) {
	taskCards, err := h.usecase.FindAll()
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, taskCards)
}

func (h *Handler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	taskCard, err := h.usecase.FindByID(uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error())
		return
	}

	response.Success(c, taskCard)
}

func (h *Handler) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	var taskCard TaskCard
	if err := c.ShouldBindJSON(&taskCard); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	taskCard.ID = uint(id)

	if err := h.usecase.Update(&taskCard); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, taskCard)
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

func (h *Handler) GetByTaskTabID(c *gin.Context) {
	taskTabIDParam := c.Param("task_tab_id")
	taskTabID, err := strconv.Atoi(taskTabIDParam)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	taskCards, err := h.usecase.FindByTaskTabID(uint(taskTabID))
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, taskCards)
}
