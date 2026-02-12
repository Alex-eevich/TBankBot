package tbank

type CandleInterval string

const (
	Interval5Sec CandleInterval = "5sec"
	Interval1Min CandleInterval = "1min"
	Interval5Min CandleInterval = "5min"
	IntervalDay  CandleInterval = "day"
)
