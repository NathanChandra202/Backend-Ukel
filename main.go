package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

var DB *sql.DB

// fungsi utama buat jalanin servernya
func main() {
	if err := koneksiDatabase(); err != nil {
		log.Fatal(
			"Gagal konek database:\n  ", err,
			"\n\nPerbaiki backend/.env lalu pastikan PostgreSQL jalan dan database kontribid + schema.sql sudah dibuat.",
		)
	}

	r := gin.Default()

	// CORS untuk Flutter Web (Chrome) dan mobile
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Accept")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	api := r.Group("/api")

	api.GET("/health", func(c *gin.Context) {
		if err := DB.Ping(); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "error",
				"pesan":  "Database tidak terhubung. Periksa PostgreSQL dan file backend/.env.",
				"detail": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"pesan":  "Backend Kontrib.ID aktif",
		})
	})

	// rute autentikasi
	auth := api.Group("/auth")
	auth.POST("/register", Register)
	auth.POST("/login", Login)

	// rute siswa yang butuh login
	siswa := api.Group("/siswa").Use(CekToken)
	siswa.GET("/profil", ProfilSiswa)
	siswa.GET("/riwayat", RiwayatPoin)
	siswa.PUT("", UpdateProfil)
	siswa.DELETE("", DeleteAkun)
	siswa.GET("/export", ExportRiwayatExcel)

	// rute jasa
	jasa := api.Group("/jasa").Use(CekToken)
	jasa.GET("", GetJasa)
	jasa.GET("/:id", DetailJasa)
	jasa.POST("", BuatJasa)
	jasa.POST("/:id/ambil", AmbilJasa)
	jasa.POST("/:id/selesai", SelesaikanJasa)

	// rute sosial
	sosial := api.Group("/sosial").Use(CekToken)
	sosial.POST("/upload", UploadSosial)
	sosial.GET("/riwayat", GetRiwayatSosial)

	// rute leaderboard
	api.GET("/leaderboard", CekToken, GetLeaderboard)

	log.Println("Server aktif → http://localhost:3000")
	log.Println("Health check → http://localhost:3000/api/health")
	r.Run(":3000")
}
