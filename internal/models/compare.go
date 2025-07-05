package models

import "time"

type CompareAudit struct {
	ID          int
	Symbol      string
	Exchange    string
	Price       float64
	RequestedAt time.Time
}
