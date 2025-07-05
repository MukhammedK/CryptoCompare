package handlers

import (
	"encoding/json"
	"github.com/MukhammedK/cryptocompare/internal/db"
	"net/http"
	"time"
)

func (h *CryptoHandler) GetAudit(w http.ResponseWriter, r *http.Request) {
	symbol := r.URL.Query().Get("symbol")
	fromStr := r.URL.Query().Get("from")
	toStr := r.URL.Query().Get("to")

	// Парсим даты или используем дефолт
	from, _ := time.Parse("2006-01-02", fromStr)
	to, _ := time.Parse("2006-01-02", toStr)
	if from.IsZero() {
		from = time.Now().AddDate(0, 0, -7) // по умолчанию — 7 дней назад
	}
	if to.IsZero() {
		to = time.Now()
	}

	audits, err := db.GetCompareAudit(h.DB, symbol, from, to)
	if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(audits)
}
