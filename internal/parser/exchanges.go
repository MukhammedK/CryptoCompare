package parser

import (
	"encoding/json"
	"fmt"
	"github.com/MukhammedK/cryptocompare/internal/models"
	"net/http"
)

func FetchFromBinance(symbol string) (models.ExchangePrice, error) {
	url := fmt.Sprintf("https://api.binance.com/api/v3/ticker/price?symbol=%sUSDT", symbol)
	resp, err := http.Get(url)
	if err != nil {
		return models.ExchangePrice{}, err
	}
	defer resp.Body.Close()

	var result struct {
		Price string `json:"price"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return models.ExchangePrice{}, err
	}

	var price float64
	fmt.Sscanf(result.Price, "%f", &price)

	return models.ExchangePrice{
		Exchange: "Binance",
		Symbol:   symbol,
		Price:    price,
	}, nil
}
