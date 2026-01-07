package labels

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
	var label TaskCardLabel
	if err := c.ShouldBindJSON(&label); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.usecase.Create(&label); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, label)
}

func (h *Handler) GetAll(c *gin.Context) {
	labels, err := h.usecase.FindAll()
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, labels)
}

func (h *Handler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil || id < 0 {
		response.Error(c, http.StatusBadRequest, "Invalid ID parameter")
		return
	}

	label, err := h.usecase.FindByID(uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error())
		return
	}

	response.Success(c, label)
}

func (h *Handler) GetByTaskCardID(c *gin.Context) {
	taskCardIDParam := c.Param("task_card_id")
	taskCardID, err := strconv.Atoi(taskCardIDParam)
	if err != nil || taskCardID < 0 {
		response.Error(c, http.StatusBadRequest, "Invalid task card ID parameter")
		return
	}

	labels, err := h.usecase.FindByTaskCardID(uint(taskCardID))
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error())
		return
	}

	response.Success(c, labels)
}

func (h *Handler) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil || id < 0 {
		response.Error(c, http.StatusBadRequest, "Invalid ID parameter")
		return
	}

	var label TaskCardLabel
	if err := c.ShouldBindJSON(&label); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	label.ID = uint(id)

	if err := h.usecase.Update(&label); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, label)
}

func (h *Handler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil || id < 0 {
		response.Error(c, http.StatusBadRequest, "Invalid ID parameter")
		return
	}

	if err := h.usecase.Delete(uint(id)); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, nil)
}
