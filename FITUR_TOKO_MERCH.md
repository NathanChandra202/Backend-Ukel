# 🛍️ Fitur Toko Merchandise + Admin

## 📋 Overview

Tambahan fitur baru:
1. **Toko Merch** - Siswa bisa tukar poin ke barang fisik
2. **Role Admin** - Admin kelola sistem

---

## 🗄️ Database Schema Baru

### 1. Tambah kolom `role` di tabel `siswa`
```sql
ALTER TABLE siswa ADD COLUMN role VARCHAR(20) DEFAULT 'siswa';
-- Possible values: 'siswa', 'admin'
```

### 2. Tabel `merch` (katalog barang)
```sql
CREATE TABLE merch (
  id            SERIAL PRIMARY KEY,
  nama          VARCHAR(100) NOT NULL,
  deskripsi     TEXT,
  foto_url      VARCHAR(255),
  harga_poin    INTEGER NOT NULL,
  stok          INTEGER DEFAULT 0,
  status        VARCHAR(20) DEFAULT 'tersedia', -- tersedia/habis
  created_at    TIMESTAMP DEFAULT NOW(),
  updated_at    TIMESTAMP DEFAULT NOW()
);
```

### 3. Tabel `pesanan_merch` (transaksi beli barang)
```sql
CREATE TABLE pesanan_merch (
  id            SERIAL PRIMARY KEY,
  siswa_id      INTEGER REFERENCES siswa(id) ON DELETE CASCADE,
  merch_id      INTEGER REFERENCES merch(id) ON DELETE CASCADE,
  jumlah        INTEGER NOT NULL DEFAULT 1,
  total_poin    INTEGER NOT NULL,
  status        VARCHAR(20) DEFAULT 'pending', -- pending/diproses/selesai/dibatalkan
  alamat        TEXT,
  catatan       TEXT,
  created_at    TIMESTAMP DEFAULT NOW(),
  updated_at    TIMESTAMP DEFAULT NOW()
);
```

---

## 🔧 Backend Implementation

### Models

#### `models/merch.go`
```go
package models

import (
	"time"
	"gorm.io/gorm"
)

type Merch struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Nama        string         `gorm:"size:100;not null" json:"nama"`
	Deskripsi   string         `gorm:"type:text" json:"deskripsi"`
	FotoURL     string         `gorm:"size:255" json:"foto_url"`
	HargaPoin   int            `gorm:"not null" json:"harga_poin"`
	Stok        int            `gorm:"default:0" json:"stok"`
	Status      string         `gorm:"size:20;default:tersedia" json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Merch) TableName() string {
	return "merch"
}
```

#### `models/pesanan_merch.go`
```go
package models

import (
	"time"
	"gorm.io/gorm"
)

type PesananMerch struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	SiswaID    uint           `gorm:"not null" json:"siswa_id"`
	MerchID    uint           `gorm:"not null" json:"merch_id"`
	Jumlah     int            `gorm:"not null;default:1" json:"jumlah"`
	TotalPoin  int            `gorm:"not null" json:"total_poin"`
	Status     string         `gorm:"size:20;default:pending" json:"status"`
	Alamat     string         `gorm:"type:text" json:"alamat"`
	Catatan    string         `gorm:"type:text" json:"catatan"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
	
	Siswa      Siswa  `gorm:"foreignKey:SiswaID" json:"siswa,omitempty"`
	Merch      Merch  `gorm:"foreignKey:MerchID" json:"merch,omitempty"`
}

func (PesananMerch) TableName() string {
	return "pesanan_merch"
}
```

---

### Controllers

