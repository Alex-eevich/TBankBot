package main

import (
	"tbankbot/internal/Graph"
	"tbankbot/internal/config"
	"tbankbot/internal/tbank"
	"time"
)

func main() {
	cfg := config.Load()

	if cfg.Token == "" {

		panic("TINKOFF_TOKEN is not set")
	}

	client := tbank.NewClient(cfg.Token, cfg.BaseURL)

	result, _ := client.Candles("BBG004730N88",
		time.Date(2023, 12, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC),
		tbank.IntervalDay)

	Graph.PrintGraph(result)
}
