package main

import (
	"log"
	"os"
	"slices"
	"strconv"
	"time"

	"github.com/AstraBert/what-a-git-year-v2/handlers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/storage/redis/v3"
	"github.com/gofiber/storage/sqlite3"
)

func main() {
	// Create a new Fiber app
	app := Setup()

	// Start the Fiber server on port 8000
	if err := app.Listen(":8000"); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

func cacheSetup() (fiber.Handler, error) {
	redisPort, err := strconv.Atoi(os.Getenv("CACHE_REDIS_PORT"))
	if err != nil {
		return nil, err
	}
	cacheStorage := redis.New(
		redis.Config{
			Host:     os.Getenv("CACHE_REDIS_HOST"),
			Password: os.Getenv("CACHE_REDIS_PASSWORD"),
			Username: os.Getenv("CACHE_REDIS_USER"),
			Port:     redisPort,
		},
	)
	cache := cache.New(
		cache.Config{
			Expiration:   1 * time.Hour,
			CacheControl: true,
			Storage:      cacheStorage,
		},
	)
	return cache, nil
}

func corsSetup(methods string) fiber.Handler {
	allowedOrigins := []string{"https://gityear.re"}
	corsHandler := cors.New(
		cors.Config{
			AllowOriginsFunc: func(origin string) bool {
				return slices.Contains(allowedOrigins, origin)
			},
			AllowMethods: methods,
		},
	)
	return corsHandler
}

func limiterSetup(reqPerMinute int) fiber.Handler {
	limiterStorage := sqlite3.New(
		sqlite3.Config{
			Database:        "ratelimiter",
			Table:           os.Getenv("RATE_LIMITING_TABLE"),
			ConnMaxLifetime: 5 * time.Second,
		},
	)
	limiter := limiter.New(
		limiter.Config{
			Max:     reqPerMinute,
			Storage: limiterStorage,
		},
	)
	return limiter
}

func Setup() *fiber.App {
	app := fiber.New()
	cache, err := cacheSetup()
	if err == nil {
		app.Use(cache)
	}
	app.Get("/signin", corsSetup("GET"), handlers.LoginRoute)
	app.Get("/signup", corsSetup("GET"), handlers.SignUpRoute)
	app.Post("/login", limiterSetup(10), corsSetup("POST"), handlers.HandleLogin)
	app.Post("/logout", limiterSetup(10), corsSetup("POST"), handlers.HandleLogout)
	app.Post("/register", limiterSetup(10), corsSetup("POST"), handlers.HandleSignUp)
	app.Get("/", handlers.HomeRoute)
	app.Get("/search", corsSetup("GET"), handlers.SearchRoute)
	app.Post("/search/gateway", limiterSetup(20), corsSetup("POST"), handlers.HandleSearchGateway)
	app.Static("/static", "./static/")
	app.Use(handlers.PageDoesNotExistRoute)
	return app
}
