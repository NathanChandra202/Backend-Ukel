package models

// kumpulan teks error/sukses biar konsisten sama flutter (field "pesan")
const (
	MsgServerError      = "Terjadi kesalahan pada server. Coba lagi."
	MsgDataTidakLengkap = "Data yang dikirim tidak lengkap atau tidak valid."

	MsgRegisterBerhasil    = "Pendaftaran berhasil. Silakan login."
	MsgRegisterEmailAda    = "Email sudah terdaftar. Gunakan email lain atau login."
	MsgLoginGagal          = "Email atau password salah."
	MsgLoginFieldKosong    = "Email dan password wajib diisi."
	MsgRegisterFieldKosong = "Lengkapi nama, email, password, kelas, dan jurusan."
	MsgTokenGagal          = "Gagal membuat sesi login. Coba lagi."

	MsgTokenTidakAda    = "Sesi login tidak ditemukan. Silakan masuk lagi."
	MsgTokenFormatSalah = "Format token tidak valid. Silakan masuk lagi."
	MsgTokenTidakValid  = "Sesi login kedaluwarsa. Silakan masuk lagi."
	MsgTokenBacaGagal   = "Sesi login tidak valid. Silakan masuk lagi."

	MsgProfilTidakAda    = "Profil siswa tidak ditemukan."
	MsgProfilUpdateGagal = "Gagal memperbarui profil."
	MsgProfilUpdateOk    = "Profil berhasil diperbarui."
	MsgProfilFieldKosong = "Nama, kelas, dan jurusan wajib diisi."
	MsgHapusAkunGagal    = "Gagal menghapus akun."
	MsgHapusAkunOk       = "Akun berhasil dihapus."

	MsgJasaListGagal          = "Gagal memuat daftar jasa."
	MsgJasaTidakAda           = "Jasa tidak ditemukan."
	MsgJasaPostGagal          = "Gagal memposting jasa."
	MsgJasaPostOk             = "Jasa berhasil diposting."
	MsgJasaFieldKosong        = "Judul, kategori, deskripsi, dan harga poin wajib diisi."
	MsgJasaTidakTersedia      = "Jasa tidak tersedia atau sudah diambil orang lain."
	MsgJasaSendiri            = "Anda tidak dapat mengambil jasa milik sendiri."
	MsgJasaPoinKurang         = "Saldo poin tidak cukup untuk mengambil jasa ini."
	MsgJasaAmbilOk            = "Berhasil mengambil jasa."
	MsgTransaksiInvalid       = "Transaksi tidak valid atau sudah selesai."
	MsgTransaksiBukanPenyedia = "Hanya penyedia jasa yang dapat menandai selesai."
	MsgJasaSelesaiOk          = "Jasa selesai. Poin telah ditransfer ke penyedia."

	MsgSosialDeskripsiKosong = "Deskripsi kegiatan wajib diisi."
	MsgSosialFotoWajib       = "Foto bukti wajib dilampirkan (JPG/PNG)."
	MsgSosialSimpanGagal     = "Gagal menyimpan foto."
	MsgSosialUploadOk        = "Aksi sosial berhasil diunggah. Poin bertambah."
	MsgSosialRiwayatGagal    = "Gagal memuat riwayat aksi sosial."

	MsgLeaderboardGagal = "Gagal memuat leaderboard."
	MsgRiwayatGagal     = "Gagal memuat riwayat transaksi."
	MsgExcelGagal       = "Gagal mengunduh file Excel."
)
