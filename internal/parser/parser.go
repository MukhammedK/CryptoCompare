package parser

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/MukhammedK/cryptocompare/internal/models"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
	"strings"
	"time"
)

var ctx = context.Background()

func FetchAndCachePrices(rdb *redis.Client, ids []string) error {
	joined := ""
	for i, id := range ids {
		if i > 0 {
			joined += ","
		}
		joined += id
	}

	url := fmt.Sprintf("https://api.coingecko.com/api/v3/simple/price?ids=%s&vs_currencies=usd", joined)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var data map[string]map[string]float64
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}

	for id, price := range data {
		key := fmt.Sprintf("price:%s", id)
		val := price["usd"]
		err := rdb.Set(ctx, key, val, 5*time.Minute).Err()
		if err != nil {
			log.Println("❌ Redis SET error for", key, ":", err)
		} else {
			log.Println("✅ Cached", key, "=", val)
		}
	}

	return nil
}

func FetchFromCoinGecko(symbol string) (models.ExchangePrice, error) {
	// Преобразуем в формат CoinGecko: BTC → bitcoin, ETH → ethereum
	idMap := map[string]string{
		"BTC": "bitcoin",
		"ETH": "ethereum",
		"SOL": "solana",
	}

	coinID, ok := idMap[strings.ToUpper(symbol)]
	if !ok {
		return models.ExchangePrice{}, fmt.Errorf("unsupported symbol for CoinGecko")
	}

	url := fmt.Sprintf("https://api.coingecko.com/api/v3/simple/price?ids=%s&vs_currencies=usd", coinID)
	resp, err := http.Get(url)
	if err != nil {
		return models.ExchangePrice{}, err
	}
	defer resp.Body.Close()

	var data map[string]map[string]float64
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return models.ExchangePrice{}, err
	}

	price := data[coinID]["usd"]
	return models.ExchangePrice{
		Exchange: "CoinGecko",
		Symbol:   symbol,
		Price:    price,
	}, nil
}
