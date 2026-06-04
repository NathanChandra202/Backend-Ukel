-- ========================================
-- MIGRATION SCRIPT: FITUR MERCH + ADMIN
-- ========================================
-- Script ini buat update database yang udah ada
-- Jalanin script ini kalau database lu udah jalan sebelumnya

-- 1. Tambah kolom 'role' ke tabel siswa (buat bedain siswa biasa vs admin)
ALTER TABLE siswa ADD COLUMN IF NOT EXISTS role VARCHAR(20) DEFAULT 'siswa';

-- 2. Update salah satu user jadi admin (ubah id=1 sesuai user lu)
-- UPDATE siswa SET role = 'admin' WHERE id = 1;

-- 3. Tabel 'merch' dan 'pesanan_merch' bakal otomatis kebuat pas lu jalanin backend
--    Karena udah ada AutoMigrate di config/database.go
--    Tapi kalau mau manual, bisa pake script dibawah:

-- Tabel merch (katalog barang)
CREATE TABLE IF NOT EXISTS merch (
  id            SERIAL PRIMARY KEY,
  nama          VARCHAR(100) NOT NULL,
  deskripsi     TEXT,
  foto_url      VARCHAR(255),
  harga_poin    INTEGER NOT NULL,
  stok          INTEGER DEFAULT 0,
  status        VARCHAR(20) DEFAULT 'tersedia',
  created_at    TIMESTAMP DEFAULT NOW(),
  updated_at    TIMESTAMP DEFAULT NOW(),
  deleted_at    TIMESTAMP
);

-- Tabel pesanan_merch (transaksi beli barang)
CREATE TABLE IF NOT EXISTS pesanan_merch (
  id            SERIAL PRIMARY KEY,
  siswa_id      INTEGER REFERENCES siswa(id) ON DELETE CASCADE,
  merch_id      INTEGER REFERENCES merch(id) ON DELETE CASCADE,
  jumlah        INTEGER NOT NULL DEFAULT 1,
  total_poin    INTEGER NOT NULL,
  status        VARCHAR(20) DEFAULT 'pending',
  alamat        TEXT,
  catatan       TEXT,
  created_at    TIMESTAMP DEFAULT NOW(),
  updated_at    TIMESTAMP DEFAULT NOW(),
  deleted_at    TIMESTAMP
);

-- Sample data merch (opsional, buat testing)
INSERT INTO merch (nama, deskripsi, foto_url, harga_poin, stok, status) VALUES
('Pulpen Pilot Hitam', 'Pulpen pilot warna hitam, smooth writing', 'https://via.placeholder.com/150', 10, 50, 'tersedia'),
('Pensil 2B Faber Castell', 'Pensil 2B premium untuk menulis dan menggambar', 'https://via.placeholder.com/150', 5, 100, 'tersedia'),
('Buku Tulis 40 Lembar', 'Buku tulis polos 40 lembar', 'https://via.placeholder.com/150', 15, 30, 'tersedia'),
('Penghapus Steadtler', 'Penghapus putih kualitas bagus', 'https://via.placeholder.com/150', 3, 80, 'tersedia'),
('Penggaris 30cm', 'Penggaris plastik transparan 30cm', 'https://via.placeholder.com/150', 7, 40, 'tersedia');

-- Done! Sekarang jalanin backend_gorm.exe
