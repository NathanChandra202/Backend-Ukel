# 🛍️ API Documentation - Merch & Admin

## 📌 Endpoint Toko Merch (Siswa)

### 1. GET `/api/merch` - List Merch
**Deskripsi**: Ambil katalog merch yang tersedia

**Headers**:
```
Authorization: Bearer {token}
```

**Response Success (200)**:
```json
{
  "data": [
    {
      "id": 1,
      "nama": "Pulpen Pilot",
      "deskripsi": "Pulpen hitam pilot",
      "foto_url": "https://example.com/foto.jpg",
      "harga_poin": 10,
      "stok": 50,
      "status": "tersedia",
      "created_at": "2026-06-03T10:00:00Z",
      "updated_at": "2026-06-03T10:00:00Z"
    }
  ]
}
```

---

### 2. GET `/api/merch/:id` - Detail Merch
**Deskripsi**: Lihat detail satu merch

**Headers**:
```
Authorization: Bearer {token}
```

**Response Success (200)**:
```json
{
  "data": {
    "id": 1,
    "nama": "Pulpen Pilot",
    "deskripsi": "Pulpen hitam pilot",
    "foto_url": "https://example.com/foto.jpg",
    "harga_poin": 10,
    "stok": 50,
    "status": "tersedia",
    "created_at": "2026-06-03T10:00:00Z",
    "updated_at": "2026-06-03T10:00:00Z"
  }
}
```

---

### 3. POST `/api/merch/beli` - Beli Merch
**Deskripsi**: Tukar poin ke barang (transaksi)

**Headers**:
```
Authorization: Bearer {token}
```

**Request Body**:
```json
{
  "merch_id": 1,
  "jumlah": 2,
  "alamat": "Jl. Kebon Jeruk No. 12, Jakarta",
  "catatan": "Kirim hari Senin ya"
}
```

**Response Success (200)**:
```json
{
  "pesan": "berhasil beli merch! tunggu diproses ya"
}
```

**Response Error (400)**:
```json
{
  "pesan": "poin lo ga cukup bro"
}
```

```json
{
  "pesan": "stok ga cukup bro"
}
```

---

### 4. GET `/api/merch/riwayat` - Riwayat Pesanan
**Deskripsi**: Lihat riwayat pesanan merch siswa

**Headers**:
```
Authorization: Bearer {token}
```

**Response Success (200)**:
```json
{
  "data": [
    {
      "id": 1,
      "siswa_id": 5,
      "merch_id": 1,
      "jumlah": 2,
      "total_poin": 20,
      "status": "pending",
      "alamat": "Jl. Kebon Jeruk No. 12, Jakarta",
      "catatan": "Kirim hari Senin ya",
      "created_at": "2026-06-03T10:00:00Z",
      "updated_at": "2026-06-03T10:00:00Z",
      "merch": {
        "id": 1,
        "nama": "Pulpen Pilot",
        "foto_url": "https://example.com/foto.jpg",
        "harga_poin": 10
      }
    }
  ]
}
```

---

## 🔐 Endpoint Admin

**NOTE**: Semua endpoint admin butuh:
1. Token JWT (`Authorization: Bearer {token}`)
2. Role user harus `admin`

---

### 1. POST `/api/admin/merch` - Buat Merch Baru
**Deskripsi**: Tambah barang baru ke katalog

**Headers**:
```
Authorization: Bearer {token}
```

**Request Body**:
```json
{
  "nama": "Pensil 2B",
  "deskripsi": "Pensil 2B Faber Castell",
  "foto_url": "https://example.com/pensil.jpg",
  "harga_poin": 5,
  "stok": 100
}
```

**Response Success (200)**:
```json
{
  "pesan": "merch berhasil ditambahkan",
  "data": {
    "id": 2,
    "nama": "Pensil 2B",
    "deskripsi": "Pensil 2B Faber Castell",
    "foto_url": "https://example.com/pensil.jpg",
    "harga_poin": 5,
    "stok": 100,
    "status": "tersedia"
  }
}
```

---

### 2. PUT `/api/admin/merch/:id` - Update Merch
**Deskripsi**: Edit data merch

**Headers**:
```
Authorization: Bearer {token}
```

**Request Body**:
```json
{
  "nama": "Pensil 2B Update",
  "deskripsi": "Pensil 2B Faber Castell Premium",
  "foto_url": "https://example.com/pensil2.jpg",
  "harga_poin": 7,
  "stok": 150
}
```

**Response Success (200)**:
```json
{
  "pesan": "merch berhasil diupdate"
}
```

---

### 3. DELETE `/api/admin/merch/:id` - Hapus Merch
**Deskripsi**: Hapus merch dari katalog

