package contact

import (
	"net/http"

	"hrm-app/internal/response"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	usecase UseCase
	bucket  string
}

func NewHandler(u UseCase, bucket string) *Handler {
	return &Handler{
		usecase: u,
		bucket:  bucket,
	}
}

// GetMyContact gets the contact of the currently logged-in user
func (h *Handler) GetMyContact(c *gin.Context) {
	// Get user_id from context (set by AuthMiddleware)
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	ctx := c.Request.Context()
	contact, err := h.usecase.GetByUserID(ctx, userID.(uint), h.bucket)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	if contact == nil {
		response.Error(c, http.StatusNotFound, "Contact not found")
		return
	}

	response.Success(c, contact)
}

// UpdateMyContact updates the contact of the currently logged-in user
func (h *Handler) UpdateMyContact(c *gin.Context) {
	// Get user_id from context (set by AuthMiddleware)
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var contact Contact
	if err := c.ShouldBindJSON(&contact); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	// Set the user_id to ensure user can only update their own contact
	contact.UserID = userID.(uint)

	ctx := c.Request.Context()
	if err := h.usecase.Update(ctx, &contact); err != nil {
		if err.Error() == "contact not found" {
			response.Error(c, http.StatusNotFound, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Fetch updated contact with public URL
	updatedContact, err := h.usecase.GetByUserID(ctx, userID.(uint), h.bucket)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, updatedContact)
}
