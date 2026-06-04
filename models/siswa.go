package models

import (
	"time"

	"gorm.io/gorm"
)

// Model Siswa - untuk data user/siswa
type Siswa struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Nama      string         `gorm:"size:100;not null" json:"nama"`
	Email     string         `gorm:"size:150;unique;not null" json:"email"`
	Password  string         `gorm:"size:255;not null" json:"-"`
	Kelas     string         `gorm:"size:50;not null" json:"kelas"`
	Jurusan   string         `gorm:"size:100;not null" json:"jurusan"`
	SaldoPoin int            `gorm:"default:7" json:"saldo_poin"`
	Role      string         `gorm:"size:20;default:siswa" json:"role"` // role: siswa / admin
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relasi ke tabel lain
	IklanJasa     []IklanJasa     `gorm:"foreignKey:PenyediaID;constraint:OnDelete:CASCADE" json:"-"`
	TransaksiBeli []TransaksiJasa `gorm:"foreignKey:PembeliID;constraint:OnDelete:CASCADE" json:"-"`
	TransaksiJual []TransaksiJasa `gorm:"foreignKey:PenyediaID;constraint:OnDelete:CASCADE" json:"-"`
	AksiSosial    []AksiSosial    `gorm:"foreignKey:SiswaID;constraint:OnDelete:CASCADE" json:"-"`
}

// Nama tabel di database
func (Siswa) TableName() string {
	return "siswa"
}
