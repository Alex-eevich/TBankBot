package models

import (
	"time"
)

type CandleInterval string

const (
	Interval1Min CandleInterval = "1min"
	Interval5Min CandleInterval = "5min"
	IntervalDay  CandleInterval = "day"
)

type Candle struct {
	Time   time.Time
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume string
}

type MarketData struct {
	Time   []time.Time
	Opens  []float64
	Highs  []float64
	Lows   []float64
	Closes []float64
}

type GetCandlesRequest struct {
	Figi     string `json:"figi"`
	From     string `json:"from"`
	To       string `json:"to"`
	Interval string `json:"interval"`
}

type GetCandlesResponse struct {
	Candles []struct {
		Time string `json:"time"`
		Open struct {
			Units string
			Nano  int32
		} `json:"open"`
		High struct {
			Units string
			Nano  int32
		} `json:"high"`
		Low struct {
			Units string
			Nano  int32
		} `json:"low"`
		Close struct {
			Units string
			Nano  int32
		} `json:"close"`
		Volume string `json:"volume"`
	} `json:"candles"`
}
