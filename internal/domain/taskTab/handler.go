package taskTab

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
	var taskTab TaskTab
	if err := c.ShouldBindJSON(&taskTab); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.usecase.Create(&taskTab); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, taskTab)
}

func (h *Handler) GetAll(c *gin.Context) {
	taskTabs, err := h.usecase.FindAll()
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, taskTabs)
}

func (h *Handler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil || id < 0 {
		response.Error(c, http.StatusBadRequest, "Invalid ID parameter")
		return
	}

	taskTab, err := h.usecase.FindByID(uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error())
		return
	}

	response.Success(c, taskTab)
}

func (h *Handler) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil || id < 0 {
		response.Error(c, http.StatusBadRequest, "Invalid ID parameter")
		return
	}

	var taskTab TaskTab
	if err := c.ShouldBindJSON(&taskTab); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	taskTab.ID = uint(id)

	if err := h.usecase.Update(&taskTab); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, taskTab)
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
