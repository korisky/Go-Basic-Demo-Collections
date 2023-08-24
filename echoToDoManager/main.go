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

	// Echo Groups -> group together routes -> that require authentication and add custom middleware check stuff
	authenticatedGroup := e.Group("/todos", func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// check for auth token
			authorization := c.Request().Header.Get("authorization")
			if authorization != "auth-token" {
				c.Error(echo.ErrUnauthorized)
				return nil
			}
			// token exists
			next(c)
			return nil
		}
	})

	authenticatedGroup.POST("/create", func(c echo.Context) error {
		// bind request's param
		reqBody := todo.CreateTodoRequest{}
		err := c.Bind(&reqBody)
		if err != nil {
			return err
		}

		newTodo := tm.Create(reqBody)

		return c.JSON(http.StatusOK, newTodo)
	})

	e.Start(":8888")
}
