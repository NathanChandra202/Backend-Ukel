package controllers

import (
	"backend/config"
	"backend/models"
	"backend/utils"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// bikin kode unik 6 digit buat ambil barang di koperasi
// contoh hasilnya: "KOP-4F2B9A"
func buatKodeAmbil() string {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	const huruf = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	kode := make([]byte, 6)
	for i := range kode {
		kode[i] = huruf[rand.Intn(len(huruf))]
	}
	return fmt.Sprintf("KOP-%s", string(kode))
}

// GetMerchList - list semua merch yang tersedia
func GetMerchList(c *gin.Context) {
	var merchList []models.Merch

	if err := config.DB.
		Where("status = ? AND stok > 0", "tersedia").
		Order("created_at DESC").
		Find(&merchList).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"pesan": "gagal load katalog merch"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": merchList})
}

// GetMerchDetail - detail satu merch
func GetMerchDetail(c *gin.Context) {
	id := c.Param("id")

	var merch models.Merch
	if err := config.DB.First(&merch, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"pesan": "merch ga ketemu"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": merch})
}

// BeliMerchInput - input buat beli merch, ga perlu alamat lagi, ambil di koperasi aja
type BeliMerchInput struct {
	MerchID uint   `json:"merch_id"`
	Jumlah  int    `json:"jumlah"`
	Catatan string `json:"catatan"`
}

// BeliMerch - tukar poin ke merch, nanti ambil di koperasi sekolah
func BeliMerch(c *gin.Context) {
	siswaID := utils.GetSiswaID(c)
	var input BeliMerchInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"pesan": "format data salah"})
		return
	}

	// validasi manual biar pesan error lebih jelas
	if input.MerchID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"pesan": "merch_id wajib diisi"})
		return
	}
	if input.Jumlah <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"pesan": "jumlah harus lebih dari 0"})
		return
	}

	// cek merch ada ga
	var merch models.Merch
	if err := config.DB.First(&merch, input.MerchID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"pesan": "merch ga ketemu"})
		return
	}

	// cek stok cukup ga
	if merch.Stok < input.Jumlah {
		c.JSON(http.StatusBadRequest, gin.H{"pesan": "stok ga cukup bro"})
		return
	}

	// hitung total poin
	totalPoin := merch.HargaPoin * input.Jumlah

	// cek saldo poin siswa
	var siswa models.Siswa
	if err := config.DB.First(&siswa, siswaID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"pesan": "error cari data siswa"})
		return
	}

	if siswa.SaldoPoin < totalPoin {
		c.JSON(http.StatusBadRequest, gin.H{"pesan": "poin lo ga cukup bro"})
		return
	}

	// generate kode ambil dulu sebelum transaksi
	kodeAmbil := buatKodeAmbil()

	var pesananBaru models.PesananMerch

	// transaksi: potong poin, kurangi stok, buat pesanan
	err := config.DB.Transaction(func(tx *gorm.DB) error {
		// potong poin siswa
		if err := tx.Model(&models.Siswa{}).Where("id = ?", siswaID).
			Update("saldo_poin", gorm.Expr("saldo_poin - ?", totalPoin)).Error; err != nil {
			return err
		}

		// kurangi stok merch
		if err := tx.Model(&merch).Update("stok", gorm.Expr("stok - ?", input.Jumlah)).Error; err != nil {
			return err
		}

		// buat pesanan dengan kode ambil unik
		pesananBaru = models.PesananMerch{
			SiswaID:   siswaID,
			MerchID:   input.MerchID,
			Jumlah:    input.Jumlah,
			TotalPoin: totalPoin,
			Status:    "pending",
			KodeAmbil: kodeAmbil,
			Catatan:   input.Catatan,
		}
		if err := tx.Create(&pesananBaru).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"pesan": "transaksi gagal"})
		return
	}

	// return kode ambil biar siswa bisa langsung liat
	c.JSON(http.StatusOK, gin.H{
		"pesan":      "berhasil tukar poin! tunjukkan kode ini ke koperasi sekolah",
		"kode_ambil": kodeAmbil,
		"pesanan_id": pesananBaru.ID,
	})
}

// GetRiwayatPesananMerch - riwayat pesanan merch siswa
func GetRiwayatPesananMerch(c *gin.Context) {
	siswaID := utils.GetSiswaID(c)

	var pesanan []models.PesananMerch
	if err := config.DB.
		Preload("Merch").
		Where("siswa_id = ?", siswaID).
		Order("created_at DESC").
		Find(&pesanan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"pesan": "gagal load riwayat"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": pesanan})
}
