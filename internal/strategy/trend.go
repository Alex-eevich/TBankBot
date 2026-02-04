package strategy

type TrendDirection int

const (
	NoTrade TrendDirection = iota
	LongOnly
	ShortOnly
)

type TrandFilter interface {
	Direct() TrendDirection
}

func (t *EMATrend) Direction() TrendDirection {
	return t.DirectionValue
}
