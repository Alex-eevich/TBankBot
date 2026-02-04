package sim

func UpdateMetrics(p *Portfolio, marketPrice float64) {
	positionValue := p.PositionQty * marketPrice
	p.Equity = p.Cash + positionValue

	if p.Equity > p.MaxEquity {
		p.MaxEquity = p.Equity
	}

	drawdown := (p.MaxEquity - p.Equity) / p.MaxEquity
	if drawdown > p.MaxDrawdown {
		p.MaxDrawdown = drawdown
	}
}
