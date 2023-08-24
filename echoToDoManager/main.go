package main

import "github.com/labstack/echo/v4"

func main() {

	// echo server's startup -> only 2 lines
	e := echo.New()
	e.Start(":8888")
}
