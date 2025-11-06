package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/hoterpoter/telegram-sync/config"
	"github.com/hoterpoter/telegram-sync/handlers"
	"github.com/hoterpoter/telegram-sync/services"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize services
	githubService := services.NewGitHubService(cfg)
	telegramService := services.NewTelegramService(cfg, githubService)

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName: "Telegram Sync v1.0",
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())

	// Routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "Telegram Sync Service is running",
		})
	})

	// Telegram webhook
	app.Post("/webhook/telegram", handlers.NewTelegramHandler(telegramService).HandleWebhook)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Starting server on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
