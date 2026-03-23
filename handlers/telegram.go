package handlers

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gofiber/fiber/v2"
	"github.com/hoterpoter/telegram-sync/services"
)

// TelegramHandler handles Telegram webhook requests
type TelegramHandler struct {
	service *services.TelegramService
}

// NewTelegramHandler creates a new Telegram handler
func NewTelegramHandler(service *services.TelegramService) *TelegramHandler {
	return &TelegramHandler{
		service: service,
	}
}

// HandleWebhook processes incoming Telegram updates
func (h *TelegramHandler) HandleWebhook(c *fiber.Ctx) error {
	var update tgbotapi.Update

	if err := c.BodyParser(&update); err != nil {
		log.Printf("Error parsing update: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Process update asynchronously
	go h.service.ProcessUpdate(&update)

	return c.JSON(fiber.Map{
		"status": "ok",
	})
}