#### `controllers/merch_controller.go` (Siswa)
```go
package controllers

import (
	"backend/config"
	"backend/models"
	"backend/utils"
	"net/http"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// list semua merch yang tersedia
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

// detail satu merch
func GetMerchDetail(c *gin.Context) {
	id := c.Param("id")
	
	var merch models.Merch
	if err := config.DB.First(&merch, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"pesan": "merch ga ketemu"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"data": merch})
}

type BeliMerchInput struct {
	MerchID uint   `json:"merch_id" binding:"required"`
	Jumlah  int    `json:"jumlah" binding:"required,gt=0"`
	Alamat  string `json:"alamat" binding:"required"`
	Catatan string `json:"catatan"`
}

// beli merch (tukar poin)
func BeliMerch(c *gin.Context) {
	siswaID := utils.GetSiswaID(c)
	var input BeliMerchInput
	
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"pesan": "data ga lengkap"})
		return
	}
	
	// cek merch ada & stok cukup
	var merch models.Merch
	if err := config.DB.First(&merch, input.MerchID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"pesan": "merch ga ketemu"})
		return
	}
	
	if merch.Stok < input.Jumlah {
		c.JSON(http.StatusBadRequest, gin.H{"pesan": "stok ga cukup"})
		return
	}
	
	// hitung total poin
	totalPoin := merch.HargaPoin * input.Jumlah
	
	// cek saldo poin siswa
	var siswa models.Siswa
	if err := config.DB.First(&siswa, siswaID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"pesan": "error"})
		return
	}
	
	if siswa.SaldoPoin < totalPoin {
		c.JSON(http.StatusBadRequest, gin.H{"pesan": "poin lo ga cukup bro"})
		return
	}
	
	// transaksi: potong poin, kurangi stok, buat pesanan
	err := config.DB.Transaction(func(tx *gorm.DB) error {
		// potong poin
		if err := tx.Model(&models.Siswa{}).Where("id = ?", siswaID).
			Update("saldo_poin", gorm.Expr("saldo_poin - ?", totalPoin)).Error; err != nil {
			return err
		}
		
		// kurangi stok
		if err := tx.Model(&merch).Update("stok", gorm.Expr("stok - ?", input.Jumlah)).Error; err != nil {
			return err
		}
		
		// buat pesanan
		pesanan := models.PesananMerch{
			SiswaID:   siswaID,
			MerchID:   input.MerchID,
			Jumlah:    input.Jumlah,
			TotalPoin: totalPoin,
			Status:    "pending",
			Alamat:    input.Alamat,
			Catatan:   input.Catatan,
		}
		if err := tx.Create(&pesanan).Error; err != nil {
			return err
		}
		
		return nil
	})
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"pesan": "transaksi gagal"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"pesan": "berhasil beli merch! tunggu diproses ya"})
}

// riwayat pesanan merch siswa
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
```

#### `controllers/admin_controller.go` (Admin)
```go
package controllers

import (
	"backend/config"
	"backend/models"
	"backend/utils"
	"net/http"
	"github.com/gin-gonic/gin"
)

// CRUD Merch (Admin only)

type CreateMerchInput struct {
	Nama       string `json:"nama" binding:"required"`
	Deskripsi  string `json:"deskripsi"`
	FotoURL    string `json:"foto_url"`
	HargaPoin  int    `json:"harga_poin" binding:"required,gt=0"`
	Stok       int    `json:"stok" binding:"required,gte=0"`
}

// buat merch baru
func AdminCreateMerch(c *gin.Context) {
	var input CreateMerchInput
	
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"pesan": "data ga lengkap"})
		return
	}
	
	merch := models.Merch{
		Nama:       input.Nama,
		Deskripsi:  input.Deskripsi,
		FotoURL:    input.FotoURL,
		HargaPoin:  input.HargaPoin,
		Stok:       input.Stok,
		Status:     "tersedia",
	}
	
	if err := config.DB.Create(&merch).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"pesan": "gagal buat merch"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"pesan": "merch berhasil ditambahkan", "data": merch})
}

// update merch
func AdminUpdateMerch(c *gin.Context) {
	id := c.Param("id")
	var input CreateMerchInput
	
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"pesan": "data ga lengkap"})
		return
	}
	
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

// hapus merch
func AdminDeleteMerch(c *gin.Context) {
	id := c.Param("id")
	
	if err := config.DB.Delete(&models.Merch{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"pesan": "gagal hapus merch"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"pesan": "merch berhasil dihapus"})
}

// list semua pesanan merch (admin)
func AdminGetPesananMerch(c *gin.Context) {
	var pesanan []models.PesananMerch
	
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

// update status pesanan (admin proses pesanan)
func AdminUpdateStatusPesanan(c *gin.Context) {
	id := c.Param("id")
	var input struct {
		Status string `json:"status" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"pesan": "status wajib diisi"})
		return
	}
	
	// validasi status
	if input.Status != "pending" && input.Status != "diproses" && 
	   input.Status != "selesai" && input.Status != "dibatalkan" {
		c.JSON(http.StatusBadRequest, gin.H{"pesan": "status ga valid"})
		return
	}
	
	if err := config.DB.Model(&models.PesananMerch{}).Where("id = ?", id).
		Update("status", input.Status).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"pesan": "gagal update status"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"pesan": "status pesanan diupdate"})
}

