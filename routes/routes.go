package routes

import (
	"backend/controllers"
	"backend/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Setup semua route API
func SetupRoutes(r *gin.Engine) {
	// Health check endpoint (cek server hidup atau mati)
	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"pesan":  "Backend Kontrib.ID aktif",
		})
	})

	// Grup /api
	api := r.Group("/api")

	// ===== AUTH ROUTES (TIDAK PERLU LOGIN) =====
	auth := api.Group("/auth")
	{
		auth.POST("/register", controllers.Register) // Daftar akun baru
		auth.POST("/login", controllers.Login)       // Login dapat token
	}

	// ===== SISWA ROUTES (PERLU LOGIN) =====
	siswa := api.Group("/siswa").Use(middleware.AuthMiddleware())
	{
		siswa.GET("/profil", controllers.GetProfil)           // Lihat profil
		siswa.PUT("", controllers.UpdateProfil)               // Update profil
		siswa.DELETE("", controllers.DeleteAkun)              // Hapus akun
		siswa.GET("/riwayat", controllers.GetRiwayatPoin)     // Riwayat transaksi
		siswa.GET("/export", controllers.ExportRiwayatExcel)  // Download Excel
	}

	// ===== JASA ROUTES (PERLU LOGIN) =====
	jasa := api.Group("/jasa").Use(middleware.AuthMiddleware())
	{
		jasa.GET("", controllers.GetJasa)                // List jasa
		jasa.GET("/:id", controllers.GetDetailJasa)      // Detail jasa
		jasa.POST("", controllers.BuatJasa)              // Post jasa baru
		jasa.POST("/:id/ambil", controllers.AmbilJasa)   // Beli jasa
		jasa.POST("/:id/selesai", controllers.SelesaikanJasa) // Tandai selesai
	}

	// ===== SOSIAL ROUTES (PERLU LOGIN) =====
	sosial := api.Group("/sosial").Use(middleware.AuthMiddleware())
	{
		sosial.POST("/upload", controllers.UploadSosial)        // Upload foto (+5 poin)
		sosial.GET("/riwayat", controllers.GetRiwayatSosial)    // Riwayat upload
	}

	// ===== LEADERBOARD (PERLU LOGIN) =====
	api.GET("/leaderboard", middleware.AuthMiddleware(), controllers.GetLeaderboard)

	// ===== MERCH ROUTES (SISWA - PERLU LOGIN) =====
	merch := api.Group("/merch").Use(middleware.AuthMiddleware())
	{
		merch.GET("", controllers.GetMerchList)                    // list merch
		merch.GET("/:id", controllers.GetMerchDetail)              // detail merch
		merch.POST("/beli", controllers.BeliMerch)                 // beli merch (tukar poin)
		merch.GET("/riwayat", controllers.GetRiwayatPesananMerch)  // riwayat pesanan
	}

	// ===== ADMIN ROUTES (KHUSUS ADMIN) =====
	admin := api.Group("/admin").Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
	{
		// Kelola Merch
		admin.POST("/merch", controllers.AdminCreateMerch)         // buat merch baru
		admin.PUT("/merch/:id", controllers.AdminUpdateMerch)      // update merch
		admin.DELETE("/merch/:id", controllers.AdminDeleteMerch)   // hapus merch

		// Kelola Pesanan Merch
		admin.GET("/pesanan", controllers.AdminGetPesananMerch)    // list semua pesanan
		admin.PUT("/pesanan/:id", controllers.AdminUpdateStatusPesanan) // update status pesanan
		admin.POST("/pesanan/konfirmasi", controllers.AdminKonfirmasiAmbil) // konfirmasi ambil via kode

		// Kelola Aksi Sosial
		admin.GET("/sosial", controllers.AdminGetAksiSosial)       // list semua aksi sosial
		admin.PUT("/sosial/:id", controllers.AdminUpdateAksiSosial) // approve/reject aksi sosial

		// Dashboard Stats
		admin.GET("/stats", controllers.AdminGetStats)             // statistik admin
	}
}
