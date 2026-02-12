package strategy

import (
	"math"
)

type EMATrend struct {
	DirectionValue TrendDirection
}

func NewEMATrend(emaFast, emaSlow, adx, closes []float64) *EMATrend {
	/*emaFast := indicators.EMA(closes, 20)
	emaSlow := indicators.EMA(closes, 50)
	adx := indicators.ADX(highs, lows, closes, 14)*/

	if emaFast == nil || emaSlow == nil || adx == nil {
		return &EMATrend{NoTrade}
	}

	i := len(closes) - 1

	if i < 0 ||
		i >= len(adx) ||
		math.IsNaN(emaFast[i]) ||
		math.IsNaN(emaSlow[i]) ||
		math.IsNaN(adx[i]) {
		return &EMATrend{NoTrade}
	}

	if adx[i] < 20 {
		return &EMATrend{NoTrade}
	}

	if emaFast[i] > emaSlow[i] {
		return &EMATrend{LongOnly}
	}

	if emaFast[i] < emaSlow[i] {
		return &EMATrend{ShortOnly}
	}

	return &EMATrend{NoTrade}
}
