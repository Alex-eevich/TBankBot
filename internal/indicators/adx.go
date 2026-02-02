package indicators

func ADX(highs, lows, closes []float64, period int) []float64 {
	if len(highs) < period+1 {
		return nil
	}

	dx := make([]float64, len(highs))

	for i := 1; i < len(highs); i++ {
		upMove := highs[i] - highs[i-1]
		downMove := lows[i-1] - lows[i]

		var plusDM, minusDM float64
		if upMove > downMove && upMove > 0 {
			plusDM = upMove
		}
		if downMove > upMove && downMove > 0 {
			minusDM = downMove
		}

		tr := max(
			highs[i]-lows[i],
			abs(highs[i]-closes[i-1]),
			abs(lows[i]-closes[i-1]),
		)

		if tr == 0 {
			continue
		}

		plusDI := 100 * (plusDM / tr)
		minusDI := 100 * (minusDM / tr)

		dx[i] = 100 * abs(plusDI-minusDI) / (plusDI + minusDI)
	}

	return EMA(dx, period)
}

func abs(v float64) float64 {
	if v < 0 {
		return -v
	}
	return v
}

func max(a, b, c float64) float64 {
	m := a
	if b > m {
		m = b
	}
	if c > m {
		m = c
	}
	return m
}
