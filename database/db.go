package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() *sql.DB {
	// Ambil konfigurasi database dari environment variables
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	if host == "" || port == "" || user == "" || password == "" || dbname == "" {
		log.Fatal("Environment variable DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME harus di-set")
	}

	psqlInfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	var err error
	DB, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("Gagal membuka koneksi:", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("Gagal ping database:", err)
	}

	fmt.Println("✅ Successfully connected to database")
	return DB
}
