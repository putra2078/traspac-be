package auth

import (
	"log"
	"net/http"
	"strings"
	"time"

	"hrm-app/config"
	"hrm-app/internal/domain/user"
	"hrm-app/internal/pkg/utils"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	userRepo user.Repository
	cfg      *config.Config
}

func NewHandler(repo user.Repository, cfg *config.Config) *Handler {
	return &Handler{userRepo: repo, cfg: cfg}
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "BadRequest",
			"message": err.Error(),
		})
		return
	}

	user, err := h.userRepo.FindByEmail(req.Email)
	if err != nil || user.ID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "Invalid email or password",
		})
		return
	}

	if !utils.CheckPasswordHash(req.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "Invalid email or password",
		})
		return
	}

	accessToken, refreshToken, err := utils.GenerateTokens(h.cfg, user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "InternalServerError",
			"message": "Failed to generate tokens",
		})
		return
	}

	log.Printf("[Login Debug] UserID: %d, Time: %v (Unix: %d)", user.ID, time.Now(), time.Now().Unix())
	log.Printf("[Login Debug] Generated Token: %s", accessToken)

	// Store session in Redis (Access Token)
	expMinutes := h.cfg.JWT.ExpiresInMinutes
	if expMinutes == 0 {
		expMinutes = 15
	}
	err = utils.SetSession(user.ID, accessToken, time.Duration(expMinutes)*time.Minute)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "InternalServerError",
			"message": "Failed to store session",
		})
		return
	}

	// Set Cookies
	c.SetCookie("access_token", accessToken, h.cfg.JWT.ExpiresInMinutes*60, "/", "", false, true)
	c.SetCookie("refresh_token", refreshToken, h.cfg.JWT.RefreshExpiresInDays*24*3600, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"message":      "Login successful",
		"access_token": accessToken, // Tetap return buat client yang nggak pake cookie
	})
}

func (h *Handler) Logout(c *gin.Context) {
	// Ambil token dari cookie atau header
	token := ""
	cookie, err := c.Cookie("access_token")
	if err == nil {
		token = cookie
	} else {
		authHeader := c.GetHeader("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			token = strings.TrimPrefix(authHeader, "Bearer ")
		}
	}

	if token != "" {
		claims, err := utils.ValidateToken(h.cfg, token)
		if err == nil {
			_ = utils.DeleteSession(claims.UserID, token)
		}
	}

	// Clear cookies
	c.SetCookie("access_token", "", -1, "/", "", false, true)
	c.SetCookie("refresh_token", "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func (h *Handler) RefreshToken(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized", "message": "Refresh token missing"})
		return
	}

	claims, err := utils.ValidateRefreshToken(h.cfg, refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized", "message": "Invalid refresh token"})
		return
	}

	accessToken, newRefreshToken, err := utils.GenerateTokens(h.cfg, claims.UserID, claims.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "InternalServerError", "message": "Failed to generate tokens"})
		return
	}

	// Store new session
	expMinutes := h.cfg.JWT.ExpiresInMinutes
	if expMinutes == 0 {
		expMinutes = 15
	}
	_ = utils.SetSession(claims.UserID, accessToken, time.Duration(expMinutes)*time.Minute)

	// Update cookies
	c.SetCookie("access_token", accessToken, h.cfg.JWT.ExpiresInMinutes*60, "/", "", false, true)
	c.SetCookie("refresh_token", newRefreshToken, h.cfg.JWT.RefreshExpiresInDays*24*3600, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"access_token": accessToken,
	})
}
