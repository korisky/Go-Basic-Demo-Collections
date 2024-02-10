package main

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"net/http"
)

func main() {
	e := echo.New()

	// add routes
	e.GET("/cats/:data", GetCats)
	e.POST("/cats", AddCat)
	e.GET("/callback", GetToken)

	// add middle ware
	e.Use(middleware.Logger())

	e.Logger.Fatal(e.Start(":8000"))
}

func GetToken(c echo.Context) error {
	token := c.QueryParam("token")
	fmt.Printf("Got token %v\n", token)
	return c.JSON(http.StatusOK, nil)
}

// GetCats is Get api, curl -X GET 'http://127.0.0.1:8000/cats/json?name=tom&type=fluffy'
func GetCats(c echo.Context) error {
	catName := c.QueryParam("name")
	catType := c.QueryParam("type")
	dataType := c.Param("data")
	if dataType == "string" {
		return c.String(http.StatusOK, fmt.Sprintf("your cat name is : %s\nand cat type is : %s\n", catName, catType))
	} else if dataType == "json" {
		return c.JSON(http.StatusOK, map[string]string{
			"name": catName,
			"type": catType})
	} else {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Please specify the data type as Sting or JSON"})
	}
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
