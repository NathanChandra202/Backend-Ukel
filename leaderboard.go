package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

type Peringkat struct {
	ID        int    `json:"id"`
	Nama      string `json:"nama"`
	Kelas     string `json:"kelas"`
	SaldoPoin int    `json:"saldo_poin"`
}

// ngambil 10 siswa paling rajin nyari poin
func GetLeaderboard(c *gin.Context) {
	rows, err := DB.Query("SELECT id, nama, kelas, saldo_poin FROM siswa ORDER BY saldo_poin DESC LIMIT 10")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"pesan": "Gagal narik data leaderboard"})
		return
	}
	defer rows.Close()

	var data []Peringkat = []Peringkat{}
	for rows.Next() {
		var p Peringkat
		if err := rows.Scan(&p.ID, &p.Nama, &p.Kelas, &p.SaldoPoin); err == nil {
			data = append(data, p)
		}
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}
