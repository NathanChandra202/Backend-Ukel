package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Middleware untuk cek JWT token
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ambil token dari header "Authorization: Bearer <token>"
		authHeader := c.GetHeader("Authorization")
		var tokenString string

		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				c.JSON(http.StatusUnauthorized, gin.H{"pesan": "Format token salah. Pakai: Bearer <token>"})
				c.Abort()
				return
			}
			tokenString = parts[1]
		} else {
			// Alternatif: ambil dari query parameter ?token=xxx
			tokenString = c.Query("token")
		}

		// Kalau tidak ada token
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"pesan": "Token tidak ada. Silakan login dulu."})
			c.Abort()
			return
		}

		// Ambil secret key dari .env
		secretKey := os.Getenv("JWT_SECRET")
		if secretKey == "" {
			secretKey = "rahasiabanget"
		}

		// Parse & validasi token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})

		// Kalau token tidak valid atau expired
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"pesan": "Token tidak valid atau sudah kedaluwarsa. Login lagi."})
			c.Abort()
			return
		}

		// Ambil data dari token (id user)
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"pesan": "Token tidak bisa dibaca."})
			c.Abort()
			return
		}

		userID, ok := claims["id"].(float64)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"pesan": "ID user tidak ada di token."})
			c.Abort()
			return
		}

		// Simpan ID user ke context, bisa dipakai di controller
		c.Set("siswaId", uint(userID))
		c.Next()
	}
}
