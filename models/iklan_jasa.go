package models

import (
	"time"

	"gorm.io/gorm"
)

// Model Iklan Jasa - untuk data jasa yang ditawarkan siswa
type IklanJasa struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	PenyediaID  uint           `gorm:"not null" json:"penyedia_id"`     // ID siswa yang posting jasa
	Judul       string         `gorm:"size:150;not null" json:"judul"`
	Kategori    string         `gorm:"size:50;not null" json:"kategori"` // Contoh: Pendidikan, Kreatif, dll
	Deskripsi   string         `gorm:"type:text;not null" json:"deskripsi"`
	HargaPoin   int            `gorm:"not null" json:"harga_poin"`       // Harga dalam poin
	Status      string         `gorm:"size:20;default:tersedia" json:"status"` // tersedia/diambil
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"-"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relasi
	Penyedia    Siswa            `gorm:"foreignKey:PenyediaID;constraint:OnDelete:CASCADE" json:"penyedia,omitempty"`
	Transaksi   []TransaksiJasa  `gorm:"foreignKey:IklanID;constraint:OnDelete:CASCADE" json:"-"`
}

// Nama tabel di database
func (IklanJasa) TableName() string {
	return "iklan_jasa"
}
