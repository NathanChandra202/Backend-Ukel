

const express   = require('express');
const pool      = require('../db');
const { cekToken } = require('./users');

const router = express.Router();


router.get('/', cekToken, async (req, res) => {
  try {
    const result = await pool.query(
      `SELECT 
         j.id,
         j.judul,
         j.kategori,
         j.deskripsi,
         j.harga_poin,
         j.status,
         j.created_at,
         s.nama     AS nama_penyedia,
         s.kelas    AS kelas_penyedia
       FROM iklan_jasa j
       JOIN siswa s ON s.id = j.penyedia_id
       WHERE j.status = 'tersedia'
       ORDER BY j.created_at DESC`
    );

    res.json({ data: result.rows });
  } catch (error) {
    console.error(error);
    res.status(500).json({ pesan: 'Terjadi kesalahan server' });
  }
});


router.get('/milik-saya', cekToken, async (req, res) => {
  try {
    const result = await pool.query(
      `SELECT 
         j.*,
         s.nama AS nama_penyedia
       FROM iklan_jasa j
       JOIN siswa s ON s.id = j.penyedia_id
       WHERE j.penyedia_id = $1
       ORDER BY j.created_at DESC`,
      [req.siswaId]
    );

    res.json({ data: result.rows });
  } catch (error) {
    console.error(error);
    res.status(500).json({ pesan: 'Terjadi kesalahan server' });
  }
});

router.get('/:id', cekToken, async (req, res) => {
  const { id } = req.params;
  try {
    const result = await pool.query(
      `SELECT 
         j.*,
         s.nama    AS nama_penyedia,
         s.kelas   AS kelas_penyedia,
         s.jurusan AS jurusan_penyedia
       FROM iklan_jasa j
       JOIN siswa s ON s.id = j.penyedia_id
       WHERE j.id = $1`,
      [id]
    );

    if (result.rows.length === 0) {
      return res.status(404).json({ pesan: 'Jasa tidak ditemukan' });
    }

    res.json({ data: result.rows[0] });
  } catch (error) {
    console.error(error);
    res.status(500).json({ pesan: 'Terjadi kesalahan server' });
  }
});


router.post('/', cekToken, async (req, res) => {
  const { judul, kategori, deskripsi, harga_poin } = req.body;

  if (!judul || !kategori || !deskripsi || !harga_poin) {
    return res.status(400).json({ pesan: 'Semua field harus diisi' });
  }

  try {
    const result = await pool.query(
      `INSERT INTO iklan_jasa (penyedia_id, judul, kategori, deskripsi, harga_poin)
       VALUES ($1, $2, $3, $4, $5)
       RETURNING *`,
      [req.siswaId, judul, kategori, deskripsi, harga_poin]
    );

    res.status(201).json({
      pesan: 'Iklan jasa berhasil dibuat!',
      jasa: result.rows[0],
    });
  } catch (error) {
    console.error(error);
    res.status(500).json({ pesan: 'Terjadi kesalahan server' });
  }
});


router.post('/:id/ambil', cekToken, async (req, res) => {
  const { id } = req.params;

  try {
    const jasaResult = await pool.query(
      'SELECT * FROM iklan_jasa WHERE id = $1',
      [id]
    );

    if (jasaResult.rows.length === 0) {
      return res.status(404).json({ pesan: 'Jasa tidak ditemukan' });
    }

    const jasa = jasaResult.rows[0];

    if (jasa.status !== 'tersedia') {
      return res.status(400).json({ pesan: 'Jasa ini sudah diambil orang lain' });
    }

    if (jasa.penyedia_id === req.siswaId) {
      return res.status(400).json({ pesan: 'Kamu tidak bisa ambil jasa milik sendiri' });
    }

    const pembeliResult = await pool.query(
      'SELECT saldo_poin FROM siswa WHERE id = $1',
      [req.siswaId]
    );
    const pembeli = pembeliResult.rows[0];

    if (pembeli.saldo_poin < jasa.harga_poin) {
      return res.status(400).json({ pesan: 'Saldo poin tidak cukup' });
    }

    await pool.query(
      'UPDATE siswa SET saldo_poin = saldo_poin - $1 WHERE id = $2',
      [jasa.harga_poin, req.siswaId]
    );

    await pool.query(
      "UPDATE iklan_jasa SET status = 'diambil' WHERE id = $1",
      [id]
    );

    const transaksiResult = await pool.query(
      `INSERT INTO transaksi_jasa (iklan_id, pembeli_id, penyedia_id, jumlah_poin)
       VALUES ($1, $2, $3, $4)
       RETURNING *`,
      [id, req.siswaId, jasa.penyedia_id, jasa.harga_poin]
    );

    res.json({
      pesan: `Berhasil ambil jasa! ${jasa.harga_poin} poin dikurangi.`,
      transaksi: transaksiResult.rows[0],
    });
  } catch (error) {
    console.error(error);
    res.status(500).json({ pesan: 'Terjadi kesalahan server' });
  }
});


router.post('/:id/selesai', cekToken, async (req, res) => {
  const { id } = req.params; // id transaksi

  try {
    // Ambil detail transaksi
    const transaksiResult = await pool.query(
      'SELECT * FROM transaksi_jasa WHERE id = $1',
      [id]
    );

    if (transaksiResult.rows.length === 0) {
      return res.status(404).json({ pesan: 'Transaksi tidak ditemukan' });
    }

    const transaksi = transaksiResult.rows[0];

    // Hanya pembeli yang bisa menandai selesai
    if (transaksi.pembeli_id !== req.siswaId) {
      return res.status(403).json({ pesan: 'Hanya pembeli yang bisa menandai selesai' });
    }

    if (transaksi.status === 'selesai') {
      return res.status(400).json({ pesan: 'Transaksi sudah selesai' });
    }

    // Kirim poin ke penyedia jasa
    await pool.query(
      'UPDATE siswa SET saldo_poin = saldo_poin + $1 WHERE id = $2',
      [transaksi.jumlah_poin, transaksi.penyedia_id]
    );

    // Update status transaksi menjadi selesai
    await pool.query(
      "UPDATE transaksi_jasa SET status = 'selesai' WHERE id = $1",
      [id]
    );

    res.json({
      pesan: `Jasa selesai! ${transaksi.jumlah_poin} poin dikirim ke penyedia.`,
    });
  } catch (error) {
    console.error(error);
    res.status(500).json({ pesan: 'Terjadi kesalahan server' });
  }
});

module.exports = router;
