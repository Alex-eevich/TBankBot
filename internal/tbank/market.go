package tbank

import (
	"strconv"
)

type InstrumentsResponse struct {
	Instruments []Instrument `json:"instruments"`
}

type Instrument struct {
	FIGI   string `json:"figi"`
	Ticker string `json:"ticker"`
	Name   string `json:"name"`
}

func (c *Client) Shares() ([]Instrument, error) {
	var resp InstrumentsResponse

	err := c.do(
		"POST",
		"tinkoff.public.invest.api.contract.v1.InstrumentsService/Shares",
		map[string]string{
			"instrumentStatus": "INSTRUMENT_STATUS_BASE",
		},
		&resp,
	)

	return resp.Instruments, err
}

func MoneyToFloat(units string, nano int32) float64 {
	u, err := strconv.ParseInt(units, 10, 64)
	if err != nil {
		return 0
	}
	return float64(u) + float64(nano)/1e9
}
