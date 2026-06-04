package models

import (
	"time"

	"gorm.io/gorm"
)

// Model Aksi Sosial - untuk upload foto kegiatan sosial (dapat poin)
type AksiSosial struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	SiswaID        uint           `gorm:"not null" json:"siswa_id"`
	Deskripsi      string         `gorm:"type:text;not null" json:"deskripsi"` // Deskripsi kegiatan
	FotoURL        string         `gorm:"size:255" json:"foto_url"`            // Path file foto
	Status         string         `gorm:"size:20;default:menunggu" json:"status"` // menunggu/disetujui
	PoinDiberikan  int            `gorm:"default:5" json:"poin_diberikan"`     // Default +5 poin
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"-"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`

	// Relasi
	Siswa          Siswa `gorm:"foreignKey:SiswaID;constraint:OnDelete:CASCADE" json:"siswa,omitempty"`
}

// Nama tabel di database
func (AksiSosial) TableName() string {
	return "aksi_sosial"
}
