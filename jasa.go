package main

import (
	"net/http"
	"time"
	"github.com/gin-gonic/gin"
)

type IklanJasa struct {
	ID          int       `json:"id"`
	Judul       string    `json:"judul"`
	Kategori    string    `json:"kategori"`
	Deskripsi   string    `json:"deskripsi"`
	HargaPoin   int       `json:"harga_poin"`
	PenyediaID  int       `json:"penyedia_id"`
	NamaPenyedia string   `json:"nama_penyedia"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

// ngambil semua jasa yang statusnya tersedia
func GetJasa(c *gin.Context) {
	rows, err := DB.Query("SELECT i.id, i.judul, i.kategori, i.deskripsi, i.harga_poin, i.penyedia_id, s.nama, i.status, i.created_at FROM iklan_jasa i JOIN siswa s ON s.id = i.penyedia_id WHERE i.status = 'tersedia' ORDER BY i.created_at DESC")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"pesan": "Gagal narik data jasa"})
		return
	}
	defer rows.Close()

	var jasas []IklanJasa = []IklanJasa{}
	for rows.Next() {
		var j IklanJasa
		if err := rows.Scan(&j.ID, &j.Judul, &j.Kategori, &j.Deskripsi, &j.HargaPoin, &j.PenyediaID, &j.NamaPenyedia, &j.Status, &j.CreatedAt); err == nil {
			jasas = append(jasas, j)
		}
	}
	c.JSON(http.StatusOK, gin.H{"data": jasas})
}

// ngambil detail satu jasa spesifik
func DetailJasa(c *gin.Context) {
	id := c.Param("id")
	var j IklanJasa
	err := DB.QueryRow("SELECT i.id, i.judul, i.kategori, i.deskripsi, i.harga_poin, i.penyedia_id, s.nama, i.status, i.created_at FROM iklan_jasa i JOIN siswa s ON s.id = i.penyedia_id WHERE i.id = $1", id).Scan(&j.ID, &j.Judul, &j.Kategori, &j.Deskripsi, &j.HargaPoin, &j.PenyediaID, &j.NamaPenyedia, &j.Status, &j.CreatedAt)
	
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"pesan": "Jasa nggak ketemu bro"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": j})
}

type BuatJasaInput struct {
	Judul     string `json:"judul" binding:"required"`
	Kategori  string `json:"kategori" binding:"required"`
	Deskripsi string `json:"deskripsi" binding:"required"`
	HargaPoin int    `json:"harga_poin" binding:"required"`
}

// bikin penawaran jasa baru
func BuatJasa(c *gin.Context) {
	siswaID := c.GetInt("siswaId")
	var input BuatJasaInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"pesan": err.Error()})
		return
	}

	_, err := DB.Exec("INSERT INTO iklan_jasa (penyedia_id, judul, kategori, deskripsi, harga_poin) VALUES ($1, $2, $3, $4, $5)",
		siswaID, input.Judul, input.Kategori, input.Deskripsi, input.HargaPoin)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"pesan": "Gagal ngepost jasa"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"pesan": "Jasa berhasil diposting!"})
}

// ngambil tugas dari orang lain buat dikerjain
func AmbilJasa(c *gin.Context) {
	pembeliID := c.GetInt("siswaId")
	jasaID := c.Param("id")

	var penyediaID, hargaPoin int
	err := DB.QueryRow("SELECT penyedia_id, harga_poin FROM iklan_jasa WHERE id = $1 AND status = 'tersedia'", jasaID).Scan(&penyediaID, &hargaPoin)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"pesan": "Jasa nggak tersedia atau udah diambil orang"})
		return
	}

	if pembeliID == penyediaID {
		c.JSON(http.StatusBadRequest, gin.H{"pesan": "Masa ngambil jasa sendiri bro"})
		return
	}

	var saldoPembeli int
	DB.QueryRow("SELECT saldo_poin FROM siswa WHERE id = $1", pembeliID).Scan(&saldoPembeli)
	if saldoPembeli < hargaPoin {
		c.JSON(http.StatusBadRequest, gin.H{"pesan": "Poin lu kurang bro buat ngambil jasa ini"})
		return
	}

	tx, _ := DB.Begin()
	tx.Exec("UPDATE siswa SET saldo_poin = saldo_poin - $1 WHERE id = $2", hargaPoin, pembeliID)
	tx.Exec("UPDATE iklan_jasa SET status = 'diambil' WHERE id = $1", jasaID)
	tx.Exec("INSERT INTO transaksi_jasa (iklan_id, pembeli_id, penyedia_id, jumlah_poin) VALUES ($1, $2, $3, $4)", jasaID, pembeliID, penyediaID, hargaPoin)
	tx.Commit()

	c.JSON(http.StatusOK, gin.H{"pesan": "Berhasil ngambil jasa, kerjain yang bener ya!"})
}

// nyelesain transaksi kalo jasanya udah kelar
func SelesaikanJasa(c *gin.Context) {
	transaksiID := c.Param("id")
	
	var penyediaID, jumlahPoin int
	err := DB.QueryRow("SELECT penyedia_id, jumlah_poin FROM transaksi_jasa WHERE id = $1 AND status = 'berjalan'", transaksiID).Scan(&penyediaID, &jumlahPoin)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"pesan": "Transaksi ga valid atau udah kelar"})
		return
	}

	tx, _ := DB.Begin()
	tx.Exec("UPDATE siswa SET saldo_poin = saldo_poin + $1 WHERE id = $2", jumlahPoin, penyediaID)
	tx.Exec("UPDATE transaksi_jasa SET status = 'selesai' WHERE id = $1", transaksiID)
	tx.Commit()

	c.JSON(http.StatusOK, gin.H{"pesan": "Jasa selesai, poin udah ditransfer ke penyedia"})
}
