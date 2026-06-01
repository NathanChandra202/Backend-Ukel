package main

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var DB *sql.DB

// fungsi utama buat jalanin servernya
func main() {
	var err error
	
	// koneksi ke database postgres (sesuaiin stringnya ya sama setup lu)
	connStr := "user=postgres password=postgres dbname=kontribid sslmode=disable"
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Gagal konek database: ", err)
	}

	r := gin.Default()

	// setting CORS biar flutter web atau hp bisa nembak api
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	api := r.Group("/api")

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

	log.Println("Server jalan di port 3000 nih bro...")
	r.Run(":3000")
}
