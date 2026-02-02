package strategy

type GridOrder struct {
	Price float64
	Side  string // BUY / SELL
}

type Grid struct {
	Level int
	Step  float64
}

func (g *Grid) build(
	price float64,
	dir TrendDirection,
) []GridOrder {

	orders := []GridOrder{}

	for i := 1; i <= g.Level; i++ {
		switch dir {

		case LongOnly:
			orders = append(orders, GridOrder{
				Price: price - float64(i)*g.Step,
				Side:  "BUY",
			})

		case ShortOnly:
			orders = append(orders, GridOrder{
				Price: price + float64(i)*g.Step,
				Side:  "SELL",
			})
		}
	}

	return orders
}
