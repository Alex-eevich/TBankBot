package tbank

import (
	"image/color"
	"tbankbot/internal/models"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func PrintGraph(candles []models.Candle) models.MarketData {
	MarketData := NewMarketData(candles)

	graph := plot.New()
	graph.Title.Text = "Graphic of Prices"
	graph.X.Label.Text = "Time"
	graph.Y.Label.Text = "Prices"

	graphPointsOpens := make(plotter.XYs, 0)
	graphPointsHighs := make(plotter.XYs, 0)
	graphPointsLows := make(plotter.XYs, 0)
	for i := 0; i < len(MarketData.Opens); i++ {
		graphPointsOpens = append(graphPointsOpens, plotter.XY{
			X: float64(i),
			Y: float64(MarketData.Opens[i]),
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

	openLine, _ := plotter.NewLine(graphPointsOpens)
	highsPoint, _ := plotter.NewScatter(graphPointsHighs)
	lowsPoint, _ := plotter.NewScatter(graphPointsLows)
	openLine.Color = color.RGBA{B: 255, A: 255}
	highsPoint.Color = color.RGBA{G: 255, A: 255}
	lowsPoint.Color = color.RGBA{R: 255, A: 255}
	grid := plotter.NewGrid()
	graph.Add(grid, openLine, highsPoint, lowsPoint)

	graph.Save(10*vg.Inch, 5*vg.Inch, "PrintGraph.png")

	return MarketData
}
