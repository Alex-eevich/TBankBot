package main

import (
	"tbankbot/internal/Graph"
	"tbankbot/internal/config"
	"tbankbot/internal/engine"
	"tbankbot/internal/risk"
	"tbankbot/internal/tbank"
	"time"
)

func main() {
	cfg := config.Load()

	if cfg.Token == "" {

		panic("TINKOFF_TOKEN is not set")
	}

	client := tbank.NewClient(cfg.Token, cfg.BaseURL)

	/*result, _ := client.Candles("BBG004730N88",
	time.Date(2023, 12, 1, 0, 0, 0, 0, time.UTC),
	time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC),
	tbank.IntervalDay)*/

	result, _ := client.Candles("BBG004730N88",
		time.Date(2026, 1, 31, 10, 0, 0, 0, time.UTC),
		time.Date(2026, 1, 31, 12, 59, 59, 0, time.UTC),
		tbank.Interval5Sec)

	MarketData := tbank.NewMarketData(result)
	closes := make([]float64, len(MarketData.Closes))
	highs := make([]float64, len(MarketData.Highs))
	lows := make([]float64, len(MarketData.Lows))
	for i, _ := range MarketData.Closes {
		closes[i] = MarketData.Closes[i]
	}
	for i, _ := range MarketData.Highs {
		highs[i] = MarketData.Highs[i]
	}
	for i, _ := range MarketData.Lows {
		lows[i] = MarketData.Lows[i]
	}

	engine := &engine.Engine{
		Risk: risk.NewRiskManager(100_000, 0.06),
	}

	engine.Run(highs, lows, closes)

	Graph.PrintGraph(result)
}
