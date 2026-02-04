package sim

import "tbankbot/internal/models"

func ExecuteOrder(p *Portfolio, o *models.Order) {
	cost := o.Price * o.Volume

	switch o.Side {

	case models.Buy:
		if p.Cash >= cost {
			p.Cash -= cost
			p.PositionQty += o.Volume
		}

	case models.Sell:
		if p.PositionQty >= o.Volume {
			p.Cash += cost
			p.PositionQty -= o.Volume
		}
	}
}
