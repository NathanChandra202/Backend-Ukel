package utils

import "github.com/gin-gonic/gin"

// Ambil ID siswa dari context (yang udah di-set di middleware)
func GetSiswaID(c *gin.Context) uint {
	// Ambil value dari context
	val, exists := c.Get("siswaId")
	if !exists {
		return 0
	}
	
	// Convert ke uint
	if id, ok := val.(uint); ok {
		return id
	}
	
	return 0
}
