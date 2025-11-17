package models

import "sync"

type Stock struct {
	Symbol       string  `json:"symbol"`
	Price        float64 `json:"price"`
	ChangePct    float64 `json:"change_pct"`
	PreviousPrice float64 `json:"previous_price"`
}

type PriceUpdate struct {
	Symbol string  `json:"symbol"`
	Price  float64 `json:"price"`
}

type Order struct {
	ID       int64   `json:"id"`
	Symbol   string  `json:"symbol"`
	Side     string  `json:"side"` // "buy" or "sell"
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
	Time     int64   `json:"time"`
}

var (
	// in-memory list of orders
	Orders      = make([]Order, 0)
	OrdersMutex sync.Mutex
	nextOrderID int64 = 1
)
