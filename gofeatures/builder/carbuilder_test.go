package builder

import (
	"fmt"
	"testing"
)

// Test_CarBuilder need to be ran by using: go test &Test_CarBuilder
func Test_CarBuilder(t *testing.T) {

	electricCarBuilder := &ElectricCarBuilder{}
	gasolineCarBuilder := &GasolineCarBuilder{}

	electricCar := CreateCar(electricCarBuilder)
	gasolineCar := CreateCar(gasolineCarBuilder)

	fmt.Printf("Electic car: %+v\n", electricCar)
	fmt.Printf("Gasoline car: %+v\n", gasolineCar)
}

// CreateCar as a mock
func CreateCar(builder CarBuilder) Car {
	return builder.
		SetMake("Toyota").
		SetModel("Camry").
		SetYear(2016).
		SetColor("white").
		SetEngineSize(2.4).
		Build()
}
