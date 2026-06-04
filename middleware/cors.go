package middleware

import (
	"github.com/gin-gonic/gin"
)

// Middleware CORS biar Flutter bisa akses API
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Allow semua origin
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		
		// Allow method apa aja
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		
		// Allow header apa aja
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Accept")
		
		// Cache preflight request 24 jam
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		
		// Kalau OPTIONS request, langsung return
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	}
}
