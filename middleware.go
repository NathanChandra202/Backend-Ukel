package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func CekToken(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	var tokenString string

	if authHeader != "" {
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"pesan": MsgTokenFormatSalah})
			c.Abort()
			return
		}
		tokenString = parts[1]
	} else {
		tokenString = c.Query("token")
	}

	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"pesan": MsgTokenTidakAda})
		c.Abort()
		return
	}

	secretKey := []byte(os.Getenv("JWT_SECRET"))
	if len(secretKey) == 0 {
		secretKey = []byte("rahasiabanget")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("metode token tidak didukung")
		}
		return secretKey, nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"pesan": MsgTokenTidakValid})
		c.Abort()
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"pesan": MsgTokenBacaGagal})
		c.Abort()
		return
	}

	idClaim, ok := claims["id"].(float64)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"pesan": MsgTokenBacaGagal})
		c.Abort()
		return
	}

	c.Set("siswaId", int(idClaim))
	c.Next()
}
