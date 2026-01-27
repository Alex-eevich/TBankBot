package main

import (
	"fmt"
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

	shares, err := client.Shares()
	if err != nil {
		panic(err)
	}

	for _, s := range shares[:10] {
		fmt.Printf("%s (%s)\n", s.Ticker, s.Name)
	}

	client.Candles("BBG004730N88",
		time.Date(2023, 12, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2023, 12, 12, 23, 59, 59, 0, time.UTC),
		tbank.IntervalDay)

}
