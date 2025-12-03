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
	Name    string
	Country string
	Website string
}

type Product struct {
	Uuid          string
	Name          string
	Description   string
	Price         float64
	StockQuantity int64
	Category      Category
	Manufacturer  Manufacturer
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type ProductsFilter struct {
	Uuids                 []string
	Names                 []string
	Categories            []Category
	ManufacturerCountries []string
}
