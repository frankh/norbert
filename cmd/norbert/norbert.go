package main

import (
	"log"

	"github.com/99designs/gqlgen/handler"
	"github.com/frankh/norbert/cmd/norbert/graph"
	"github.com/frankh/norbert/cmd/norbert/models"
	"github.com/frankh/norbert/cmd/norbert/plugins"
	"github.com/frankh/norbert/cmd/norbert/runner"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	plugins.LoadAll()
	log.Println(runner.RunCheck(models.Check{
		CheckRunner: "http",
		Vars: map[string]interface{}{
			"url":      "https://www.google.com",
			"expected": []int{200},
		},
	}))

	// Echo instance
	e := echo.New()
	e.Static("/", "./public")

	// Middleware
	e.Use(middleware.Logger())

	e.GET("/query-playground", echo.WrapHandler(
		handler.Playground("Norbert", "/query"),
	))

	e.POST("/query", echo.WrapHandler(handler.GraphQL(
		graph.NewExecutableSchema(graph.Config{
			Resolvers: graph.NewResolver(),
		}),
	)))

	// Start server
	e.Logger.Fatal(e.Start(":8000"))
}
