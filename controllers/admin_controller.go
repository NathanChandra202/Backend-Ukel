package controllers

import (
	"backend/config"
	"backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ===== CRUD MERCH (ADMIN ONLY) =====

// CreateMerchInput - input buat bikin merch baru
type CreateMerchInput struct {
	Nama      string `json:"nama" binding:"required"`
	Deskripsi string `json:"deskripsi"`
	FotoURL   string `json:"foto_url"`
	HargaPoin int    `json:"harga_poin" binding:"required,gt=0"`
	Stok      int    `json:"stok" binding:"required,gte=0"`
}

// AdminCreateMerch - buat merch baru
func AdminCreateMerch(c *gin.Context) {
	var input CreateMerchInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"pesan": "data ga lengkap"})
		return
	}

	// bikin merch baru
	merch := models.Merch{
		Nama:      input.Nama,
		Deskripsi: input.Deskripsi,
		FotoURL:   input.FotoURL,
		HargaPoin: input.HargaPoin,
		Stok:      input.Stok,
		Status:    "tersedia",
	}

	if err := config.DB.Create(&merch).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"pesan": "gagal buat merch"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"pesan": "merch berhasil ditambahkan", "data": merch})
}

// AdminUpdateMerch - update merch
func AdminUpdateMerch(c *gin.Context) {
	id := c.Param("id")
	var input CreateMerchInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"pesan": "data ga lengkap"})
		return
	}

	// update data merch
	if err := config.DB.Model(&models.Merch{}).Where("id = ?", id).Updates(map[string]interface{}{
		"nama":       input.Nama,
		"deskripsi":  input.Deskripsi,
		"foto_url":   input.FotoURL,
		"harga_poin": input.HargaPoin,
		"stok":       input.Stok,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"pesan": "gagal update merch"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"pesan": "merch berhasil diupdate"})
}

// AdminDeleteMerch - hapus merch
func AdminDeleteMerch(c *gin.Context) {
	id := c.Param("id")

	// hapus merch (soft delete)
	if err := config.DB.Delete(&models.Merch{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"pesan": "gagal hapus merch"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"pesan": "merch berhasil dihapus"})
}

// ===== KELOLA PESANAN MERCH (ADMIN) =====

// AdminGetPesananMerch - list semua pesanan merch
func AdminGetPesananMerch(c *gin.Context) {
	var pesanan []models.PesananMerch

	// ambil semua pesanan, join sama data siswa & merch
	if err := config.DB.
		Preload("Siswa").
		Preload("Merch").
		Order("created_at DESC").
		Find(&pesanan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"pesan": "gagal load pesanan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": pesanan})
}

// AdminUpdateStatusPesanan - update status pesanan (pending -> siap_ambil -> sudah_diambil)
func AdminUpdateStatusPesanan(c *gin.Context) {
	id := c.Param("id")
	var input struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"pesan": "status wajib diisi"})
		return
	}

	// validasi status yang bener - sesuai alur koperasi sekolah
	statusValid := map[string]bool{
		"pending":       true,
		"siap_ambil":    true, // admin udah siapkan barang di koperasi
		"sudah_diambil": true, // siswa udah ambil barang
		"dibatalkan":    true,
	}
	if !statusValid[input.Status] {
		c.JSON(http.StatusBadRequest, gin.H{"pesan": "status ga valid. pilih: pending / siap_ambil / sudah_diambil / dibatalkan"})
		return
	}

	if err := config.DB.Model(&models.PesananMerch{}).Where("id = ?", id).
		Update("status", input.Status).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"pesan": "gagal update status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"pesan": "status pesanan diupdate ke: " + input.Status})
}

// AdminKonfirmasiAmbil - konfirmasi siswa ambil barang berdasarkan kode unik
// Admin tinggal scan / input kode yang ditunjukin siswa di koperasi
func AdminKonfirmasiAmbil(c *gin.Context) {
	var input struct {
		KodeAmbil string `json:"kode_ambil" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"pesan": "kode_ambil wajib diisi"})
		return
	}

	// cari pesanan berdasarkan kode
	var pesanan models.PesananMerch
	if err := config.DB.
		Preload("Siswa").
		Preload("Merch").
		Where("kode_ambil = ?", input.KodeAmbil).
		First(&pesanan).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"pesan": "kode tidak ditemukan"})
		return
	}

	// cek statusnya, harus siap_ambil dulu
	if pesanan.Status == "sudah_diambil" {
		c.JSON(http.StatusBadRequest, gin.H{"pesan": "barang ini udah diambil sebelumnya"})
		return
	}
	if pesanan.Status == "dibatalkan" {
		c.JSON(http.StatusBadRequest, gin.H{"pesan": "pesanan ini sudah dibatalkan"})
		return
	}

	// tandai sudah diambil
	if err := config.DB.Model(&pesanan).Update("status", "sudah_diambil").Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"pesan": "gagal konfirmasi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"pesan": "konfirmasi berhasil! barang sudah diambil",
		"data": gin.H{
			"nama_siswa": pesanan.Siswa.Nama,
			"nama_merch": pesanan.Merch.Nama,
			"jumlah":     pesanan.Jumlah,
			"total_poin": pesanan.TotalPoin,
		},
	})
}

// ===== KELOLA AKSI SOSIAL (ADMIN) =====

// AdminUpdateAksiSosial - approve/reject aksi sosial
func AdminUpdateAksiSosial(c *gin.Context) {
	id := c.Param("id")
	var input struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"pesan": "status wajib diisi"})
		return
	}

	// validasi status harus disetujui atau ditolak
	if input.Status != "disetujui" && input.Status != "ditolak" {
		c.JSON(http.StatusBadRequest, gin.H{"pesan": "status harus disetujui atau ditolak"})
		return
	}

	// update status aksi sosial
	if err := config.DB.Model(&models.AksiSosial{}).Where("id = ?", id).
		Update("status", input.Status).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"pesan": "gagal update status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"pesan": "status aksi sosial diupdate"})
}

// AdminGetAksiSosial - list semua aksi sosial (buat admin review)
func AdminGetAksiSosial(c *gin.Context) {
	var aksiList []models.AksiSosial

	// ambil semua aksi sosial, join sama data siswa
	if err := config.DB.
		Preload("Siswa").
		Order("created_at DESC").
		Find(&aksiList).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"pesan": "gagal load aksi sosial"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": aksiList})
}

// ===== DASHBOARD ADMIN =====

// AdminGetStats - statistik dashboard admin
func AdminGetStats(c *gin.Context) {
	var stats struct {
		TotalSiswa        int64 `json:"total_siswa"`
		TotalPoinBeredar  int   `json:"total_poin_beredar"`
		TotalJasa         int64 `json:"total_jasa"`
		TotalPesananMerch int64 `json:"total_pesanan_merch"`
	}

	// hitung total siswa
	config.DB.Model(&models.Siswa{}).Count(&stats.TotalSiswa)

	// hitung total poin beredar (semua saldo siswa dijumlah)
	config.DB.Model(&models.Siswa{}).Select("COALESCE(SUM(saldo_poin), 0)").Scan(&stats.TotalPoinBeredar)

	// hitung total jasa
	config.DB.Model(&models.IklanJasa{}).Count(&stats.TotalJasa)

	// hitung total pesanan merch
	config.DB.Model(&models.PesananMerch{}).Count(&stats.TotalPesananMerch)

	c.JSON(http.StatusOK, gin.H{"data": stats})
}
