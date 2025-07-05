package main

import (
	"github.com/MukhammedK/cryptocompare/internal/db"
	"log"
)

func main() {
	_, err := db.ConnectPostgres()
	if err != nil {
		log.Fatalf("❌ Не удалось подключиться к PostgreSQL: %v", err)
	}
	log.Println("✅ Подключение к PostgreSQL успешно!")
}
