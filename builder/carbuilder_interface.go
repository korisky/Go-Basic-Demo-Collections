package builder

type CarBuilder interface {
	SetMake(make string) CarBuilder
	SetModel(model string) CarBuilder
	SetYear(year int) CarBuilder
	SetColor(color string) CarBuilder
	SetEngineSize(engineSize float64) CarBuilder
	Build() Car
}
