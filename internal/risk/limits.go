package risk

type RiskManager struct {
	MaxDrawdown float64
	EquityPeak  float64
	EquityNow   float64
	Stopped     bool
}

func NewRiskManager(startCapital, maxDrawdown float64) *RiskManager {
	return &RiskManager{
		MaxDrawdown: maxDrawdown,
		EquityPeak:  startCapital,
		EquityNow:   startCapital,
	}
}

func (r *RiskManager) Update(pnl float64) {
	if r.Stopped {
		return
	}

	r.EquityNow += pnl

	if r.EquityNow > r.EquityPeak {
		r.EquityPeak = r.EquityNow
	}

	drawdown := (r.EquityPeak - r.EquityNow) / r.EquityPeak
	if drawdown >= r.MaxDrawdown {
		r.Stopped = true
	}
}

func (r *RiskManager) Allowed() bool {
	return !r.Stopped
}
