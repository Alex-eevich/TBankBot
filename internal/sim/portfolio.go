package sim

type Portfolio struct {
	Cash        float64 // свободные деньги
	PositionQty float64 // количество актива
	AvgPrice    float64 // средняя цена позиции
	Equity      float64 // общая стоимость
	MaxEquity   float64
	MaxDrawdown float64
}
