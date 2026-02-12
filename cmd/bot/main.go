package main

import (
	"log"
	"tbankbot/internal/config"
	"tbankbot/internal/tbank"
)

func main() {
	cfg := config.Load()

	if cfg.Token == "" {

		panic("TINKOFF_TOKEN is not set")
	}

	client := tbank.NewClient(cfg.Token, cfg.BaseURL)

	/*accountID, err := client.OpenSandboxAccount()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("New sandbox account:", accountID)*/

	accountID, err := client.GetSandboxAccounts()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Sandbox accounts:", accountID[0])

	/*err = client.SandboxPayIn(accountID[0], 100000)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Sandbox funded")*/

	/*result, _ := client.Candles("BBG004730N88",
	time.Date(2026, 1, 31, 10, 0, 0, 0, time.UTC),
	time.Date(2026, 1, 31, 12, 59, 59, 0, time.UTC),
	tbank.Interval5Sec)*/

	/*MarketData := tbank.NewMarketData(result)
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

	engine := &engine.BacktestEngine{
		Risk: risk.NewRiskManager(100_000, 0.06),
	}*/

	/*accountID := "sandbox-account-id"

	strat := &strategy.GridTrendStrategy{
		FastEMA: 20,
		SlowEMA: 50,
	}
	broker := broker.NewTBankBroker(client, accountID)

	engine := &engine.Engine{
		Broker:   broker,
		Client:   client,
		Strategy: strat,
		Figi:     "BBG004730N88",
	}

	engine.Run()

	Graph.PrintGraph(result)*/
}
