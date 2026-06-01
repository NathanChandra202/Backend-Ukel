package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// fungsi buat cek token jwt pas ada request masuk, kalau nggak valid tolak
func CekToken(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	var tokenString string

	if authHeader != "" {
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"pesan": "Format token salah"})
			c.Abort()
			return
		}
		tokenString = parts[1]
	} else {
		tokenString = c.Query("token")
	}

	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"pesan": "Token tidak ada, harap login"})
		c.Abort()
		return
	}

	secretKey := []byte(os.Getenv("JWT_SECRET"))

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return secretKey, nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"pesan": "Token tidak valid atau sudah expired"})
		c.Abort()
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"pesan": "Gagal membaca isi token"})
		c.Abort()
		return
	}

	siswaID := int(claims["id"].(float64))
	c.Set("siswaId", siswaID)

	c.Next()
}
