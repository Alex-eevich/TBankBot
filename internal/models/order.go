package models

type OrderSide int

const (
	Buy OrderSide = iota
	Sell
)

type Order struct {
	Price  float64
	Volume float64
	Side   OrderSide
}
