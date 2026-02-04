package engine

import (
	"log"

	"tbankbot/internal/indicators"
	"tbankbot/internal/risk"
	"tbankbot/internal/sim"
	"tbankbot/internal/strategy"
)

type Engine struct {
	Risk *risk.RiskManager
}

func (e *Engine) Run(
	highs, lows, closes []float64,
) {
	// === индикаторы ===
	atr := indicators.ATR(highs, lows, closes, 14)
	emaFast := indicators.EMA(closes, 20)
	emaSlow := indicators.EMA(closes, 50)
	adx := indicators.ADX(highs, lows, closes, 14)

	// === портфель (paper trading) ===
	portfolio := &sim.Portfolio{
		Cash:      100_000, // стартовый капитал
		MaxEquity: 100_000,
	}

	for i := 50; i < len(closes); i++ {

		// --- риск-менеджер ---
		if !e.Risk.Allowed() {
			log.Println("Trading stopped by risk manager")
			return
		}

		// --- определение тренда ---
		trend := strategy.NewEMATrend(
			highs[:i+1],
			lows[:i+1],
			closes[:i+1],
		)

		// --- если нет тренда — просто обновляем equity ---
		if trend.Direction() == strategy.NoTrade {
			sim.UpdateMetrics(portfolio, closes[i])

			log.Printf(
				"[%d] NO TRADE | EMA20=%.2f EMA50=%.2f ADX=%.2f Equity=%.2f",
				i,
				emaFast[i],
				emaSlow[i],
				adx[i],
				portfolio.Equity,
			)
			continue
		}

		// --- шаг сетки ---
		step := atr[i] * 0.5

		// --- строим grid ---
		grid := strategy.BuildGrid(
			closes[i],
			trend.Direction(),
			strategy.GridConfig{
				Levels: 5,
				Step:   step,
				Volume: 1,
			},
		)

		// --- исполняем grid ---
		for _, order := range grid {
			sim.ExecuteOrder(portfolio, &order)
		}

		// --- обновляем метрики ---
		sim.UpdateMetrics(portfolio, closes[i])

		log.Printf(
			"[%d] Trend=%v Price=%.2f Orders=%d Cash=%.2f Pos=%.2f Equity=%.2f",
			i,
			trend.Direction(),
			closes[i],
			len(grid),
			portfolio.Cash,
			portfolio.PositionQty,
			portfolio.Equity,
		)
	}

	// === финальный отчёт ===
	log.Println("===== SIMULATION FINISHED =====")
	log.Printf("Final Cash: %.2f", portfolio.Cash)
	log.Printf("Final Position: %.4f", portfolio.PositionQty)
	log.Printf("Final Equity: %.2f", portfolio.Equity)
	log.Printf("Max Drawdown: %.2f%%", portfolio.MaxDrawdown*100)
}
