package indicators

func EMA(values []float64, period int) []float64 {
	if len(values) < period {
		return nil
	}

	k := 2.0 / float64(period+1)
	ema := make([]float64, len(values))

	// старт с SMA
	var sum float64
	for i := 0; i < period; i++ {
		sum += values[i]
	}
	ema[period-1] = sum / float64(period)

	for i := period; i < len(values); i++ {
		ema[i] = values[i]*k + ema[i-1]*(1-k)
	}

	return ema
}
