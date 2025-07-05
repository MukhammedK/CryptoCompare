package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/MukhammedK/cryptocompare/internal/models"
	"github.com/MukhammedK/cryptocompare/internal/parser"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
	"time"

	"github.com/MukhammedK/cryptocompare/internal/db"
	"github.com/go-chi/chi/v5"
)

type CryptoHandler struct {
	DB    *sql.DB
	Redis *redis.Client
}

func NewCryptoHandler(db *sql.DB, redis *redis.Client) *CryptoHandler {
	return &CryptoHandler{DB: db, Redis: redis}
}

func (h *CryptoHandler) GetAllCryptos(w http.ResponseWriter, r *http.Request) {
	cryptos, err := db.GetAllCryptos(h.DB)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(cryptos)
}

func (h *CryptoHandler) GetCryptoBySymbol(w http.ResponseWriter, r *http.Request) {
	symbol := chi.URLParam(r, "symbol")
	crypto, err := db.GetCryptoBySymbol(h.DB, symbol)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(crypto)
}
func (h *CryptoHandler) GetCryptoPrice(w http.ResponseWriter, r *http.Request) {
	symbol := chi.URLParam(r, "symbol")

	crypto, err := db.GetCryptoBySymbol(h.DB, symbol)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	log.Println("🔍 symbol:", symbol)
	log.Println("🔍 coingecko_id:", crypto.CoingeckoID)
	key := fmt.Sprintf("price:%s", crypto.CoingeckoID)
	log.Println("🔍 Redis key:", key)

	val, err := h.Redis.Get(context.Background(), key).Result()
	if err != nil {
		http.Error(w, "price not cached", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"symbol": symbol,
		"price":  val,
	})
}
func (h *CryptoHandler) CompareCryptoPrices(w http.ResponseWriter, r *http.Request) {
	symbol := chi.URLParam(r, "symbol")
	key := fmt.Sprintf("compare:%s", symbol)

	// Кэш
	cached, err := h.Redis.Get(context.Background(), key).Result()
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(cached))
		return
	}

	var results []models.ExchangePrice

	// Попробуем Binance
	if p, err := parser.FetchFromBinance(symbol); err == nil {
		results = append(results, p)
	}

	// Если Binance ничего не дал — пробуем CoinGecko
	if len(results) == 0 {
		if p, err := parser.FetchFromCoinGecko(symbol); err == nil {
			results = append(results, p)
		}
	}

	if len(results) == 0 {
		http.Error(w, "Could not fetch prices from any source", http.StatusInternalServerError)
		return
	}

	// Логируем в БД
	_ = db.SaveCompareResults(h.DB, results)

	// Кэшируем
	jsonData, _ := json.Marshal(results)
	_ = h.Redis.Set(context.Background(), key, jsonData, time.Minute).Err()

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
