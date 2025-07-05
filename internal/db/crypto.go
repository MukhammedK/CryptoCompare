package db

import (
	"context"
	"database/sql"
	_ "database/sql"
	"github.com/MukhammedK/cryptocompare/internal/models"
)

func InsertCrypto(db *sql.DB, c models.Crypto) error {
	_, err := db.ExecContext(context.Background(),
		`INSERT INTO cryptos (symbol, name) VALUES ($1, $2) ON CONFLICT (symbol) DO NOTHING`,
		c.Symbol, c.Name)
	return err
}

func GetAllCryptos(db *sql.DB) ([]models.Crypto, error) {
	rows, err := db.QueryContext(context.Background(),
		`SELECT id, symbol, name FROM cryptos`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cryptos []models.Crypto
	for rows.Next() {
		var c models.Crypto
		if err := rows.Scan(&c.ID, &c.Symbol, &c.Name); err != nil {
			return nil, err
		}
		cryptos = append(cryptos, c)
	}
	return cryptos, nil
}

func GetCryptoBySymbol(db *sql.DB, symbol string) (*models.Crypto, error) {
	var c models.Crypto
	row := db.QueryRowContext(context.Background(),
		`SELECT id, symbol, name, coingecko_id FROM cryptos WHERE symbol = $1`, symbol)
	err := row.Scan(&c.ID, &c.Symbol, &c.Name, &c.CoingeckoID)
	if err != nil {
		return nil, err
	}
	return &c, nil
}
