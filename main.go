package main

import (
	"log"

	"github.com/AstraBert/what-a-git-year-v2/handlers"
	"github.com/gofiber/fiber/v2"
)

func main() {
	// Create a new Fiber app
	app := Setup()

	// Start the Fiber server on port 8000
	if err := app.Listen(":8000"); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

func Setup() *fiber.App {
	app := fiber.New()
	app.Get("/signin", handlers.LoginRoute)
	app.Get("/signup", handlers.SignUpRoute)
	app.Post("/login", handlers.HandleLogin)
	app.Post("/logout", handlers.HandleLogout)
	app.Post("/register", handlers.HandleSignUp)
	app.Get("/", handlers.HomeRoute)
	app.Get("/search", handlers.SearchRoute)
	app.Post("/search/gateway", handlers.HandleSearchGateway)
	app.Static("/static", "./static/")
	app.Use(handlers.PageDoesNotExistRoute)
	return app
}
