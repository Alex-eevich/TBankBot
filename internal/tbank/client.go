package tbank

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"tbankbot/internal/models"
	"time"
)

type Client struct {
	token      string
	baseURL    string
	httpClient *http.Client
}

type SandboxAccountsResponse struct {
	Accounts []struct {
		Id string `json:"id"`
	} `json:"accounts"`
}

type SandboxAccountInfo struct {
	AccountId string `json:"accountId"`
	Balance   struct {
		Units string `json:"units"`
		Nano  int32  `json:"nano"`
	} `json:"balance"`
}

type MoneyValue struct {
	Currency string `json:"currency"`
	Units    string `json:"units"`
	Nano     int32  `json:"nano"`
}

// SandboxPortfolioResponse — часть ответа с портфелем
type SandboxPortfolioResponse struct {
	TotalAmountCurrencies *MoneyValue `json:"totalAmountCurrencies"`
	TotalAmountShares     *MoneyValue `json:"totalAmountShares"`
	TotalAmountBonds      *MoneyValue `json:"totalAmountBonds"`
}

type PortfolioResponse struct {
	TotalAmountCurrencies MoneyValue `json:"totalAmountCurrencies"`
	TotalAmountPortfolio  MoneyValue `json:"totalAmountPortfolio"`
	TotalAmountShares     MoneyValue `json:"totalAmountShares"`
}

func NewClient(token, baseURL string) *Client {
	return &Client{
		token:      token,
		baseURL:    baseURL,
		httpClient: &http.Client{},
	}
}

func (c *Client) do(method, path string, body any, out any) error {
	var buf *bytes.Buffer

	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return err
		}
		buf = bytes.NewBuffer(b)
	} else {
		buf = bytes.NewBuffer(nil)
	}

	req, err := http.NewRequest(
		method,
		fmt.Sprintf("%s/%s", c.baseURL, path),
		buf,
	)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("API error: %s", resp.Status)
	}

	if out != nil {
		return json.NewDecoder(resp.Body).Decode(out)
	}

	return nil
}

func (c *Client) GetSandboxAccounts() ([]string, error) {

	var resp SandboxAccountsResponse

	err := c.do(
		"POST",
		"tinkoff.public.invest.api.contract.v1.SandboxService/GetSandboxAccounts",
		map[string]string{},
		&resp,
	)

	if err != nil {
		return nil, err
	}

	var ids []string
	for _, acc := range resp.Accounts {
		ids = append(ids, acc.Id)
	}

	return ids, nil
}

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
			Interval5Sec: "CANDLE_INTERVAL_5_SEC",
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
			Open:   MoneyToFloat(c.Open.Units, c.Open.Nano),
			High:   MoneyToFloat(c.High.Units, c.High.Nano),
			Low:    MoneyToFloat(c.Low.Units, c.Low.Nano),
			Close:  MoneyToFloat(c.Close.Units, c.Close.Nano),
			Volume: c.Volume,
		})
	}
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

func (c *Client) PostOrder(accountID, figi string, qty int64, direction string) error {

	req := map[string]any{
		"figi":      figi,
		"quantity":  qty,
		"direction": direction,
		"accountID": accountID,
		"orderType": "ORDER_TYPE_MARKET",
	}

	return c.do(
		"POST",
		"tinkoff.public.invest.api.contract.v1.OrdersService/PostOrder",
		req,
		nil,
	)
}

func (c *Client) GetPortfolio(accountID string, out any) error {
	return c.do(
		"POST",
		"tinkoff.public.invest.api.contract.v1.OperationsService/GetPortfolio",
		map[string]string{
			"accountId": accountID,
			"currency":  "RUB",
		},
		out,
	)
}

// Создание аккаунта на SandBox
func (c *Client) OpenSandboxAccount() (string, error) {

	var resp struct {
		AccountId string `json:"accountId"`
	}
	err := c.do(
		"POST",
		"tinkoff.public.invest.api.contract.v1.SandboxService/OpenSandboxAccount",
		map[string]string{},
		&resp,
	)
	if err != nil {
		return "", err
	}

	return resp.AccountId, nil
}

// Инфа о счете
func (c *Client) GetSandboxPortfolio(accountID, token, baseURL string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	url := baseURL + "/tinkoff.public.invest.api.contract.v1.SandboxService/GetSandboxPortfolio"

	// Тело запроса
	body := map[string]string{
		"accountId": accountID,
		"currency":  "RUB",
	}
	jsonBody, _ := json.Marshal(body)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	// Выполняем запрос
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status %s", res.Status)
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	//fmt.Println("RAW:", string(bodyBytes))
	var parsed PortfolioResponse
	if err := json.Unmarshal(bodyBytes, &parsed); err != nil {
		return err
	}

	Units := parsed.TotalAmountCurrencies.Units
	Nano := parsed.TotalAmountCurrencies.Nano / 1000000
	Currency := parsed.TotalAmountCurrencies.Currency

	InCycleUnits := parsed.TotalAmountShares.Units
	InCycleNano := parsed.TotalAmountShares.Nano / 1000000

	SumPortfolioUnits := parsed.TotalAmountPortfolio.Units
	SumPortfolioNano := parsed.TotalAmountPortfolio.Nano / 1000000

	log.Printf("Доступные средства: %s.%d %s \n",
		Units,
		Nano,
		Currency,
	)
	log.Printf("Средств в акицях: %s.%d %s \n",
		InCycleUnits,
		InCycleNano,
		Currency,
	)
	log.Printf("Сумма портфеля: %s.%d %s \n",
		SumPortfolioUnits,
		SumPortfolioNano,
		Currency,
	)

	return nil
}

// Пополнение счета
func (c *Client) SandboxPayIn(accountID, token, baseURL, amount string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	url := baseURL + "/tinkoff.public.invest.api.contract.v1.SandboxService/SandboxPayIn"

	requestBody := map[string]interface{}{
		"accountId": accountID,
		"amount": MoneyValue{
			Currency: "RUB",
			Units:    amount,
			Nano:     0,
		},
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %s", res.Status)
	}

	fmt.Println("Счёт успешно пополнен 💰")
	return nil
}
