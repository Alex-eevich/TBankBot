package strategy

type TrendDirection int

const (
	NoTrade TrendDirection = iota
	LongOnly
	ShortOnly
)

type TrandFilter interface {
	Direction() TrendDirection
}
