package controllers

import (
	"backend/config"
	"backend/models"
	"backend/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

type UpdateProfilInput struct {
	Nama    string `json:"nama" binding:"required"`
	Kelas   string `json:"kelas" binding:"required"`
	Jurusan string `json:"jurusan" binding:"required"`
}

func GetProfil(c *gin.Context) {
	siswaID := utils.GetSiswaID(c)

	var siswa models.Siswa
	if err := config.DB.First(&siswa, siswaID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"pesan": "Profil siswa tidak ditemukan."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"siswa": siswa})
}

func UpdateProfil(c *gin.Context) {
	siswaID := utils.GetSiswaID(c)
	var input UpdateProfilInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"pesan": "Nama, kelas, dan jurusan wajib diisi."})
		return
	}

	if err := config.DB.Model(&models.Siswa{}).Where("id = ?", siswaID).Updates(map[string]interface{}{
		"nama":    input.Nama,
		"kelas":   input.Kelas,
		"jurusan": input.Jurusan,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"pesan": "Gagal memperbarui profil."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"pesan": "Profil berhasil diperbarui."})
}

func DeleteAkun(c *gin.Context) {
	siswaID := utils.GetSiswaID(c)

	if err := config.DB.Delete(&models.Siswa{}, siswaID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"pesan": "Gagal menghapus akun."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"pesan": "Akun berhasil dihapus."})
}

func GetRiwayatPoin(c *gin.Context) {
	siswaID := utils.GetSiswaID(c)

	var transaksi []models.TransaksiJasa
	if err := config.DB.
		Preload("Iklan").
		Preload("Pembeli").
		Preload("Penyedia").
		Where("pembeli_id = ? OR penyedia_id = ?", siswaID, siswaID).
		Order("created_at DESC").
		Find(&transaksi).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"pesan": "Gagal memuat riwayat transaksi."})
		return
	}

	type RiwayatResponse struct {
		ID           uint   `json:"id"`
		JumlahPoin   int    `json:"jumlah_poin"`
		Status       string `json:"status"`
		CreatedAt    string `json:"created_at"`
		NamaJasa     string `json:"nama_jasa"`
		Kategori     string `json:"kategori"`
		NamaPembeli  string `json:"nama_pembeli"`
		NamaPenyedia string `json:"nama_penyedia"`
		ArahPoin     string `json:"arah_poin"`
	}

	var riwayat []RiwayatResponse
	for _, t := range transaksi {
		arahPoin := "masuk"
		if t.PembeliID == siswaID {
			arahPoin = "keluar"
		}

		riwayat = append(riwayat, RiwayatResponse{
			ID:           t.ID,
			JumlahPoin:   t.JumlahPoin,
			Status:       t.Status,
			CreatedAt:    t.CreatedAt.Format("2006-01-02 15:04:05"),
			NamaJasa:     t.Iklan.Judul,
			Kategori:     t.Iklan.Kategori,
			NamaPembeli:  t.Pembeli.Nama,
			NamaPenyedia: t.Penyedia.Nama,
			ArahPoin:     arahPoin,
		})
	}

	c.JSON(http.StatusOK, gin.H{"riwayat": riwayat})
}

func ExportRiwayatExcel(c *gin.Context) {
	siswaID := utils.GetSiswaID(c)

	var transaksi []models.TransaksiJasa
	if err := config.DB.
		Preload("Iklan").
		Preload("Pembeli").
		Preload("Penyedia").
		Where("pembeli_id = ? OR penyedia_id = ?", siswaID, siswaID).
		Order("created_at DESC").
		Find(&transaksi).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"pesan": "Gagal memuat riwayat transaksi."})
		return
	}

	f := excelize.NewFile()
	sheetName := "RiwayatPoin"
	f.SetSheetName("Sheet1", sheetName)

	headers := []string{"Tanggal", "Judul Jasa", "Pembeli", "Penyedia", "Jumlah Poin", "Status", "Jenis Transaksi"}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetName, cell, header)
	}

	rowIdx := 2
	for _, t := range transaksi {
		jenisTx := "Poin Masuk (Pendapatan)"
		if t.PembeliID == siswaID {
			jenisTx = "Poin Keluar (Bayar)"
		}

		f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowIdx), t.CreatedAt.Format("2006-01-02 15:04:05"))
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", rowIdx), t.Iklan.Judul)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", rowIdx), t.Pembeli.Nama)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", rowIdx), t.Penyedia.Nama)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", rowIdx), t.JumlahPoin)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", rowIdx), t.Status)
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", rowIdx), jenisTx)
		rowIdx++
	}

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=riwayat_poin.xlsx")

	if err := f.Write(c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"pesan": "Gagal mengunduh file Excel."})
	}
}
