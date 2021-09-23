package main

import (
	"2021_2_LostPointer/handlers"
	"github.com/labstack/echo"
)

func main() {
	e := echo.New()

	e.GET("/", handlers.GetHomePageHandler())
	e.POST("/signup", handlers.CreateUserHandler())
	e.POST("/signin", handlers.LoginUserHandler())

	e.Logger.Fatal(e.Start(":3030"))
}