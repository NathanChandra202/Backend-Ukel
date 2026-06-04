package controllers

import (
	"backend/config"
	"backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetLeaderboard(c *gin.Context) {
	var siswa []models.Siswa

	if err := config.DB.
		Select("id, nama, kelas, saldo_poin").
		Order("saldo_poin DESC").
		Limit(10).
		Find(&siswa).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"pesan": "Gagal memuat leaderboard."})
		return
	}

	type Peringkat struct {
		ID        uint   `json:"id"`
		Nama      string `json:"nama"`
		Kelas     string `json:"kelas"`
		SaldoPoin int    `json:"saldo_poin"`
	}

	var data []Peringkat
	for _, s := range siswa {
		data = append(data, Peringkat{
			ID:        s.ID,
			Nama:      s.Nama,
			Kelas:     s.Kelas,
			SaldoPoin: s.SaldoPoin,
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": data})
}
