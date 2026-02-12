package broker

import "tbankbot/internal/tbank"

type TBankBroker struct {
	client    *tbank.Client
	accountID string
}

type MoneyValue struct {
	Units int64 `json:"units"`
	Nano  int32 `json:"nano"`
}

type PortfolioResponse struct {
	TotalAmountPortfolio MoneyValue `json:"totalAmountPortfolio"`

	Positions []struct {
		Figi     string     `json:"figi"`
		Quantity MoneyValue `json:"quantity"`
	} `json:"positions"`
}

func (m MoneyValue) ToFloat() float64 {
	return float64(m.Units) + float64(m.Nano)/1e9
}

func NewTBankBroker(client *tbank.Client, accountID string) *TBankBroker {
	return &TBankBroker{
		client:    client,
		accountID: accountID,
	}
}

func (b *TBankBroker) GetBalance() (float64, error) {
	var resp PortfolioResponse

	err := b.client.GetPortfolio(b.accountID, &resp)
	if err != nil {
		return 0, err
	}

	return resp.TotalAmountPortfolio.ToFloat(), nil
}

func (b *TBankBroker) GetPosition(figi string) (float64, error) {
	var resp PortfolioResponse
	err := b.client.GetPortfolio(b.accountID, &resp)
	if err != nil {
		return 0, err
	}

	for _, p := range resp.Positions {
		if p.Figi == figi {
			return p.Quantity.ToFloat(), nil
		}
	}

	return 0, nil
}

func (b *TBankBroker) PlaceMarketOrder(figi string, qty int64, side Side) error {

	direction := "ORDER_DIRECTION_BUY"
	if side == Sell {
		direction = "ORDER_DIRECTION_SELL"
	}

	return b.client.PostOrder(b.accountID, figi, qty, direction)
}
