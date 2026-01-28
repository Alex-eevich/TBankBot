package tbank

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"tbankbot/internal/models"
	"time"
)

func (c *Client) Candles(
	figi string,
	from, to time.Time,
	interval CandleInterval,
) ([]models.Candle, error) {

	reqBody := models.GetCandlesRequest{
		Figi: figi,
		From: from.UTC().Format(time.RFC3339),
		To:   to.UTC().Format(time.RFC3339),
		Interval: map[CandleInterval]string{
			Interval1Min: "CANDLE_INTERVAL_1_MIN",
			Interval5Min: "CANDLE_INTERVAL_5_MIN",
			IntervalDay:  "CANDLE_INTERVAL_DAY",
		}[interval],
	}

	data, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(
		"POST",
		c.baseURL+"/tinkoff.public.invest.api.contract.v1.MarketDataService/GetCandles",
		bytes.NewBuffer(data),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %s", resp.Status)
	}

	var apiResp models.GetCandlesResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}

	var result []models.Candle
	for _, c := range apiResp.Candles {
		t, _ := time.Parse(time.RFC3339, c.Time)

		result = append(result, models.Candle{
			Time:   t,
			Open:   moneyToFloat(c.Open.Units, c.Open.Nano),
			High:   moneyToFloat(c.High.Units, c.High.Nano),
			Low:    moneyToFloat(c.Low.Units, c.Low.Nano),
			Close:  moneyToFloat(c.Close.Units, c.Close.Nano),
			Volume: c.Volume,
		})
	}

	/*for i, _ := range result {
		Open = append(Open, )
	}*/

	return result, nil
}

func NewMarketData(candles []models.Candle) models.MarketData {
	md := models.MarketData{
		Time:   make([]time.Time, 0, len(candles)),
		Opens:  make([]float64, 0, len(candles)),
		Highs:  make([]float64, 0, len(candles)),
		Lows:   make([]float64, 0, len(candles)),
		Closes: make([]float64, 0, len(candles)),
	}

	for _, c := range candles {
		md.Time = append(md.Time, c.Time)
		md.Opens = append(md.Opens, c.Open)
		md.Highs = append(md.Highs, c.High)
		md.Lows = append(md.Lows, c.Low)
		md.Closes = append(md.Closes, c.Close)
	}

	return md
}
