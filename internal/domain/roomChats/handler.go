package room_chats

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

func (h *Handler) Create(c *gin.Context) {
	var roomChats RoomsChats
	if err := c.ShouldBindJSON(&roomChats); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	userID := c.MustGet("user_id").(uint)
	roomChats.CreatedBy = userID

	if err := h.usecase.Create(&roomChats); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, roomChats)
}

func (h *Handler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid ID parameter")
		return
	}
	roomChats, err := h.usecase.GetByID(uint(id))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, roomChats)
}

func (h *Handler) GetByWorkspaceID(c *gin.Context) {
	workspaceIDParam := c.Param("workspace_id")
	workspaceID, err := strconv.ParseUint(workspaceIDParam, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid workspace ID parameter")
		return
	}
	roomChats, err := h.usecase.GetByWorkspaceID(uint(workspaceID))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, roomChats)
}

func (h *Handler) GetAll(c *gin.Context) {
	roomChats, err := h.usecase.GetAll()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, roomChats)
}

func (h *Handler) Update(c *gin.Context) {
	var roomChats RoomsChats
	if err := c.ShouldBindJSON(&roomChats); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.usecase.Update(&roomChats); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, roomChats)
}

func (h *Handler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid ID parameter")
		return
	}
	if err := h.usecase.Delete(uint(id)); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.DeleteSuccess(c, "room chat deleted successfully")
}

func (h *Handler) UploadAttachment(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.Error(c, http.StatusBadRequest, "file is required")
		return
	}

	url, err := h.usecase.UploadAttachment(c, file)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, gin.H{
		"url": url,
	})
}
