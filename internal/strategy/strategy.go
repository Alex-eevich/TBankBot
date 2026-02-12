package strategy

type Signal int

const (
	Hold Signal = iota
	Enterlong
	Exitlong
)

type GridTrendStrategy struct {
	FastEMA int
	SlowEMA int
}

type Strategy interface {
	Evaluate(closes []float64, position float64, ema20 []float64, ema50 []float64, atr []float64, adx []float64) Signal
}

func (s *GridTrendStrategy) Evaluate(closes []float64, position float64, ema20 []float64, ema50 []float64, atr []float64, adx []float64) Signal {

	last := len(closes) - 1
	if last < 50 {
		return Hold
	}

	trending := adx[last] > 20
	crossUp := ema20[last-1] < ema50[last-1] &&
		ema20[last] > ema50[last]
	crossDown := ema20[last-1] > ema50[last-1] &&
		ema20[last] < ema50[last]

	if trending && crossUp && position == 0 {
		return Enterlong
	}
	if crossDown && position > 0 {
		return Exitlong
	}
	return Hold
}
