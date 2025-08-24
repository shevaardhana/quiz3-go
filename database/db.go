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
	var psqlInfo string
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL != "" {
		psqlInfo = databaseURL
	} else {
		host := os.Getenv("PGHOST")
		port := os.Getenv("PGPORT")
		user := os.Getenv("PGUSER")
		password := os.Getenv("PGPASSWORD")
		dbname := os.Getenv("PGDATABASE")

		if host == "" || port == "" || user == "" || password == "" || dbname == "" {
			log.Fatal("Environment variable DATABASE_URL atau PGHOST, PGPORT, PGUSER, PGPASSWORD, PGDATABASE harus di-set")
		}

		psqlInfo = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			host, port, user, password, dbname,
		)
	}

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
