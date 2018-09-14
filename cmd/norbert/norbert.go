package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	// Echo instance
	e := echo.New()
	e.Static("/", "./public")

	// Middleware
	e.Use(middleware.Logger())

	// Start server
	e.Logger.Fatal(e.Start(":8000"))
}
