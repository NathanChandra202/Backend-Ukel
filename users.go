package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

type ProfilResponse struct {
	ID        int       `json:"id"`
	Nama      string    `json:"nama"`
	Email     string    `json:"email"`
	Kelas     string    `json:"kelas"`
	Jurusan   string    `json:"jurusan"`
	SaldoPoin int       `json:"saldo_poin"`
	CreatedAt time.Time `json:"created_at"`
}

// ngambil data profil siswa dari db
func ProfilSiswa(c *gin.Context) {
	siswaID := c.GetInt("siswaId")

	var profil ProfilResponse
	err := DB.QueryRow(
		"SELECT id, nama, email, kelas, jurusan, saldo_poin, created_at FROM siswa WHERE id = $1",
		siswaID,
	).Scan(&profil.ID, &profil.Nama, &profil.Email, &profil.Kelas, &profil.Jurusan, &profil.SaldoPoin, &profil.CreatedAt)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"pesan": "Siswa tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"siswa": profil})
}

type RiwayatResponse struct {
	ID           int       `json:"id"`
	JumlahPoin   int       `json:"jumlah_poin"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	NamaJasa     string    `json:"nama_jasa"`
	Kategori     string    `json:"kategori"`
	NamaPembeli  string    `json:"nama_pembeli"`
	NamaPenyedia string    `json:"nama_penyedia"`
	ArahPoin     string    `json:"arah_poin"`
}

// ngambil riwayat transaksi poin masuk / keluar
func RiwayatPoin(c *gin.Context) {
	siswaID := c.GetInt("siswaId")

	rows, err := DB.Query(`
		SELECT 
			t.id, t.jumlah_poin, t.status, t.created_at,
			j.judul AS nama_jasa, j.kategori,
			p.nama AS nama_pembeli, py.nama AS nama_penyedia,
			CASE WHEN t.pembeli_id = $1 THEN 'keluar' ELSE 'masuk' END AS arah_poin
		FROM transaksi_jasa t
		JOIN iklan_jasa j ON j.id = t.iklan_id
		JOIN siswa p ON p.id = t.pembeli_id
		JOIN siswa py ON py.id = t.penyedia_id
		WHERE t.pembeli_id = $1 OR t.penyedia_id = $1
		ORDER BY t.created_at DESC`,
		siswaID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"pesan": "Terjadi kesalahan server"})
		return
	}
	defer rows.Close()

	var riwayat []RiwayatResponse = []RiwayatResponse{}
	for rows.Next() {
		var r RiwayatResponse
		if err := rows.Scan(&r.ID, &r.JumlahPoin, &r.Status, &r.CreatedAt, &r.NamaJasa, &r.Kategori, &r.NamaPembeli, &r.NamaPenyedia, &r.ArahPoin); err == nil {
			riwayat = append(riwayat, r)
		}
	}

	c.JSON(http.StatusOK, gin.H{"riwayat": riwayat})
}

type UpdateProfilInput struct {
	Nama    string `json:"nama" binding:"required"`
	Kelas   string `json:"kelas" binding:"required"`
	Jurusan string `json:"jurusan" binding:"required"`
}

// update profil user kayak nama, kelas, jurusan
func UpdateProfil(c *gin.Context) {
	siswaID := c.GetInt("siswaId")
	var input UpdateProfilInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"pesan": err.Error()})
		return
	}

	_, err := DB.Exec("UPDATE siswa SET nama = $1, kelas = $2, jurusan = $3 WHERE id = $4", input.Nama, input.Kelas, input.Jurusan, siswaID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"pesan": "Gagal mengupdate profil"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"pesan": "Profil berhasil diupdate"})
}

// fungsi buat delete akun permanen
func DeleteAkun(c *gin.Context) {
	siswaID := c.GetInt("siswaId")

	_, err := DB.Exec("DELETE FROM siswa WHERE id = $1", siswaID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"pesan": "Gagal menghapus akun"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"pesan": "Akun berhasil dihapus"})
}

// fungsi buat export excel, pakenya token yang dipassing ke param biar gampang
func ExportRiwayatExcel(c *gin.Context) {
	siswaID := c.GetInt("siswaId")

	rows, err := DB.Query(`
		SELECT 
			t.created_at, j.judul, p.nama, py.nama, t.jumlah_poin, t.status,
			CASE WHEN t.pembeli_id = $1 THEN 'keluar' ELSE 'masuk' END AS arah_poin
		FROM transaksi_jasa t
		JOIN iklan_jasa j ON j.id = t.iklan_id
		JOIN siswa p ON p.id = t.pembeli_id
		JOIN siswa py ON py.id = t.penyedia_id
		WHERE t.pembeli_id = $1 OR t.penyedia_id = $1
		ORDER BY t.created_at DESC`,
		siswaID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"pesan": "Gagal mengambil data"})
		return
	}
	defer rows.Close()

	f := excelize.NewFile()
	sheetName := "RiwayatPoin"
	f.SetSheetName("Sheet1", sheetName)
	
	// Set Headers
	headers := []string{"Tanggal", "Judul Jasa", "Pembeli", "Penyedia", "Jumlah Poin", "Status", "Jenis Transaksi"}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetName, cell, header)
	}

	rowIdx := 2
	for rows.Next() {
		var (
			createdAt time.Time
			judul, pembeli, penyedia, status, arahPoin string
			jumlahPoin int
		)
		if err := rows.Scan(&createdAt, &judul, &pembeli, &penyedia, &jumlahPoin, &status, &arahPoin); err == nil {
			f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowIdx), createdAt.Format("2006-01-02 15:04:05"))
			f.SetCellValue(sheetName, fmt.Sprintf("B%d", rowIdx), judul)
			f.SetCellValue(sheetName, fmt.Sprintf("C%d", rowIdx), pembeli)
			f.SetCellValue(sheetName, fmt.Sprintf("D%d", rowIdx), penyedia)
			f.SetCellValue(sheetName, fmt.Sprintf("E%d", rowIdx), jumlahPoin)
			f.SetCellValue(sheetName, fmt.Sprintf("F%d", rowIdx), status)
			
			jenisTx := "Poin Keluar (Bayar)"
			if arahPoin == "masuk" {
				jenisTx = "Poin Masuk (Pendapatan)"
			}
			f.SetCellValue(sheetName, fmt.Sprintf("G%d", rowIdx), jenisTx)
			rowIdx++
		}
	}

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=riwayat_poin.xlsx")
	
	if err := f.Write(c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"pesan": "Gagal membuat file excel"})
	}
}
