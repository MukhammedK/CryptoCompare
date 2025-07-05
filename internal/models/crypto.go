package models

type Crypto struct {
	ID          int64  `json:"id"`
	Symbol      string `json:"symbol"`
	Name        string `json:"name"`
	CoingeckoID string `json:"coingecko_id"`
}
