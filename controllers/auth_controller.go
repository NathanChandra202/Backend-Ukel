package controllers

import (
	"backend/config"
	"backend/models"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Struct untuk input register
type RegisterInput struct {
	Nama     string `json:"nama" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Kelas    string `json:"kelas" binding:"required"`
	Jurusan  string `json:"jurusan" binding:"required"`
}

// Struct untuk input login
type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Daftar akun baru
func Register(c *gin.Context) {
	var input RegisterInput
	
	// Validasi input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"pesan": "Data tidak lengkap. Isi semua field."})
		return
	}

	// Hash password pakai bcrypt (biar aman)
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"pesan": "Error saat hash password."})
		return
	}

	// Bikin object siswa baru
	siswa := models.Siswa{
		Nama:      input.Nama,
		Email:     strings.ToLower(strings.TrimSpace(input.Email)),
		Password:  string(hashPassword),
		Kelas:     input.Kelas,
		Jurusan:   input.Jurusan,
		SaldoPoin: 7, // Saldo awal 7 poin
	}

	// Simpan ke database
	if err := config.DB.Create(&siswa).Error; err != nil {
		c.JSON(http.StatusConflict, gin.H{"pesan": "Email sudah dipakai. Pakai email lain."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"pesan": "Daftar berhasil! Silakan login."})
}

// Login dan dapat token
func Login(c *gin.Context) {
	var input LoginInput
	
	// Validasi input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"pesan": "Email dan password wajib diisi."})
		return
	}

	// Cari user di database by email
	var siswa models.Siswa
	email := strings.ToLower(strings.TrimSpace(input.Email))
	
	if err := config.DB.Where("email = ?", email).First(&siswa).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"pesan": "Email atau password salah."})
		return
	}

	// Cek password cocok atau tidak
	if err := bcrypt.CompareHashAndPassword([]byte(siswa.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"pesan": "Email atau password salah."})
		return
	}

	// Buat JWT token (berlaku 72 jam)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  siswa.ID,
		"exp": time.Now().Add(72 * time.Hour).Unix(),
	})

	// Ambil secret key dari .env
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		secretKey = "rahasiabanget"
	}

	// Sign token jadi string
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"pesan": "Gagal buat token."})
		return
	}

	// Kirim response token + data siswa
	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
		"siswa": gin.H{
			"id":         siswa.ID,
			"nama":       siswa.Nama,
			"saldo_poin": siswa.SaldoPoin,
		},
	})
}
