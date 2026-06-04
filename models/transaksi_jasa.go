package models

import (
	"time"

	"gorm.io/gorm"
)

// Model Transaksi Jasa - untuk catat transaksi jual-beli jasa
type TransaksiJasa struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	IklanID     uint           `gorm:"not null" json:"iklan_id"`
	PembeliID   uint           `gorm:"not null" json:"pembeli_id"`   // Siswa yang beli jasa
	PenyediaID  uint           `gorm:"not null" json:"penyedia_id"`  // Siswa yang jual jasa
	JumlahPoin  int            `gorm:"not null" json:"jumlah_poin"`  // Berapa poin transaksinya
	Status      string         `gorm:"size:20;default:berjalan" json:"status"` // berjalan/selesai
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"-"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relasi
	Iklan       IklanJasa `gorm:"foreignKey:IklanID;constraint:OnDelete:CASCADE" json:"iklan,omitempty"`
	Pembeli     Siswa     `gorm:"foreignKey:PembeliID;constraint:OnDelete:CASCADE" json:"pembeli,omitempty"`
	Penyedia    Siswa     `gorm:"foreignKey:PenyediaID;constraint:OnDelete:CASCADE" json:"penyedia,omitempty"`
}

// Nama tabel di database
func (TransaksiJasa) TableName() string {
	return "transaksi_jasa"
}
