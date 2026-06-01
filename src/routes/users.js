
const express = require('express');
const jwt     = require('jsonwebtoken');
const pool    = require('../db');

const router = express.Router();

function cekToken(req, res, next) {
  const header = req.headers['authorization'];
  if (!header) {
    return res.status(401).json({ pesan: 'Token tidak ada, harap login' });
  }

  const token = header.split(' ')[1]; // Ambil bagian setelah "Bearer "
  try {
    const data = jwt.verify(token, process.env.JWT_SECRET);
    req.siswaId = data.id; // Simpan ID siswa di request
    next();
  } catch (error) {
    return res.status(401).json({ pesan: 'Token tidak valid atau sudah expired' });
  }
}


router.get('/profil', cekToken, async (req, res) => {
  try {
    const result = await pool.query(
      'SELECT id, nama, email, kelas, jurusan, saldo_poin, created_at FROM siswa WHERE id = $1',
      [req.siswaId]
    );

    if (result.rows.length === 0) {
      return res.status(404).json({ pesan: 'Siswa tidak ditemukan' });
    }

    res.json({ siswa: result.rows[0] });
  } catch (error) {
    console.error(error);
    res.status(500).json({ pesan: 'Terjadi kesalahan server' });
  }
});


router.get('/riwayat', cekToken, async (req, res) => {
  try {

    const result = await pool.query(
      `SELECT 
         t.id,
         t.jumlah_poin,
         t.status,
         t.created_at,
         j.judul       AS nama_jasa,
         j.kategori,
         p.nama        AS nama_pembeli,
         py.nama       AS nama_penyedia,
         CASE 
           WHEN t.pembeli_id = $1 THEN 'keluar'
           ELSE 'masuk'
         END           AS arah_poin
       FROM transaksi_jasa t
       JOIN iklan_jasa j  ON j.id = t.iklan_id
       JOIN siswa p       ON p.id = t.pembeli_id
       JOIN siswa py      ON py.id = t.penyedia_id
       WHERE t.pembeli_id = $1 OR t.penyedia_id = $1
       ORDER BY t.created_at DESC`,
      [req.siswaId]
    );

    res.json({ riwayat: result.rows });
  } catch (error) {
    console.error(error);
    res.status(500).json({ pesan: 'Terjadi kesalahan server' });
  }
});

module.exports = router;
module.exports.cekToken = cekToken;
