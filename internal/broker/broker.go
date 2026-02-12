package broker

type Side int

const (
	Buy Side = iota
	Sell
)

type Broker interface {
	GetBalance() (float64, error)
	GetPosition(figi string) (float64, error)
	PlaceMarketOrder(figi string, qty int64, side Side) error
}
