
const express = require('express');
const multer  = require('multer');
const path    = require('path');
const pool    = require('../db');
const { cekToken } = require('./users');

const router = express.Router();

const storage = multer.diskStorage({

  destination: (req, file, cb) => {
    cb(null, path.join(__dirname, '..', '..', 'uploads'));
  },
  filename: (req, file, cb) => {
    const namaUnik = Date.now() + '-' + file.originalname;
    cb(null, namaUnik);
  },
});

const filterGambar = (req, file, cb) => {
  const tipeGambar = ['image/jpeg', 'image/png', 'image/jpg'];
  if (tipeGambar.includes(file.mimetype)) {
    cb(null, true);
  } else {
    cb(new Error('Hanya file gambar yang diizinkan (JPG/PNG)'), false);
  }
};

const upload = multer({ storage, fileFilter: filterGambar });

router.post('/upload', cekToken, upload.single('foto'), async (req, res) => {
  const { deskripsi } = req.body;

  if (!deskripsi) {
    return res.status(400).json({ pesan: 'Deskripsi kegiatan harus diisi' });
  }

  const fotoUrl = req.file
    ? `http://localhost:${process.env.PORT || 3000}/uploads/${req.file.filename}`
    : null;

  try {
    const result = await pool.query(
      `INSERT INTO aksi_sosial (siswa_id, deskripsi, foto_url)
       VALUES ($1, $2, $3)
       RETURNING *`,
      [req.siswaId, deskripsi, fotoUrl]
    );

    res.status(201).json({
      pesan: 'Laporan berhasil dikirim! Tunggu verifikasi admin.',
      aksi: result.rows[0],
    });
  } catch (error) {
    console.error(error);
    res.status(500).json({ pesan: 'Terjadi kesalahan server' });
  }
});

router.get('/riwayat', cekToken, async (req, res) => {
  try {
    const result = await pool.query(
      `SELECT * FROM aksi_sosial
       WHERE siswa_id = $1
       ORDER BY created_at DESC`,
      [req.siswaId]
    );

    res.json({ data: result.rows });
  } catch (error) {
    console.error(error);
    res.status(500).json({ pesan: 'Terjadi kesalahan server' });
  }
});

module.exports = router;
