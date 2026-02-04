package Graph

import (
	"image/color"
	"math"
	"tbankbot/internal/indicators"
	"tbankbot/internal/models"
	"tbankbot/internal/tbank"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func PrintGraph(candles []models.Candle) models.MarketData {
	MarketData := tbank.NewMarketData(candles)
	closes := make([]float64, len(MarketData.Closes))
	for i, _ := range MarketData.Closes {
		closes[i] = MarketData.Closes[i]
	}

	ema20 := indicators.EMA(closes, 20)
	ema50 := indicators.EMA(closes, 50)

	graph := plot.New()
	graph.Title.Text = "Graphic of Prices"
	graph.X.Label.Text = "Time"
	graph.Y.Label.Text = "Prices"

	graphPointsOpens := make(plotter.XYs, 0)
	graphPointsHighs := make(plotter.XYs, 0)
	graphPointsLows := make(plotter.XYs, 0)
	graphPointsClose := make(plotter.XYs, 0)
	emaPoints20 := make(plotter.XYs, 0)
	emaPoints50 := make(plotter.XYs, 0)
	for i := 0; i < len(MarketData.Opens); i++ {
		graphPointsOpens = append(graphPointsOpens, plotter.XY{
			X: float64(i),
			Y: float64(MarketData.Opens[i]),
		})
		graphPointsClose = append(graphPointsClose, plotter.XY{
			X: float64(i),
			Y: float64(MarketData.Closes[i]),
		})
		graphPointsHighs = append(graphPointsHighs, plotter.XY{
			X: float64(i) + 0.5,
			Y: float64(MarketData.Highs[i]),
		})
		graphPointsLows = append(graphPointsLows, plotter.XY{
			X: float64(i) + 0.5,
			Y: float64(MarketData.Lows[i]),
		})
	}
	for i := range ema20 {
		if ema20[i] == 0 {
			continue
		}
		if math.IsNaN(ema20[i]) {
			continue
		}

		emaPoints20 = append(emaPoints20, plotter.XY{
			X: float64(i),
			Y: ema20[i],
		})
	}
	for i := range ema50 {
		if ema50[i] == 0 {
			continue
		}
		if math.IsNaN(ema50[i]) {
			continue
		}

		emaPoints50 = append(emaPoints50, plotter.XY{
			X: float64(i),
			Y: ema50[i],
		})
	}

	//openLine, _ := plotter.NewLine(graphPointsOpens)
	closeLine, _ := plotter.NewLine(graphPointsClose)
	ema20line, _ := plotter.NewLine(emaPoints20)
	ema50line, _ := plotter.NewLine(emaPoints50)
	highsPoint, _ := plotter.NewScatter(graphPointsHighs)
	lowsPoint, _ := plotter.NewScatter(graphPointsLows)
	ema20line.Color = color.RGBA{B: 255, A: 255}
	ema50line.Color = color.RGBA{R: 255, A: 255}
	//openLine.Color = color.RGBA{B: 255, A: 255}
	highsPoint.Color = color.RGBA{G: 255, A: 255}
	lowsPoint.Color = color.RGBA{R: 255, A: 255}
	grid := plotter.NewGrid()
	graph.Add(grid /*openLine,*/, closeLine, highsPoint, lowsPoint, ema20line, ema50line)

	graph.Save(100*vg.Inch, 45*vg.Inch, "PrintGraph.png")

	return MarketData
}
