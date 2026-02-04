package indicators

import "math"

func ADX(highs, lows, closes []float64, period int) []float64 {
	n := len(highs)
	if n < period*2 {
		return nil
	}

	tr := make([]float64, n)
	plusDM := make([]float64, n)
	minusDM := make([]float64, n)

	// === 1. True Range и Directional Movement ===
	for i := 1; i < n; i++ {
		highDiff := highs[i] - highs[i-1]
		lowDiff := lows[i-1] - lows[i]

		if highDiff > lowDiff && highDiff > 0 {
			plusDM[i] = highDiff
		} else {
			plusDM[i] = 0
		}

		if lowDiff > highDiff && lowDiff > 0 {
			minusDM[i] = lowDiff
		} else {
			minusDM[i] = 0
		}

		tr1 := highs[i] - lows[i]
		tr2 := math.Abs(highs[i] - closes[i-1])
		tr3 := math.Abs(lows[i] - closes[i-1])

		tr[i] = math.Max(tr1, math.Max(tr2, tr3))
	}

	// === 2. Wilder smoothing ===
	smoothedTR := make([]float64, n)
	smoothedPlusDM := make([]float64, n)
	smoothedMinusDM := make([]float64, n)

	var trSum, plusSum, minusSum float64
	for i := 1; i <= period; i++ {
		trSum += tr[i]
		plusSum += plusDM[i]
		minusSum += minusDM[i]
	}

	smoothedTR[period] = trSum
	smoothedPlusDM[period] = plusSum
	smoothedMinusDM[period] = minusSum

	for i := period + 1; i < n; i++ {
		smoothedTR[i] = smoothedTR[i-1] - (smoothedTR[i-1] / float64(period)) + tr[i]
		smoothedPlusDM[i] = smoothedPlusDM[i-1] - (smoothedPlusDM[i-1] / float64(period)) + plusDM[i]
		smoothedMinusDM[i] = smoothedMinusDM[i-1] - (smoothedMinusDM[i-1] / float64(period)) + minusDM[i]
	}

	// === 3. DI+, DI-, DX ===
	dx := make([]float64, n)

	for i := period; i < n; i++ {
		if smoothedTR[i] == 0 {
			dx[i] = math.NaN()
			continue
		}

		diPlus := 100 * (smoothedPlusDM[i] / smoothedTR[i])
		diMinus := 100 * (smoothedMinusDM[i] / smoothedTR[i])

		denom := diPlus + diMinus
		if denom == 0 {
			dx[i] = math.NaN()
			continue
		}

		dx[i] = math.Abs(diPlus-diMinus) / denom * 100
	}

	// === 4. ADX (Wilder smoothing of DX) ===
	adx := make([]float64, n)

	var dxSum float64
	start := period * 2
	for i := period; i < start; i++ {
		dxSum += dx[i]
	}
	adx[start] = dxSum / float64(period)

	for i := start + 1; i < n; i++ {
		adx[i] = (adx[i-1]*(float64(period-1)) + dx[i]) / float64(period)
	}

	return adx
}
