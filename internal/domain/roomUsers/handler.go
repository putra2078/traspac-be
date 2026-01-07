package roomUsers

import (
	"net/http"
	"strconv"

	"hrm-app/internal/response"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	usecase UseCase
}

func NewHandler(usecase UseCase) *Handler {
	return &Handler{usecase: usecase}
}

func (h *Handler) GetUsersByRoom(c *gin.Context) {
	roomIDParam := c.Param("room_id")
	roomID, err := strconv.ParseUint(roomIDParam, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid room ID parameter")
		return
	}

	users, err := h.usecase.GetUsersByRoom(uint(roomID))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, users)
}

func (h *Handler) Join(c *gin.Context) {
	var body struct {
		RoomID uint `json:"room_id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	userID := c.MustGet("user_id").(uint)
	if err := h.usecase.Join(body.RoomID, userID); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, "success join room")
}
