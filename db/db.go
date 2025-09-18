package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://neondb_owner:npg_eoz9XJMp8LrE@ep-morning-recipe-adn7rvwi-pooler.c-2.us-east-1.aws.neon.tech/neondb?sslmode=require"
	}

	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Ошибка подключения к БД:", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("БД не отвечает:", err)
	}

	log.Println("✅ Подключено к БД!")
}
