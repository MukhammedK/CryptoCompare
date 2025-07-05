package main

import (
	"github.com/MukhammedK/cryptocompare/internal/db"
	"github.com/MukhammedK/cryptocompare/internal/handlers"
	"github.com/MukhammedK/cryptocompare/internal/parser"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"

	"net/http"
	"time"
)

func main() {
	// Подключение к БД
	pg, err := db.ConnectPostgres()
	if err != nil {
		log.Fatalf("DB error: %v", err)
	}
	redisClient := db.ConnectRedis()

	err = parser.FetchAndCachePrices(redisClient, []string{"bitcoin"})
	if err != nil {
		log.Println("🔥 Manual fetch error:", err)
	} else {
		log.Println("✅ Manual fetch success")
	}

	go func() {
		ids := []string{"bitcoin", "ethereum", "solana"}
		for {
			err := parser.FetchAndCachePrices(redisClient, ids)
			if err != nil {
				log.Println("❌ Fetch error:", err)
			}
			time.Sleep(1 * time.Minute)
		}
	}()
	handler := handlers.NewCryptoHandler(pg, redisClient)
	// Роутинг
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/cryptos/{symbol}/price", handler.GetCryptoPrice)

	// Передаём БД в обработчики

	r.Get("/cryptos", handler.GetAllCryptos)
	r.Get("/cryptos/{symbol}", handler.GetCryptoBySymbol)
	r.Get("/cryptos/{symbol}/compare", handler.CompareCryptoPrices)
	r.Get("/audit", handler.GetAudit)

	fs := http.FileServer(http.Dir("./frontend-static"))
	r.Handle("/*", fs)

	log.Println("🚀 Server running on :8080")
	http.ListenAndServe(":8080", r)
}
