package models

import (
	"time"

	"gorm.io/gorm"
)

// Model Merch - barang yang bisa dituker poin
type Merch struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Nama      string         `gorm:"size:100;not null" json:"nama"`          // nama barang
	Deskripsi string         `gorm:"type:text" json:"deskripsi"`             // deskripsi barang
	FotoURL   string         `gorm:"size:255" json:"foto_url"`               // foto barang
	HargaPoin int            `gorm:"not null" json:"harga_poin"`             // harga pake poin
	Stok      int            `gorm:"default:0" json:"stok"`                  // stok barang
	Status    string         `gorm:"size:20;default:tersedia" json:"status"` // tersedia / habis
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // soft delete
}

// Nama tabel di database
func (Merch) TableName() string {
	return "merch"
}
