package main

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
)

func main() {
	e := echo.New()

	// add routes
	e.GET("/cats/:data", GetCats)
	e.POST("/cats", AddCat)

	e.Logger.Fatal(e.Start(":8000"))
}

// GetCats is Get api, http://localhost:8000/cats?name=a&type=fluffy
func GetCats(c echo.Context) error {
	catName := c.QueryParam("name")
	catType := c.QueryParam("type")
	return c.String(http.StatusOK, fmt.Sprintf("your cat name is: %s\nand cat type is: %s\n", catName, catType))
}

func AddCat(c echo.Context) error {
	type Cat struct {
		Name string `json:"name"`
		Type string `json:"type"`
	}

	cat := Cat{}
	defer c.Request().Body.Close()

	err := json.NewDecoder(c.Request().Body).Decode(&cat)
	if err != nil {
		log.Fatalf("error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.String(http.StatusOK, "We got your cat")
}
