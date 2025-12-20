package servers

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

func(h *Handler) Create(c *gin.Context) {
	var server Server
	if err := c.ShouldBindJSON(&server); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.usecase.Create(&server); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, server)
}

func (h *Handler) GetAll(c *gin.Context) {
	servers, err := h.usecase.GetAll()
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c,servers)
}

func (h *Handler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	server, err := h.usecase.GetByID(uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error())
		return
	}

	response.Success(c, server)
}