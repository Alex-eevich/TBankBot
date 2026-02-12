package engine

import (
	"log"
	"tbankbot/internal/broker"
	"tbankbot/internal/indicators"
	"tbankbot/internal/strategy"
	"tbankbot/internal/tbank"
	"time"
)

type Engine struct {
	Broker   broker.Broker
	Client   *tbank.Client
	Strategy strategy.Strategy
	Figi     string
}

func (e *Engine) Run() {

	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {

		<-ticker.C
		log.Println("New cycle")

		// 1. получаем последние свечи
		from := time.Now().Add(-120 * time.Minute)
		to := time.Now()
		candles, err := e.Client.Candles(
			e.Figi,
			from,
			to,
			tbank.Interval1Min,
		)
		if err != nil {
			log.Println("Error fetching candles", err)
			continue
		}
		if len(candles) < 60 {
			log.Println("Not enough candles")
			continue
		}

		md := tbank.NewMarketData(candles)
		closes := md.Closes
		lows := md.Lows
		highs := md.Highs

		// 2. считаем EMA
		ema20 := indicators.EMA(closes, 20)
		ema50 := indicators.EMA(closes, 50)
		atr := indicators.ATR(highs, lows, closes, 14)
		adx := indicators.ADX(highs, lows, closes, 14)

		// 3. если сигнал → отправляем market order
		balance, _ := e.Broker.GetBalance()
		position, _ := e.Broker.GetPosition(e.Figi)

		log.Printf("Balance: %.2f Position: %.2f", balance, position)

		signal := e.Strategy.Evaluate(closes, position, ema20, ema50, atr, adx)
		e.executeSignal(signal)

		log.Println("Tick...")
	}
}

func (e *Engine) executeSignal(signal strategy.Signal) {

	switch signal {

	case strategy.Enterlong:
		log.Println("ENTER LONG")
		e.Broker.PlaceMarketOrder(e.Figi, 1, broker.Buy)

	case strategy.Exitlong:
		log.Println("EXIT LONG")
		e.Broker.PlaceMarketOrder(e.Figi, 1, broker.Sell)

	case strategy.Hold:
		// ничего
	}
}
