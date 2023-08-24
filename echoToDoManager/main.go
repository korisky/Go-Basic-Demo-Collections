package main

import (
	"example/todomanager/todo"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

func main() {

	// get instance
	tm := todo.NewTodoManager()

	// echo server's startup -> only 2 lines
	e := echo.New()
	e.Use(middleware.Logger())

	e.GET("/", func(c echo.Context) error {
		todos := tm.GetAll()
		return c.JSON(http.StatusOK, todos) // really simple json back stuff
	})

	e.Start(":8888")
}
