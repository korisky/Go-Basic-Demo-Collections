package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"os"
	"runtime"
)

var (
	home   = os.Getenv("HOME")
	user   = os.Getenv("USER")
	gopath = os.Getenv("GOPATH")
)

func main() {
	fmt.Println("Now get into Main")
	app := fiber.New()

	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.SendString("Hello, Fiber!")
	})

	app.Listen(":8902")
}

func init() {
	var numCPU = runtime.NumCPU()
	fmt.Printf("\nVar part finished, home:%s\nuser:%s\ngopath:%s\nnumCpu:%d\n",
		home, user, gopath, numCPU)
	fmt.Println("\nNow get into init")
}
