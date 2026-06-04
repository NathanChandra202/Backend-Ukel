package controllers

import (
	"backend/config"
	"backend/models"
	"backend/utils"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UploadSosial(c *gin.Context) {
	siswaID := utils.GetSiswaID(c)
	deskripsi := strings.TrimSpace(c.PostForm("deskripsi"))

	if deskripsi == "" {
		c.JSON(http.StatusBadRequest, gin.H{"pesan": "Deskripsi kegiatan wajib diisi."})
		return
	}

	file, err := c.FormFile("foto")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"pesan": "Foto bukti wajib dilampirkan (JPG/PNG)."})
		return
	}

	if err := os.MkdirAll("uploads", 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"pesan": "Gagal menyimpan foto."})
		return
	}

	filename := filepath.Base(file.Filename)
	path := filepath.Join("uploads", filename)
	if err := c.SaveUploadedFile(file, path); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"pesan": "Gagal menyimpan foto."})
		return
	}

	poinDiberikan := 5

	err = config.DB.Transaction(func(tx *gorm.DB) error {
		// Simpan aksi sosial
		aksi := models.AksiSosial{
			SiswaID:       siswaID,
			Deskripsi:     deskripsi,
			FotoURL:       path,
			Status:        "disetujui",
			PoinDiberikan: poinDiberikan,
		}
		if err := tx.Create(&aksi).Error; err != nil {
			return err
		}

		// Tambah poin siswa
		if err := tx.Model(&models.Siswa{}).Where("id = ?", siswaID).
			Update("saldo_poin", gorm.Expr("saldo_poin + ?", poinDiberikan)).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"pesan": "Terjadi kesalahan pada server. Coba lagi."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"pesan": "Aksi sosial berhasil diunggah. Poin bertambah."})
}

func GetRiwayatSosial(c *gin.Context) {
	siswaID := utils.GetSiswaID(c)

	var aksiSosial []models.AksiSosial
	if err := config.DB.Where("siswa_id = ?", siswaID).Order("created_at DESC").Find(&aksiSosial).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"pesan": "Gagal memuat riwayat aksi sosial."})
		return
	}

	type AksiResponse struct {
		ID        uint   `json:"id"`
		Deskripsi string `json:"deskripsi"`
		FotoURL   string `json:"foto_url"`
		Status    string `json:"status"`
		Poin      int    `json:"poin_diberikan"`
		CreatedAt string `json:"created_at"`
	}

	var response []AksiResponse
	for _, a := range aksiSosial {
		response = append(response, AksiResponse{
			ID:        a.ID,
			Deskripsi: a.Deskripsi,
			FotoURL:   a.FotoURL,
			Status:    a.Status,
			Poin:      a.PoinDiberikan,
			CreatedAt: a.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}
