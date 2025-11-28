package models

import (
	"time"
)

type Category int

const (
	UNKNOWN Category = iota
	VEGETABLES
	FRUITS
	MEATS
)

type Manufacturer struct {
	Name    string `json:"name"`
	Country string `json:"country"`
	Website string `json:"website"`
}

type Product struct {
	Uuid          string       `json:"uuid"`
	Name          string       `json:"name"`
	Description   string       `json:"description"`
	Price         float64      `json:"price"`
	StockQuantity int64        `json:"stock_quantity"`
	Category      Category     `json:"category"`
	Manufacturer  Manufacturer `json:"manufacturer"`
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
}
