package engine

/*
import (
	"log"

	"tbankbot/internal/indicators"
	"tbankbot/internal/risk"
	"tbankbot/internal/sim"
	"tbankbot/internal/strategy"
)

type BacktestEngine struct {
	Risk *risk.RiskManager
}

func (e *BacktestEngine) Run(
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

	initial := portfolio.Cash

	for i := 50; i < len(closes); i++ {

		// --- риск-менеджер ---
		if !e.Risk.Allowed() {
			log.Println("Trading stopped by risk manager")
			return
		}

		// --- определение тренда ---
		trend := strategy.NewEMATrend(
			emaFast[:i+1],
			emaSlow[:i+1],
			adx[:i+1],
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

	final := portfolio.Equity

	growth := (final - initial) / initial * 100

	log.Printf("Return: %.2f%%", growth)
}*/
