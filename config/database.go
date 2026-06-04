package config

import (
	"backend/models"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// Sambungkan ke database PostgreSQL
func ConnectDatabase() error {
	// Load file .env
	godotenv.Load()

	// Ambil config dari .env, kalau kosong pakai default
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}
	
	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5432"
	}
	
	user := os.Getenv("DB_USER")
	if user == "" {
		user = "postgres"
	}
	
	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		password = "postgres"
	}
	
	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		dbname = "kontribid"
	}

	// Buat connection string
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Sambung ke database
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		return fmt.Errorf("gagal koneksi ke database: %w", err)
	}

	// Auto buat tabel kalau belum ada
	DB.AutoMigrate(
		&models.Siswa{},
		&models.IklanJasa{},
		&models.TransaksiJasa{},
		&models.AksiSosial{},
		&models.Merch{},
		&models.PesananMerch{},
	)

	log.Printf("Database terhubung: %s@%s:%s/%s", user, host, port, dbname)
	return nil
}
