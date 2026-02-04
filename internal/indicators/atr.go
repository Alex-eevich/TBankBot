package indicators

import "math"

func ATR(highs, lows, closes []float64, period int) []float64 {
	if len(highs) < period+1 {
		return nil
	}

	tr := make([]float64, len(highs))
	atr := make([]float64, len(highs))

	// True Range
	for i := 1; i < len(highs); i++ {
		h_l := highs[i] - lows[i]
		h_c := math.Abs(highs[i] - closes[i-1])
		l_c := math.Abs(lows[i] - closes[i-1])

		tr[i] = math.Max(h_l, math.Max(h_c, l_c))
	}

	// первые period-1 значений — ATR не существует
	for i := 0; i < period-1; i++ {
		atr[i] = math.NaN()
	}

	// начальное ATR = SMA(TR)
	var sum float64
	for i := 1; i <= period; i++ {
		sum += tr[i]
	}
	atr[period] = sum / float64(period)

	// EMA-like сглаживание
	for i := period + 1; i < len(tr); i++ {
		atr[i] = (atr[i-1]*(float64(period-1)) + tr[i]) / float64(period)
	}

	return atr
}
