package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.File("index.html")
	})
	e.GET("/local-public-ca.pem", func(c echo.Context) error {
		return c.Attachment("local-public-ca.pem", "local-public-ca.pem")
	})

	e.Logger.Fatal(e.Start(":8081"))
}