**Headers**:
```
Authorization: Bearer {token}
```

**Response Success (200)**:
```json
{
  "pesan": "merch berhasil dihapus"
}
```

---

### 4. GET `/api/admin/pesanan` - List Semua Pesanan
**Deskripsi**: Lihat semua pesanan merch (semua siswa)

**Headers**:
```
Authorization: Bearer {token}
```

**Response Success (200)**:
```json
{
  "data": [
    {
      "id": 1,
      "siswa_id": 5,
      "merch_id": 1,
      "jumlah": 2,
      "total_poin": 20,
      "status": "pending",
      "alamat": "Jl. Kebon Jeruk No. 12, Jakarta",
      "catatan": "Kirim hari Senin ya",
      "created_at": "2026-06-03T10:00:00Z",
      "updated_at": "2026-06-03T10:00:00Z",
      "siswa": {
        "id": 5,
        "nama": "Nathan Update",
        "email": "nathan@example.com",
        "kelas": "XI",
        "jurusan": "PPLG"
      },
      "merch": {
        "id": 1,
        "nama": "Pulpen Pilot",
        "foto_url": "https://example.com/foto.jpg"
      }
    }
  ]
}
```

---

### 5. PUT `/api/admin/pesanan/:id` - Update Status Pesanan
**Deskripsi**: Ubah status pesanan (proses, selesai, batal)

**Headers**:
```
Authorization: Bearer {token}
```

**Request Body**:
```json
{
  "status": "diproses"
}
```

**Status Valid**: `pending`, `diproses`, `selesai`, `dibatalkan`

**Response Success (200)**:
```json
{
  "pesan": "status pesanan diupdate"
}
```

---

### 6. GET `/api/admin/sosial` - List Semua Aksi Sosial
**Deskripsi**: Lihat semua aksi sosial yang di-upload siswa

**Headers**:
```
Authorization: Bearer {token}
```

**Response Success (200)**:
```json
{
  "data": [
    {
      "id": 1,
      "siswa_id": 5,
      "foto_url": "/uploads/foto.jpg",
      "deskripsi": "Bersih-bersih kelas",
      "status": "pending",
      "created_at": "2026-06-03T10:00:00Z",
      "siswa": {
        "id": 5,
        "nama": "Nathan Update",
        "email": "nathan@example.com"
      }
    }
  ]
}
```

---

### 7. PUT `/api/admin/sosial/:id` - Approve/Reject Aksi Sosial
**Deskripsi**: Setujui atau tolak aksi sosial

**Headers**:
```
Authorization: Bearer {token}
```

**Request Body**:
```json
{
  "status": "disetujui"
}
```

**Status Valid**: `disetujui`, `ditolak`

**Response Success (200)**:
```json
{
  "pesan": "status aksi sosial diupdate"
}
```

---

### 8. GET `/api/admin/stats` - Dashboard Stats
**Deskripsi**: Statistik untuk dashboard admin

**Headers**:
```
Authorization: Bearer {token}
```

**Response Success (200)**:
```json
{
  "data": {
    "total_siswa": 150,
    "total_poin_beredar": 5000,
    "total_jasa": 45,
    "total_pesanan_merch": 23
  }
}
```

---

## 🔒 Error Responses

### 401 Unauthorized (Token ga ada/salah)
```json
{
  "pesan": "user ga ketemu"
}
```

### 403 Forbidden (Bukan admin)
```json
{
  "pesan": "akses ditolak, khusus admin aja cuy"
}
```

### 400 Bad Request
```json
{
  "pesan": "data ga lengkap"
}
```

### 404 Not Found
```json
{
  "pesan": "merch ga ketemu"
}
```

### 500 Internal Server Error
```json
{
  "pesan": "gagal load katalog merch"
}
```

---

## 📝 Notes

1. **Cara jadi Admin**: Update manual di database:
   ```sql
   UPDATE siswa SET role = 'admin' WHERE id = 1;
   ```

2. **Flow Transaksi Merch**:
   - Siswa pilih merch & jumlah
   - Sistem cek stok & saldo poin
   - Kalau cukup: potong poin + kurangi stok + buat pesanan
   - Admin proses pesanan di dashboard

3. **Flow Aksi Sosial**:
   - Siswa upload foto aksi sosial
   - Admin review & approve/reject
   - Kalau approve: siswa dapat poin (+5)

4. **Status Pesanan**:
   - `pending` = Baru masuk, belum diproses
   - `diproses` = Admin lagi kirim barang
   - `selesai` = Barang udah sampai
   - `dibatalkan` = Dibatalkan admin/system
