package builder

type GasolineCarBuilder struct {
	Car
}

func (b *GasolineCarBuilder) Build() Car {
	return b.Car
}

func (b *GasolineCarBuilder) SetMake(make string) CarBuilder {
	b.Make = make
	return b
}

func (b *GasolineCarBuilder) SetModel(model string) CarBuilder {
	b.Model = model
	return b
}

func (b *GasolineCarBuilder) SetYear(year int) CarBuilder {
	b.Year = year
	return b
}

func (b *GasolineCarBuilder) SetColor(color string) CarBuilder {
	b.Color = color
	return b
}

func (b *GasolineCarBuilder) SetEngineSize(engineSize float64) CarBuilder {
	b.EnginSize = engineSize
	return b
}
