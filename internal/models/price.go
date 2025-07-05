package models

type ExchangePrice struct {
	Symbol   string  `json:"symbol"`
	Price    float64 `json:"price"`
	Exchange string  `json:"exchange"`
}
