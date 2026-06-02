package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type IklanJasa struct {
	ID           int       `json:"id"`
	Judul        string    `json:"judul"`
	Kategori     string    `json:"kategori"`
	Deskripsi    string    `json:"deskripsi"`
	HargaPoin    int       `json:"harga_poin"`
	PenyediaID   int       `json:"penyedia_id"`
	NamaPenyedia string    `json:"nama_penyedia"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
}

func GetJasa(c *gin.Context) {
	rows, err := DB.Query(`
		SELECT i.id, i.judul, i.kategori, i.deskripsi, i.harga_poin,
		       i.penyedia_id, s.nama, i.status, i.created_at
		FROM iklan_jasa i
		JOIN siswa s ON s.id = i.penyedia_id
		WHERE i.status = 'tersedia'
		ORDER BY i.created_at DESC`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"pesan": MsgJasaListGagal})
		return
	}
	defer rows.Close()

	jasas := []IklanJasa{}
	for rows.Next() {
		var j IklanJasa
		if err := rows.Scan(&j.ID, &j.Judul, &j.Kategori, &j.Deskripsi, &j.HargaPoin,
			&j.PenyediaID, &j.NamaPenyedia, &j.Status, &j.CreatedAt); err == nil {
			jasas = append(jasas, j)
		}
	}
	c.JSON(http.StatusOK, gin.H{"data": jasas})
}

func DetailJasa(c *gin.Context) {
	id := c.Param("id")
	var j IklanJasa
	err := DB.QueryRow(`
		SELECT i.id, i.judul, i.kategori, i.deskripsi, i.harga_poin,
		       i.penyedia_id, s.nama, i.status, i.created_at
		FROM iklan_jasa i
		JOIN siswa s ON s.id = i.penyedia_id
		WHERE i.id = $1`, id).Scan(
		&j.ID, &j.Judul, &j.Kategori, &j.Deskripsi, &j.HargaPoin,
		&j.PenyediaID, &j.NamaPenyedia, &j.Status, &j.CreatedAt)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"pesan": MsgJasaTidakAda})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": j})
}

type BuatJasaInput struct {
	Judul     string `json:"judul" binding:"required"`
	Kategori  string `json:"kategori" binding:"required"`
	Deskripsi string `json:"deskripsi" binding:"required"`
	HargaPoin int    `json:"harga_poin" binding:"required,gt=0"`
}

func BuatJasa(c *gin.Context) {
	siswaID := c.GetInt("siswaId")
	var input BuatJasaInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"pesan": MsgJasaFieldKosong})
		return
	}

	_, err := DB.Exec(`
		INSERT INTO iklan_jasa (penyedia_id, judul, kategori, deskripsi, harga_poin, status)
		VALUES ($1, $2, $3, $4, $5, 'tersedia')`,
		siswaID, input.Judul, input.Kategori, input.Deskripsi, input.HargaPoin)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"pesan": MsgJasaPostGagal})
		return
	}
	c.JSON(http.StatusOK, gin.H{"pesan": MsgJasaPostOk})
}

func AmbilJasa(c *gin.Context) {
	pembeliID := c.GetInt("siswaId")
	jasaID := c.Param("id")

	var penyediaID, hargaPoin int
	err := DB.QueryRow(
		"SELECT penyedia_id, harga_poin FROM iklan_jasa WHERE id = $1 AND status = 'tersedia'",
		jasaID,
	).Scan(&penyediaID, &hargaPoin)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"pesan": MsgJasaTidakTersedia})
		return
	}

	if pembeliID == penyediaID {
		c.JSON(http.StatusBadRequest, gin.H{"pesan": MsgJasaSendiri})
		return
	}

	var saldoPembeli int
	if err := DB.QueryRow("SELECT saldo_poin FROM siswa WHERE id = $1", pembeliID).Scan(&saldoPembeli); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"pesan": MsgServerError})
		return
	}
	if saldoPembeli < hargaPoin {
		c.JSON(http.StatusBadRequest, gin.H{"pesan": MsgJasaPoinKurang})
		return
	}

	tx, err := DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"pesan": MsgServerError})
		return
	}
	defer tx.Rollback()

	if _, err = tx.Exec("UPDATE siswa SET saldo_poin = saldo_poin - $1 WHERE id = $2", hargaPoin, pembeliID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"pesan": MsgServerError})
		return
	}
	if _, err = tx.Exec("UPDATE iklan_jasa SET status = 'diambil' WHERE id = $1", jasaID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"pesan": MsgServerError})
		return
	}
	if _, err = tx.Exec(
		"INSERT INTO transaksi_jasa (iklan_id, pembeli_id, penyedia_id, jumlah_poin) VALUES ($1, $2, $3, $4)",
		jasaID, pembeliID, penyediaID, hargaPoin,
	); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"pesan": MsgServerError})
		return
	}
	if err = tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"pesan": MsgServerError})
		return
	}

	c.JSON(http.StatusOK, gin.H{"pesan": MsgJasaAmbilOk})
}

func SelesaikanJasa(c *gin.Context) {
	siswaID := c.GetInt("siswaId")
	transaksiID := c.Param("id")

	var penyediaID, jumlahPoin int
	err := DB.QueryRow(
		"SELECT penyedia_id, jumlah_poin FROM transaksi_jasa WHERE id = $1 AND status = 'berjalan'",
		transaksiID,
	).Scan(&penyediaID, &jumlahPoin)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"pesan": MsgTransaksiInvalid})
		return
	}

	if penyediaID != siswaID {
		c.JSON(http.StatusForbidden, gin.H{"pesan": MsgTransaksiBukanPenyedia})
		return
	}

	tx, err := DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"pesan": MsgServerError})
		return
	}
	defer tx.Rollback()

	if _, err = tx.Exec("UPDATE siswa SET saldo_poin = saldo_poin + $1 WHERE id = $2", jumlahPoin, penyediaID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"pesan": MsgServerError})
		return
	}
	if _, err = tx.Exec("UPDATE transaksi_jasa SET status = 'selesai' WHERE id = $1", transaksiID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"pesan": MsgServerError})
		return
	}
	if err = tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"pesan": MsgServerError})
		return
	}

	c.JSON(http.StatusOK, gin.H{"pesan": MsgJasaSelesaiOk})
}
