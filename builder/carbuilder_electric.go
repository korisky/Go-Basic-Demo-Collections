package builder

// ElectricCarBuilder is a specific car builder
type ElectricCarBuilder struct {
	Car
}

func (b *ElectricCarBuilder) Build() Car {
	return b.Car
}

func (b *ElectricCarBuilder) SetModel(model string) CarBuilder {
	b.Model = model
	return b
}

func (b *ElectricCarBuilder) SetYear(year int) CarBuilder {
	b.Year = year
	return b
}

func (b *ElectricCarBuilder) SetColor(color string) CarBuilder {
	b.Color = color
	return b
}

func (b *ElectricCarBuilder) SetEngineSize(engineSize float64) CarBuilder {
	b.EnginSize = engineSize
	return b
}

func (b *ElectricCarBuilder) SetMake(make string) CarBuilder {
	b.Make = make
	return b
}
