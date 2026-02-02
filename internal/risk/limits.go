package risk

type RiskManager interface {
	AllowTrade() bool
	AllowNewPosition(positionValue float64) bool
	UpdateEquity(currentEquity float64)
}

type FixedLimits struct {
	StartEquity float64
	MaxDrawdown float64 // например 0.08 = 8%
	MaxExposure float64 // например 0.4 = 40%

	PeakEquity    float64
	CurrentEquity float64
	Exposure      float64
	Blocked       bool
}

func NewFixedLimits(startEquity, maxDD, maxExposure float64) *FixedLimits {
	return &FixedLimits{
		StartEquity:   startEquity,
		PeakEquity:    startEquity,
		CurrentEquity: startEquity,
		MaxDrawdown:   maxDD,
		MaxExposure:   maxExposure,
	}
}

func (r *FixedLimits) UpdateEquity(current float64) {
	r.CurrentEquity = current

	if current > r.PeakEquity {
		r.PeakEquity = current
	}

	drawdown := (r.PeakEquity - current) / r.PeakEquity
	if drawdown >= r.MaxDrawdown {
		r.Blocked = true
	}
}

func (r *FixedLimits) AllowNewPosition(positionValue float64) bool {
	if r.Blocked {
		return false
	}

	newExposure := r.Exposure + positionValue
	return newExposure <= r.StartEquity*r.MaxExposure
}

func (r *FixedLimits) AddExposure(value float64) {
	r.Exposure += value
}

func (r *FixedLimits) ReduceExposure(value float64) {
	r.Exposure -= value
	if r.Exposure < 0 {
		r.Exposure = 0
	}
}
