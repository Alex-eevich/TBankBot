package sim

type OrderSide int

const (
	Buy OrderSide = iota
	Sell
)

type Order struct {
	Price  float64
	Qty    float64
	Side   OrderSide
	Filled bool
}
