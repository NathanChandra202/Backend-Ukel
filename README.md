# Backend Kontrib.ID - GORM Edition

Backend REST API untuk aplikasi Kontrib.ID yang dibangun dengan **Go + Gin + GORM + PostgreSQL**.

## 📁 Struktur Folder (MVC Pattern)

```
backend/
├── config/           # Konfigurasi database
│   └── database.go   # Setup GORM connection
├── controllers/      # Business logic handlers
│   ├── auth_controller.go
│   ├── user_controller.go
│   ├── jasa_controller.go
│   ├── sosial_controller.go
│   └── leaderboard_controller.go
├── middleware/       # HTTP middlewares
│   ├── auth.go       # JWT authentication
│   └── cors.go       # CORS configuration
├── models/           # Database models (GORM)
│   ├── siswa.go
│   ├── iklan_jasa.go
│   ├── transaksi_jasa.go
│   └── aksi_sosial.go
├── routes/           # Route definitions
│   └── routes.go     # Semua endpoint API
├── uploads/          # Folder untuk foto aksi sosial
├── main.go           # Entry point
├── .env              # Environment variables
└── go.mod            # Dependencies
```

## 🚀 Tech Stack

- **Go 1.26.2** - Programming language
- **Gin** - Web framework
- **GORM** - ORM untuk database operations
- **PostgreSQL** - Database
- **JWT** - Authentication
- **Bcrypt** - Password hashing
- **Excelize** - Excel export

## 📦 Dependencies

```go
github.com/gin-gonic/gin           // Web framework
gorm.io/gorm                       // ORM
gorm.io/driver/postgres            // PostgreSQL driver
github.com/golang-jwt/jwt/v5       // JWT authentication
golang.org/x/crypto                // Bcrypt password hashing
github.com/joho/godotenv           // Load .env file
github.com/xuri/excelize/v2        // Excel export
```

## ⚙️ Setup & Installation

### 1. Install Dependencies
```bash
go mod download
```

### 2. Setup Database
Pastikan PostgreSQL sudah jalan, lalu buat database:
```sql
CREATE DATABASE kontribid;
```

GORM akan otomatis membuat tabel saat aplikasi pertama kali jalan (Auto Migration).

### 3. Konfigurasi .env
Edit file `.env`:
```env
PORT=3000

DB_HOST=localhost
DB_PORT=5432
DB_NAME=kontribid
DB_USER=postgres
DB_PASSWORD=your_password

JWT_SECRET=kontribid_rahasia_2026
```

### 4. Run Server
```bash
go run main.go
```

Atau build dulu:
```bash
go build -o backend_gorm.exe
.\backend_gorm.exe
```

Server akan jalan di `http://localhost:3000`

## 📡 API Endpoints

### Health Check
```
GET /api/health
```

### Authentication
```
POST /api/auth/register   # Daftar akun baru
POST /api/auth/login      # Login & dapat token
```

### User Profile (butuh auth)
```
GET  /api/siswa/profil    # Lihat profil
PUT  /api/siswa           # Update profil
DELETE /api/siswa         # Hapus akun
GET  /api/siswa/riwayat   # Riwayat transaksi poin
GET  /api/siswa/export    # Download riwayat (Excel)
```

### Jasa Marketplace (butuh auth)
```
GET  /api/jasa            # List semua jasa tersedia
GET  /api/jasa/:id        # Detail jasa
POST /api/jasa            # Post jasa baru
POST /api/jasa/:id/ambil  # Ambil/beli jasa
POST /api/jasa/:id/selesai # Tandai jasa selesai
```

### Aksi Sosial (butuh auth)
```
POST /api/sosial/upload   # Upload foto kegiatan (+5 poin)
GET  /api/sosial/riwayat  # Riwayat upload
```

### Leaderboard (butuh auth)
```
GET /api/leaderboard      # Top 10 siswa
```

## 🔐 Authentication

Semua endpoint (kecuali `/auth/*` dan `/health`) butuh JWT token di header:

```
Authorization: Bearer <your_jwt_token>
```

Token berlaku **72 jam** setelah login.

## 🗄️ Database Models

### Siswa
- ID, Nama, Email (unique), Password (hashed)
- Kelas, Jurusan
- SaldoPoin (default: 7)

### IklanJasa
- Judul, Kategori, Deskripsi, HargaPoin
- Status: `tersedia`, `diambil`

### TransaksiJasa
- IklanID, PembeliID, PenyediaID
- JumlahPoin
- Status: `berjalan`, `selesai`

### AksiSosial
- Deskripsi, FotoURL
- Status: `menunggu`, `disetujui`
- PoinDiberikan (default: 5)

## 🔄 Keunggulan GORM vs SQL Manual

| Fitur | SQL Manual | GORM |
|---|---|---|
| Query Builder | ❌ Manual string | ✅ Method chaining |
| Migrations | ❌ Manual | ✅ Auto migrate |
| Relations | ❌ Manual JOIN | ✅ Preload() |
| Type Safety | ⚠️ Partial | ✅ Full |
| Transactions | ⚠️ Manual | ✅ tx.Transaction() |
| Code Length | 🔴 Panjang | 🟢 Singkat |

## 🛠️ Development

### Build
```bash
go build -o backend_gorm.exe
```

### Run Tests (TODO)
```bash
go test ./...
```

### Format Code
```bash
go fmt ./...
```

## 📝 Migration dari SQL ke GORM

File lama yang sudah dihapus:
- ❌ `auth.go`, `users.go`, `jasa.go`, dll (diganti dengan controllers/)
- ❌ `db.go` (diganti dengan config/database.go)
- ❌ `middleware.go` (dipindah ke middleware/)
- ❌ `pesan.go` (pesan error langsung di controller)

File baru:
- ✅ `models/*` - GORM models
- ✅ `controllers/*` - Business logic
- ✅ `routes/routes.go` - Centralized routing
- ✅ `middleware/*` - Reusable middlewares
- ✅ `config/*` - Database config

## 🚨 Troubleshooting

### Error: Port already in use
```bash
# Matikan server lama dulu
taskkill /F /IM backend.exe
```

### Error: Database connection failed
- Pastikan PostgreSQL jalan
- Cek username/password di `.env`
- Pastikan database `kontribid` sudah dibuat

### Error: Migration failed
Kalau ada constraint conflict, drop database lalu buat ulang:
```sql
DROP DATABASE kontribid;
CREATE DATABASE kontribid;
```

## 📄 License

Ini project uji kelayakan, bebas dipakai untuk pembelajaran.

---

**Dibuat dengan ❤️ menggunakan Go + GORM**
