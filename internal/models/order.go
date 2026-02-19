package models

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"log"
	"net/http"
	"time"
)

type OrderSide int

const (
	Buy OrderSide = iota
	Sell
)

type Order struct {
	Price  float64
	Volume float64
	Side   OrderSide
}

type Price struct {
	Nano  int32  `json:"nano"`
	Units string `json:"units"`
}

func PostOrder(accountID, token, baseURL, figi string, Side OrderSide) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var direction string
	if Side == Buy {
		direction = "ORDER_DIRECTION_BUY"
	} else if Side == Sell {
		direction = "ORDER_DIRECTION_SELL"
	}

	url := baseURL + "/tinkoff.public.invest.api.contract.v1.OrdersService/PostOrder"
	requestBody := map[string]interface{}{
		"quantity":     1,
		"direction":    direction,
		"accountId":    accountID,
		"orderType":    "ORDER_TYPE_MARKET",
		"orderId":      uuid.New(),
		"instrumentId": figi,
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
	//body, _ := io.ReadAll(res.Body)
	//fmt.Println("Response:", string(body))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %s", res.Status)
	}

	log.Println("Post Order is completed " + figi + " " + direction)
	return nil
}
