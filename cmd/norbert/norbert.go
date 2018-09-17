package main

import (
	"log"
	"os"

	"github.com/99designs/gqlgen/handler"
	"github.com/frankh/norbert/cmd/norbert/config"
	"github.com/frankh/norbert/cmd/norbert/graph"
	"github.com/frankh/norbert/cmd/norbert/plugins"
	"github.com/frankh/norbert/cmd/norbert/repository"
	"github.com/frankh/norbert/cmd/norbert/runner"
	"github.com/frankh/norbert/pkg/leader"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	plugins.LoadAll()
	InitMessageQueue()

	dbURI := os.Getenv("NORBERT_DATABASE_URI")
	if dbURI == "" {
		dbURI = "postgres://root@localhost/norbert?sslmode=disable"
	}

	db, err := repository.NewRepository(dbURI)
	if err != nil {
		log.Fatal(err)
	}

	elector, err := leader.NewElector(db.DB.DB, leader.DefaultConfig)
	if err != nil {
		log.Fatal(err)
	}
	elector.Start()
	runner.Start(nc, elector, db, config.Loaded.Checks)

	// Echo instance
	e := echo.New()
	e.Static("/", "./public")

	// Middleware
	e.Use(middleware.Logger())

	e.GET("/query-playground", echo.WrapHandler(
		handler.Playground("Norbert", "/query"),
	))

	e.Match([]string{"GET", "POST"}, "/query", echo.WrapHandler(handler.GraphQL(
		graph.NewExecutableSchema(graph.Config{
			Resolvers: graph.NewResolver(db),
		}),
	)))

	// Start server
	e.Logger.Fatal(e.Start(":8000"))
}
