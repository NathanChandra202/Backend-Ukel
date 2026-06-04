package controllers

import (
	"backend/config"
	"backend/models"
	"backend/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BuatJasaInput struct {
	Judul     string `json:"judul" binding:"required"`
	Kategori  string `json:"kategori" binding:"required"`
	Deskripsi string `json:"deskripsi" binding:"required"`
	HargaPoin int    `json:"harga_poin" binding:"required,gt=0"`
}

func GetJasa(c *gin.Context) {
	var jasas []models.IklanJasa

	if err := config.DB.
		Preload("Penyedia").
		Where("status = ?", "tersedia").
		Order("created_at DESC").
		Find(&jasas).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"pesan": "Gagal memuat daftar jasa."})
		return
	}

	type JasaResponse struct {
		ID           uint   `json:"id"`
		Judul        string `json:"judul"`
		Kategori     string `json:"kategori"`
		Deskripsi    string `json:"deskripsi"`
		HargaPoin    int    `json:"harga_poin"`
		PenyediaID   uint   `json:"penyedia_id"`
		NamaPenyedia string `json:"nama_penyedia"`
		Status       string `json:"status"`
		CreatedAt    string `json:"created_at"`
	}

	var response []JasaResponse
	for _, j := range jasas {
		response = append(response, JasaResponse{
			ID:           j.ID,
			Judul:        j.Judul,
			Kategori:     j.Kategori,
			Deskripsi:    j.Deskripsi,
			HargaPoin:    j.HargaPoin,
			PenyediaID:   j.PenyediaID,
			NamaPenyedia: j.Penyedia.Nama,
			Status:       j.Status,
			CreatedAt:    j.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}

func GetDetailJasa(c *gin.Context) {
	id := c.Param("id")

	var jasa models.IklanJasa
	if err := config.DB.Preload("Penyedia").First(&jasa, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"pesan": "Jasa tidak ditemukan."})
		return
	}

	response := gin.H{
		"id":            jasa.ID,
		"judul":         jasa.Judul,
		"kategori":      jasa.Kategori,
		"deskripsi":     jasa.Deskripsi,
		"harga_poin":    jasa.HargaPoin,
		"penyedia_id":   jasa.PenyediaID,
		"nama_penyedia": jasa.Penyedia.Nama,
		"status":        jasa.Status,
		"created_at":    jasa.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}

func BuatJasa(c *gin.Context) {
	siswaID := utils.GetSiswaID(c)
	var input BuatJasaInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"pesan": "Judul, kategori, deskripsi, dan harga poin wajib diisi."})
		return
	}

	jasa := models.IklanJasa{
		PenyediaID: siswaID,
		Judul:      input.Judul,
		Kategori:   input.Kategori,
		Deskripsi:  input.Deskripsi,
		HargaPoin:  input.HargaPoin,
		Status:     "tersedia",
	}

	if err := config.DB.Create(&jasa).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"pesan": "Gagal memposting jasa."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"pesan": "Jasa berhasil diposting."})
}

func AmbilJasa(c *gin.Context) {
	pembeliID := utils.GetSiswaID(c)
	jasaID := c.Param("id")

	var jasa models.IklanJasa
	if err := config.DB.Where("id = ? AND status = ?", jasaID, "tersedia").First(&jasa).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"pesan": "Jasa tidak tersedia atau sudah diambil orang lain."})
		return
	}

	if pembeliID == jasa.PenyediaID {
		c.JSON(http.StatusBadRequest, gin.H{"pesan": "Anda tidak dapat mengambil jasa milik sendiri."})
		return
	}

	var pembeli models.Siswa
	if err := config.DB.First(&pembeli, pembeliID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"pesan": "Terjadi kesalahan pada server. Coba lagi."})
		return
	}

	if pembeli.SaldoPoin < jasa.HargaPoin {
		c.JSON(http.StatusBadRequest, gin.H{"pesan": "Saldo poin tidak cukup untuk mengambil jasa ini."})
		return
	}

	err := config.DB.Transaction(func(tx *gorm.DB) error {
		// Potong poin pembeli
		if err := tx.Model(&models.Siswa{}).Where("id = ?", pembeliID).
			Update("saldo_poin", gorm.Expr("saldo_poin - ?", jasa.HargaPoin)).Error; err != nil {
			return err
		}

		// Update status jasa
		if err := tx.Model(&jasa).Update("status", "diambil").Error; err != nil {
			return err
		}

		// Buat transaksi
		transaksi := models.TransaksiJasa{
			IklanID:    jasa.ID,
			PembeliID:  pembeliID,
			PenyediaID: jasa.PenyediaID,
			JumlahPoin: jasa.HargaPoin,
			Status:     "berjalan",
		}
		if err := tx.Create(&transaksi).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"pesan": "Terjadi kesalahan pada server. Coba lagi."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"pesan": "Berhasil mengambil jasa."})
}

func SelesaikanJasa(c *gin.Context) {
	siswaID := utils.GetSiswaID(c)
	transaksiID := c.Param("id")

	var transaksi models.TransaksiJasa
	if err := config.DB.Where("id = ? AND status = ?", transaksiID, "berjalan").First(&transaksi).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"pesan": "Transaksi tidak valid atau sudah selesai."})
		return
	}

	if transaksi.PenyediaID != siswaID {
		c.JSON(http.StatusForbidden, gin.H{"pesan": "Hanya penyedia jasa yang dapat menandai selesai."})
		return
	}

	err := config.DB.Transaction(func(tx *gorm.DB) error {
		// Tambah poin penyedia
		if err := tx.Model(&models.Siswa{}).Where("id = ?", transaksi.PenyediaID).
			Update("saldo_poin", gorm.Expr("saldo_poin + ?", transaksi.JumlahPoin)).Error; err != nil {
			return err
		}

		// Update status transaksi
		if err := tx.Model(&transaksi).Update("status", "selesai").Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"pesan": "Terjadi kesalahan pada server. Coba lagi."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"pesan": "Jasa selesai. Poin telah ditransfer ke penyedia."})
}
