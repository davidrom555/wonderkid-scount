package main

import (
	"fifa-scout/database"
	"fifa-scout/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	database.Connect()

	app := fiber.New(fiber.Config{
		AppName: "FIFA Wonderkid Scout v2",
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,DELETE,OPTIONS",
		AllowHeaders: "Content-Type",
	}))
	app.Use(logger.New())

	app.Static("/sound", "./sound")

	api := app.Group("/api")
	api.Get("/players", handlers.GetPlayers)
	api.Get("/players/filter", handlers.FilterPlayers)
	api.Get("/predict/:id", handlers.PredictPlayer)
	api.Get("/watchlist", handlers.GetWatchlist)
	api.Post("/watchlist/:id", handlers.AddToWatchlist)
	api.Delete("/watchlist/:id", handlers.RemoveFromWatchlist)

	app.Listen(":8080")
}
