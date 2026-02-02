package engine

import "time"

type Candle struct {
	Time   time.Time
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume int64
}

type MarketDataProvider interface {
	GetCandles(figi string, from, to time.Time) ([]Candle, error)
}

type TrendDirection int

const (
	NoTrend TrendDirection = iota
	Long
	Short
)

type TrendStrategy interface {
	Direction(candles []Candle) TrendDirection
}

type GridOrder struct {
	Price  float64
	Volume int64
	IsBuy  bool
}

type GridStrategy interface {
	BuildGrid(
		price float64,
		direction TrendDirection,
	) []GridOrder
}

type RiskManager interface {
	AllowTrade() bool
	AllowNewPosition(value float64) bool
}

type Broker interface {
	LastPrice(figi string) (float64, error)
	PlaceOrder(order GridOrder) error
	CancelAll(figi string) error
}

type Engine struct {
	Market MarketDataProvider
	Trend  TrendStrategy
	Grid   GridStrategy
	Risk   RiskManager
	Broker Broker

	Figi string
}

func (e *Engine) RunOnce() error {
	if !e.Risk.AllowTrade() {
		_ = e.Broker.CancelAll(e.Figi)
		return nil
	}

	to := time.Now()
	from := to.Add(-200 * time.Hour)

	candles, err := e.Market.GetCandles(e.Figi, from, to)
	if err != nil || len(candles) < 50 {
		return err
	}

	direction := e.Trend.Direction(candles)
	if direction == NoTrend {
		_ = e.Broker.CancelAll(e.Figi)
		return nil
	}

	price, err := e.Broker.LastPrice(e.Figi)
	if err != nil {
		return err
	}

	orders := e.Grid.BuildGrid(price, direction)

	for _, order := range orders {
		orderValue := float64(order.Volume) * order.Price

		if !e.Risk.AllowNewPosition(orderValue) {
			continue
		}

		_ = e.Broker.PlaceOrder(order)
	}

	return nil
}
