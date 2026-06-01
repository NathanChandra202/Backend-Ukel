package main

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type RegisterInput struct {
	Nama     string `json:"nama" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	Kelas    string `json:"kelas" binding:"required"`
	Jurusan  string `json:"jurusan" binding:"required"`
}

// fungsi buat daftar user baru
func Register(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"pesan": err.Error()})
		return
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)

	_, err := DB.Exec("INSERT INTO siswa (nama, email, password, kelas, jurusan) VALUES ($1, $2, $3, $4, $5)",
		input.Nama, input.Email, string(hash), input.Kelas, input.Jurusan)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"pesan": "Email udah kepake atau error lain"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"pesan": "Berhasil daftar, silakan login!"})
}

type LoginInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// fungsi buat login dan generate token
func Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"pesan": err.Error()})
		return
	}

	var id int
	var nama, hash string
	err := DB.QueryRow("SELECT id, nama, password FROM siswa WHERE email = $1", input.Email).Scan(&id, &nama, &hash)
	
	if err != nil || bcrypt.CompareHashAndPassword([]byte(hash), []byte(input.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"pesan": "Email atau password salah nih"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":   id,
		"exp":  time.Now().Add(time.Hour * 72).Unix(),
	})

	secretKey := []byte(os.Getenv("JWT_SECRET"))
	if len(secretKey) == 0 {
		secretKey = []byte("rahasiabanget")
	}

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"pesan": "Gagal bikin token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
		"siswa": gin.H{
			"id":   id,
			"nama": nama,
		},
	})
}
