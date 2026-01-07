package middleware

import (
	"net/http"
	"strings"
	"time"

	"hrm-app/config"
	"hrm-app/internal/pkg/utils"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := ""

		// 1. Coba ambil dari Cookie
		if cookie, err := c.Cookie("access_token"); err == nil {
			tokenStr = cookie
		}

		// 2. Jika tidak ada di cookie, coba ambil dari Authorization header
		if tokenStr == "" {
			authHeader := c.GetHeader("Authorization")
			if authHeader != "" {
				parts := strings.SplitN(authHeader, " ", 2)
				if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
					tokenStr = parts[1]
				}
			}
		}

		// 3. Jika masih tidak ada, coba ambil dari query parameter (untuk WebSocket)
		if tokenStr == "" {
			tokenStr = c.Query("token")
		}

		if tokenStr == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "Authentication required",
			})
			return
		}

		// 3. Validasi JWT
		claims, err := utils.ValidateToken(cfg, tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "Invalid or expired token",
			})
			return
		}

		// 4. Validasi Session di Redis
		_, err = utils.GetSession(claims.UserID, tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "Session expired or invalid",
			})
			return
		}

		// Simpan seluruh claims ke context biar bisa diakses langsung di handler lain
		c.Set("claims", claims)
		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)

		// 5. Extend Session (Activity Check / Sliding Window)
		// Kita perpanjang session setiap kali ada request agar tidak expired selama user aktif
		expMinutes := cfg.JWT.ExpiresInMinutes
		if expMinutes == 0 {
			expMinutes = 15
		}
		utils.ExtendSession(claims.UserID, tokenStr, time.Duration(expMinutes)*time.Minute)

		// 6. Perbarui Cookie Access Token agar browser tidak menghapus cookie sebelum Redis expired
		c.SetCookie("access_token", tokenStr, expMinutes*60, "/", "", false, true)

		c.Next()
	}
}
