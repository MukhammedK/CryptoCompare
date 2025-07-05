package db

import (
	"database/sql"
	"github.com/MukhammedK/cryptocompare/internal/models"
	"time"
)

func SaveCompareResults(db *sql.DB, records []models.ExchangePrice) error {
	query := `INSERT INTO compare_audit (symbol, exchange, price) VALUES ($1, $2, $3)`
	for _, r := range records {
		_, err := db.Exec(query, r.Symbol, r.Exchange, r.Price)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetCompareAudit(db *sql.DB, symbol string, from, to time.Time) ([]models.CompareAudit, error) {
	query := `SELECT id, symbol, exchange, price, requested_at
			  FROM compare_audit
			  WHERE ($1 = '' OR symbol = $1)
			  AND requested_at BETWEEN $2 AND $3
			  ORDER BY requested_at DESC
			  LIMIT 100`

	rows, err := db.Query(query, symbol, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var audits []models.CompareAudit
	for rows.Next() {
		var a models.CompareAudit
		err := rows.Scan(&a.ID, &a.Symbol, &a.Exchange, &a.Price, &a.RequestedAt)
		if err != nil {
			return nil, err
		}
		audits = append(audits, a)
	}
	return audits, nil
}
