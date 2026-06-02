package main

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func UploadSosial(c *gin.Context) {
	siswaID := c.GetInt("siswaId")
	deskripsi := strings.TrimSpace(c.PostForm("deskripsi"))
	if deskripsi == "" {
		c.JSON(http.StatusBadRequest, gin.H{"pesan": MsgSosialDeskripsiKosong})
		return
	}

	file, err := c.FormFile("foto")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"pesan": MsgSosialFotoWajib})
		return
	}

	if err := os.MkdirAll("uploads", 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"pesan": MsgSosialSimpanGagal})
		return
	}

	filename := filepath.Base(file.Filename)
	path := filepath.Join("uploads", filename)
	if err := c.SaveUploadedFile(file, path); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"pesan": MsgSosialSimpanGagal})
		return
	}

	poinDiberikan := 5
	tx, err := DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"pesan": MsgServerError})
		return
	}
	defer tx.Rollback()

	if _, err = tx.Exec(
		"INSERT INTO aksi_sosial (siswa_id, deskripsi, foto_url, status, poin_diberikan) VALUES ($1, $2, $3, 'disetujui', $4)",
		siswaID, deskripsi, path, poinDiberikan,
	); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"pesan": MsgServerError})
		return
	}
	if _, err = tx.Exec(
		"UPDATE siswa SET saldo_poin = saldo_poin + $1 WHERE id = $2",
		poinDiberikan, siswaID,
	); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"pesan": MsgServerError})
		return
	}
	if err = tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"pesan": MsgServerError})
		return
	}

	c.JSON(http.StatusOK, gin.H{"pesan": MsgSosialUploadOk})
}

type AksiSosial struct {
	ID        int       `json:"id"`
	Deskripsi string    `json:"deskripsi"`
	FotoURL   string    `json:"foto_url"`
	Status    string    `json:"status"`
	Poin      int       `json:"poin_diberikan"`
	CreatedAt time.Time `json:"created_at"`
}

func GetRiwayatSosial(c *gin.Context) {
	siswaID := c.GetInt("siswaId")
	rows, err := DB.Query(
		"SELECT id, deskripsi, foto_url, status, poin_diberikan, created_at FROM aksi_sosial WHERE siswa_id = $1 ORDER BY created_at DESC",
		siswaID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"pesan": MsgSosialRiwayatGagal})
		return
	}
	defer rows.Close()

	data := []AksiSosial{}
	for rows.Next() {
		var a AksiSosial
		if err := rows.Scan(&a.ID, &a.Deskripsi, &a.FotoURL, &a.Status, &a.Poin, &a.CreatedAt); err == nil {
			data = append(data, a)
		}
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}
