package middleware

import (
	"backend/config"
	"backend/models"
	"backend/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AdminMiddleware - middleware cek apakah user adalah admin
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// ambil id siswa dari token (udah lewat AuthMiddleware)
		siswaID := utils.GetSiswaID(c)

		// cek role siswa di database
		var siswa models.Siswa
		if err := config.DB.Select("role").First(&siswa, siswaID).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"pesan": "user ga ketemu"})
			c.Abort()
			return
		}

		// cek role nya admin atau bukan
		if siswa.Role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"pesan": "akses ditolak, khusus admin aja cuy"})
			c.Abort()
			return
		}

		// kalau admin, lanjut gas
		c.Next()
	}
}
