package contact

// import (
// 	"net/http"
// 	"strconv"

// 	"github.com/gin-gonic/gin"
// )

// type Handler struct {
// 	usecase UseCase
// }

// func NewHandler(u UseCase) *Handler {
// 	return &Handler{usecase: u}
// }

// func (h *Handler) Register(c *gin.Context) {
// 	var contact Contact
// 	if err := c.ShouldBindJSON(&contact); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	if err := h.usecase.Register(&contact); err != nil {
// 		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Contact registered successfully"})
// }

// func (h *Handler) GetAll(c *gin.Context) {
// 	data, err := h.usecase.GetAll()
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, data)
// }

// func (h *Handler) GetByID(c *gin.Context) {
// 	id, _ :=  strconv.Atoi(c.Param("id"))
// 	data ,err := h.usecase.GetByID(uint(id))

// 	if err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Contact not found"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, data)
// }

// func (h *Handler) Delete(c *gin.Context) {
// 	id, _ := strconv.Atoi(c.Param("id"))
// 	err := h.usecase.DeleteByID(uint(id))
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Contact deleted successfully"})
// }
