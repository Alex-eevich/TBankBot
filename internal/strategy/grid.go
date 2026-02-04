package strategy

import "tbankbot/internal/models"

type GridConfig struct {
	Levels int
	Step   float64
	Volume float64
}

func BuildGrid(
	price float64,
	direction TrendDirection,
	cfg GridConfig,
) []models.Order {

	var orders []models.Order

	for i := 1; i <= cfg.Levels; i++ {
		offset := float64(i) * cfg.Step

		switch direction {

		case LongOnly:
			orders = append(orders,
				models.Order{
					Price:  price - offset,
					Volume: cfg.Volume,
					Side:   models.Buy,
				},
				models.Order{
					Price:  price + offset,
					Volume: cfg.Volume,
					Side:   models.Sell,
				},
			)

		case ShortOnly:
			orders = append(orders,
				models.Order{
					Price:  price + offset,
					Volume: cfg.Volume,
					Side:   models.Sell,
				},
				models.Order{
					Price:  price - offset,
					Volume: cfg.Volume,
					Side:   models.Buy,
				},
			)
		}
	}

	return orders
}