// approve/reject aksi sosial (admin)
func AdminUpdateAksiSosial(c *gin.Context) {
	id := c.Param("id")
	var input struct {
		Status string `json:"status" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"pesan": "status wajib diisi"})
		return
	}
	
	// validasi status
	if input.Status != "disetujui" && input.Status != "ditolak" {
		c.JSON(http.StatusBadRequest, gin.H{"pesan": "status harus disetujui atau ditolak"})
		return
	}
	
	if err := config.DB.Model(&models.AksiSosial{}).Where("id = ?", id).
		Update("status", input.Status).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"pesan": "gagal update status"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"pesan": "status aksi sosial diupdate"})
}

// list semua aksi sosial (admin)
func AdminGetAksiSosial(c *gin.Context) {
	var aksiList []models.AksiSosial
	
	if err := config.DB.
		Preload("Siswa").
		Order("created_at DESC").
		Find(&aksiList).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"pesan": "gagal load aksi sosial"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"data": aksiList})
}

// statistik dashboard admin
func AdminGetStats(c *gin.Context) {
	var stats struct {
		TotalSiswa      int64 `json:"total_siswa"`
		TotalPoinBeredar int   `json:"total_poin_beredar"`
		TotalJasa       int64 `json:"total_jasa"`
		TotalPesananMerch int64 `json:"total_pesanan_merch"`
	}
	
	config.DB.Model(&models.Siswa{}).Count(&stats.TotalSiswa)
	config.DB.Model(&models.Siswa{}).Select("COALESCE(SUM(saldo_poin), 0)").Scan(&stats.TotalPoinBeredar)
	config.DB.Model(&models.IklanJasa{}).Count(&stats.TotalJasa)
	config.DB.Model(&models.PesananMerch{}).Count(&stats.TotalPesananMerch)
	
	c.JSON(http.StatusOK, gin.H{"data": stats})
}
```

---

### Middleware Admin

#### `middleware/admin.go`
```go
package middleware

import (
	"backend/config"
	"backend/models"
	"backend/utils"
	"net/http"
	"github.com/gin-gonic/gin"
)

// middleware cek apakah user adalah admin
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		siswaID := utils.GetSiswaID(c)
		
		var siswa models.Siswa
		if err := config.DB.Select("role").First(&siswa, siswaID).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"pesan": "user ga ketemu"})
			c.Abort()
			return
		}
		
		if siswa.Role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"pesan": "akses ditolak, khusus admin"})
			c.Abort()
			return
		}
		
		c.Next()
	}
}
```

---

### Routes Update

Update `routes/routes.go`:
```go
// tambah di akhir function SetupRoutes

// ===== MERCH ROUTES (SISWA) =====
merch := api.Group("/merch").Use(middleware.AuthMiddleware())
{
	merch.GET("", controllers.GetMerchList)              // list merch
	merch.GET("/:id", controllers.GetMerchDetail)        // detail merch
	merch.POST("/beli", controllers.BeliMerch)           // beli merch
	merch.GET("/riwayat", controllers.GetRiwayatPesananMerch) // riwayat pesanan
}

// ===== ADMIN ROUTES =====
admin := api.Group("/admin").Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
{
	// Kelola Merch
	admin.POST("/merch", controllers.AdminCreateMerch)      // buat merch
	admin.PUT("/merch/:id", controllers.AdminUpdateMerch)   // update merch
	admin.DELETE("/merch/:id", controllers.AdminDeleteMerch) // hapus merch
	
	// Kelola Pesanan
	admin.GET("/pesanan", controllers.AdminGetPesananMerch) // list pesanan
	admin.PUT("/pesanan/:id", controllers.AdminUpdateStatusPesanan) // update status
	
	// Kelola Aksi Sosial
	admin.GET("/sosial", controllers.AdminGetAksiSosial)    // list aksi sosial
	admin.PUT("/sosial/:id", controllers.AdminUpdateAksiSosial) // approve/reject
	
	// Dashboard Stats
	admin.GET("/stats", controllers.AdminGetStats)          // statistik
}
```

---

## ✅ Checklist Backend

- [ ] Update model Siswa (tambah field `role`)
- [ ] Buat model Merch
- [ ] Buat model PesananMerch
- [ ] Auto migrate di database.go
- [ ] Buat merch_controller.go
- [ ] Buat admin_controller.go
- [ ] Buat admin middleware
- [ ] Update routes.go
- [ ] Test endpoints dengan Postman

---

## 📱 Flutter Implementation (Basic)

### Screens Needed:
1. `merch_screen.dart` - Katalog merch (siswa)
2. `detail_merch_screen.dart` - Detail & form beli
3. `riwayat_merch_screen.dart` - Riwayat pesanan
4. `admin_dashboard_screen.dart` - Dashboard admin
5. `admin_merch_screen.dart` - Kelola merch
6. `admin_pesanan_screen.dart` - Kelola pesanan

Mau gua buatin implementasi lengkapnya ga? 🚀
