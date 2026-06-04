package models

import (
	"time"

	"gorm.io/gorm"
)

// Model PesananMerch - transaksi tukar poin ke merch, diambil di koperasi sekolah
type PesananMerch struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	SiswaID   uint           `gorm:"not null" json:"siswa_id"`
	MerchID   uint           `gorm:"not null" json:"merch_id"`
	Jumlah    int            `gorm:"not null;default:1" json:"jumlah"`
	TotalPoin int            `gorm:"not null" json:"total_poin"`
	Status    string         `gorm:"size:20;default:pending" json:"status"` // pending / siap_ambil / sudah_diambil / dibatalkan
	KodeAmbil string         `gorm:"size:10" json:"kode_ambil"`             // kode unik buat claim di koperasi
	Catatan   string         `gorm:"type:text" json:"catatan"`              // catatan tambahan (opsional)
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relasi
	Siswa Siswa `gorm:"foreignKey:SiswaID" json:"siswa,omitempty"`
	Merch Merch `gorm:"foreignKey:MerchID" json:"merch,omitempty"`
}

func (PesananMerch) TableName() string {
	return "pesanan_merch"
}
