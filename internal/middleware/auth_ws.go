package middleware

import (
	"log"
	"net/http"
	"strings"
	"time"

	"hrm-app/config"
	"hrm-app/internal/pkg/utils"

	"github.com/gin-gonic/gin"
)

// AuthMiddlewareWS is a WebSocket-specific authentication middleware
// that only uses query parameter tokens and avoids cookie operations
// to prevent handshake failures in production environments
func AuthMiddlewareWS(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := ""

		// 1. Try Query Parameter (Standard for WebSockets)
		if qToken := c.Query("token"); qToken != "" {
			tokenStr = qToken
		}

		// 2. Try Authorization Header (Fallback)
		if tokenStr == "" {
			authHeader := c.GetHeader("Authorization")
			if authHeader != "" {
				parts := strings.SplitN(authHeader, " ", 2)
				if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
					tokenStr = parts[1]
				}
			}
		}

		// 3. Try Cookie (Browser Fallback)
		if tokenStr == "" {
			if cookie, err := c.Cookie("access_token"); err == nil {
				tokenStr = cookie
			}
		}

		// Handle optional "Bearer " prefix in query param if present
		tokenStr = strings.TrimSpace(tokenStr)
		if strings.HasPrefix(strings.ToLower(tokenStr), "bearer ") {
			tokenStr = tokenStr[7:]
		}

		if tokenStr == "" {
			log.Println("[WS Auth] Rejecting: Token missing in all sources (Query, Header, Cookie)")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "Authentication token required",
			})
			return
		}

		log.Printf("[WS Auth Debug] Current Server Time: %v (Unix: %d)", time.Now(), time.Now().Unix())

		// 2. Validate JWT
		claims, err := utils.ValidateToken(cfg, tokenStr)
		if err != nil {
			log.Printf("[WS Auth] Rejecting: JWT validation failed (Token: %s...): %v", tokenStr[:10], err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "Invalid or expired token",
			})
			return
		}

		log.Printf("[WS Auth Debug] Claims: UserID=%d, Exp=%v, Sub=%s", claims.UserID, claims.ExpiresAt, claims.Subject)

		// 3. Validate session in Redis
		_, err = utils.GetSession(claims.UserID, tokenStr)
		if err != nil {
			log.Printf("[WS Auth] Session validation failed for UserID %d: %v", claims.UserID, err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "Session expired or invalid",
			})
			return
		}

		// 4. Store in context for WebSocket handler
		c.Set("claims", claims)
		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)

		log.Printf("[WS Auth] Successful handshake for UserID %d", claims.UserID)

		// 5. Extend Session (Activity Check)
		expMinutes := cfg.JWT.ExpiresInMinutes
		if expMinutes == 0 {
			expMinutes = 60 // Default can follow config or 60
		}
		_ = utils.ExtendSession(claims.UserID, tokenStr, time.Duration(expMinutes)*time.Minute)

		c.Next()
	}
}
