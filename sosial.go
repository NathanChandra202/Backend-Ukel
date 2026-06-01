package main

import (
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

// fungsi buat upload foto aksi sosial 
func UploadSosial(c *gin.Context) {
	siswaID := c.GetInt("siswaId")
	deskripsi := c.PostForm("deskripsi")

	file, err := c.FormFile("foto")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"pesan": "Foto harus dilampirin bro"})
		return
	}

	filename := filepath.Base(file.Filename)
	path := "uploads/" + filename
	if err := c.SaveUploadedFile(file, path); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"pesan": "Gagal nyimpen foto"})
		return
	}

	poinDiberikan := 5
	tx, _ := DB.Begin()
	tx.Exec("INSERT INTO aksi_sosial (siswa_id, deskripsi, foto_url, status, poin_diberikan) VALUES ($1, $2, $3, 'disetujui', $4)", siswaID, deskripsi, path, poinDiberikan)
	tx.Exec("UPDATE siswa SET saldo_poin = saldo_poin + $1 WHERE id = $2", poinDiberikan, siswaID)
	tx.Commit()

	c.JSON(http.StatusOK, gin.H{"pesan": "Aksi sosial diupload, poin lu nambah!"})
}

type AksiSosial struct {
	ID        int       `json:"id"`
	Deskripsi string    `json:"deskripsi"`
	FotoURL   string    `json:"foto_url"`
	Status    string    `json:"status"`
	Poin      int       `json:"poin_diberikan"`
	CreatedAt time.Time `json:"created_at"`
}

// ngambil riwayat aksi sosial yang udah diupload
func GetRiwayatSosial(c *gin.Context) {
	siswaID := c.GetInt("siswaId")
	rows, err := DB.Query("SELECT id, deskripsi, foto_url, status, poin_diberikan, created_at FROM aksi_sosial WHERE siswa_id = $1 ORDER BY created_at DESC", siswaID)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"pesan": "Gagal narik riwayat"})
		return
	}
	defer rows.Close()

	var data []AksiSosial = []AksiSosial{}
	for rows.Next() {
		var a AksiSosial
		if err := rows.Scan(&a.ID, &a.Deskripsi, &a.FotoURL, &a.Status, &a.Poin, &a.CreatedAt); err == nil {
			data = append(data, a)
		}
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}
