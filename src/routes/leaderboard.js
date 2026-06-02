

const express = require('express');
const pool    = require('../db');
const { cekToken } = require('./users');

const router = express.Router();
router.get('/', cekToken, async (req, res) => {
  try {
    const result = await pool.query(
      `SELECT 
         s.id,
         s.nama,
         s.kelas,
         s.jurusan,
         s.saldo_poin,
         -- Hitung berapa kali jasa yang disediakan sudah selesai
         COUNT(t.id) AS jumlah_jasa_selesai,
         -- Total poin yang pernah diterima
         COALESCE(SUM(CASE WHEN t.status = 'selesai' THEN t.jumlah_poin ELSE 0 END), 0) AS total_poin_diterima
       FROM siswa s
       LEFT JOIN transaksi_jasa t ON t.penyedia_id = s.id
       GROUP BY s.id
       ORDER BY total_poin_diterima DESC, s.saldo_poin DESC
       LIMIT 20`,
    );

    res.json({ data: result.rows });
  } catch (error) {
    console.error(error);
    res.status(500).json({ pesan: 'Terjadi kesalahan server' });
  }
});

module.exports = router;
