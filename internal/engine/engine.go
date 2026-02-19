package engine

import (
	"log"
	"strconv"
	"tbankbot/internal/Graph"
	"tbankbot/internal/broker"
	"tbankbot/internal/indicators"
	"tbankbot/internal/sim"
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

func (e *Engine) Run(accountID, token, baseURL string, client *tbank.Client) {

	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	cycle := 0

	for {
		log.Println("Cycle " + strconv.Itoa(cycle))
		cycle += 1

		// 1. получаем последние свечи
		from := time.Now().Add(-180 * time.Minute)
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

		i := len(closes) - 1
		trend := strategy.NewEMATrend(
			ema20,
			ema50,
			adx,
			closes,
		)

		if trend.Direction() == strategy.NoTrade {
			return
		}

		step := atr[i] * 0.5

		grid := strategy.BuildGrid(
			closes[i],
			trend.Direction(),
			strategy.GridConfig{
				Levels: 5,
				Step:   step,
				Volume: 1,
			},
		)

		for _, order := range grid {
			sim.ExecuteOrder(&order, accountID, token, baseURL)
		}
		if err := client.GetSandboxPortfolio(accountID, token, baseURL); err != nil {
			log.Println("Ошибка:", err)
		}
		Graph.PrintGraph(candles)
		<-ticker.C
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
