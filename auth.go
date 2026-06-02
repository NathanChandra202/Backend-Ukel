package main

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type RegisterInput struct {
	Nama     string `json:"nama" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Kelas    string `json:"kelas" binding:"required"`
	Jurusan  string `json:"jurusan" binding:"required"`
}

func Register(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"pesan": MsgRegisterFieldKosong})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"pesan": MsgServerError})
		return
	}

	_, err = DB.Exec(
		"INSERT INTO siswa (nama, email, password, kelas, jurusan) VALUES ($1, $2, $3, $4, $5)",
		input.Nama, strings.TrimSpace(strings.ToLower(input.Email)), string(hash), input.Kelas, input.Jurusan,
	)

	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"pesan": MsgRegisterEmailAda})
		return
	}

	c.JSON(http.StatusOK, gin.H{"pesan": MsgRegisterBerhasil})
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"pesan": MsgLoginFieldKosong})
		return
	}

	var id, saldoPoin int
	var nama, hash string
	email := strings.TrimSpace(strings.ToLower(input.Email))
	err := DB.QueryRow(
		"SELECT id, nama, password, saldo_poin FROM siswa WHERE email = $1",
		email,
	).Scan(&id, &nama, &hash, &saldoPoin)

	if err != nil || bcrypt.CompareHashAndPassword([]byte(hash), []byte(input.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"pesan": MsgLoginGagal})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	})

	secretKey := []byte(os.Getenv("JWT_SECRET"))
	if len(secretKey) == 0 {
		secretKey = []byte("rahasiabanget")
	}

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"pesan": MsgTokenGagal})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
		"siswa": gin.H{
			"id":         id,
			"nama":       nama,
			"saldo_poin": saldoPoin,
		},
	})
}
