package boardsUsers

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
	var boardUsers BoardsUsers
	if err := c.ShouldBindJSON(&boardUsers); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	// Get user_id from context (set by AuthMiddleware)
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	if err := h.usecase.Create(&boardUsers, userID.(uint)); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, boardUsers)
}

func (h *Handler) GetByBoardID(c *gin.Context) {
	boardIDParam := c.Param("board_id")
	boardID, err := strconv.Atoi(boardIDParam)
	if err != nil || boardID < 0 {
		response.Error(c, http.StatusBadRequest, "Invalid board_id parameter")
		return
	}

	boardUsers, err := h.usecase.GetByBoardID(uint(boardID))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, boardUsers)
}

func (h *Handler) GetByUserID(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	boardUsers, err := h.usecase.GetByUserID(userID.(uint))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, boardUsers)
}

func (h *Handler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil || id < 0 {
		response.Error(c, http.StatusBadRequest, "Invalid id parameter")
		return
	}

	boardUsers, err := h.usecase.GetByID(uint(id))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, boardUsers)
}

func (h *Handler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil || id < 0 {
		response.Error(c, http.StatusBadRequest, "Invalid id parameter")
		return
	}

	if err := h.usecase.Delete(uint(id)); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to delete board user")
		return
	}

	response.DeleteSuccess(c, "Board user deleted successfully")
}

func (h *Handler) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil || id < 0 {
		response.Error(c, http.StatusBadRequest, "Invalid id parameter")
		return
	}

	var boardUsers BoardsUsers
	if err := c.ShouldBindJSON(&boardUsers); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	boardUsers.ID = uint(id)

	if err := h.usecase.Update(&boardUsers); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, boardUsers)
}

func (h *Handler) Join(c *gin.Context) {
	var req struct {
		Token string `json:"token"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	// Get user_id from context (set by AuthMiddleware)
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	if err := h.usecase.Join(userID.(uint), req.Token); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "Joined board successfully"})
}

func (h *Handler) GenerateJoinToken(c *gin.Context) {
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

	token, err := h.usecase.GenerateJoinToken(uint(id), userID.(uint))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, gin.H{"token": token})
}
