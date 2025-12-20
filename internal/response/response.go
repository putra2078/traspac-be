package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Error(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{"error": message})
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func DeleteSuccess(c *gin.Context, message string) {
	c.JSON(http.StatusOK, gin.H{"message": message})
}
