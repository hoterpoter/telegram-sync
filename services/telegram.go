package services

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/hoterpoter/telegram-sync/config"
)

// TelegramService handles Telegram bot operations
type TelegramService struct {
	bot           *tgbotapi.BotAPI
	config        *config.Config
	githubService *GitHubService
}

// NewTelegramService creates a new Telegram service
func NewTelegramService(cfg *config.Config, githubService *GitHubService) *TelegramService {
	bot, err := tgbotapi.NewBotAPI(cfg.TelegramBotToken)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	return &TelegramService{
		bot:           bot,
		config:        cfg,
		githubService: githubService,
	}
}

// ProcessUpdate processes a Telegram update
func (s *TelegramService) ProcessUpdate(update *tgbotapi.Update) {
	if update.Message == nil {
		return
	}

	message := update.Message
	chatID := message.Chat.ID
	messageID := message.MessageID
	username := message.From.UserName
	if username == "" {
		username = fmt.Sprintf("%s %s", message.From.FirstName, message.From.LastName)
	}
	date := time.Unix(int64(message.Date), 0)

	ctx := context.Background()

	// Handle text messages
	if message.Text != "" {
		if err := s.githubService.SaveMessage(ctx, chatID, messageID, username, message.Text, date); err != nil {
			log.Printf("Failed to save message: %v", err)
			s.sendErrorResponse(chatID, "Failed to save message to GitHub")
			return
		}
		s.sendSuccessResponse(chatID, "Message saved to GitHub successfully!")
	}

	// Handle photos
	if message.Photo != nil && len(message.Photo) > 0 {
		// Get the largest photo
		photo := message.Photo[len(message.Photo)-1]

		photoData, err := s.downloadFile(photo.FileID)
		if err != nil {
			log.Printf("Failed to download photo: %v", err)
			s.sendErrorResponse(chatID, "Failed to download photo")
			return
		}

		// Telegram photos are typically JPEG
		caption := message.Caption
		fileExt := ".jpg"
		if err := s.githubService.SavePhoto(ctx, chatID, messageID, username, photoData, date, caption, fileExt); err != nil {
			log.Printf("Failed to save photo: %v", err)
			s.sendErrorResponse(chatID, "Failed to save photo to GitHub")
			return
		}
		s.sendSuccessResponse(chatID, "Photo saved to GitHub successfully!")
	}

	// Handle documents (images sent as files)
	if message.Document != nil && isImageMimeType(message.Document.MimeType) {
		photoData, err := s.downloadFile(message.Document.FileID)
		if err != nil {
			log.Printf("Failed to download document: %v", err)
			s.sendErrorResponse(chatID, "Failed to download document")
			return
		}

		caption := message.Caption
		fileExt := getExtensionFromMimeType(message.Document.MimeType)
		if err := s.githubService.SavePhoto(ctx, chatID, messageID, username, photoData, date, caption, fileExt); err != nil {
			log.Printf("Failed to save document: %v", err)
			s.sendErrorResponse(chatID, "Failed to save document to GitHub")
			return
		}
		s.sendSuccessResponse(chatID, "Document saved to GitHub successfully!")
	}
}

// downloadFile downloads a file from Telegram
func (s *TelegramService) downloadFile(fileID string) ([]byte, error) {
	file, err := s.bot.GetFile(tgbotapi.FileConfig{FileID: fileID})
	if err != nil {
		return nil, fmt.Errorf("failed to get file: %w", err)
	}

	fileURL := file.Link(s.config.TelegramBotToken)

	resp, err := http.Get(fileURL)
	if err != nil {
		return nil, fmt.Errorf("failed to download file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %s", resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return data, nil
}

// sendSuccessResponse sends a success message to the user
func (s *TelegramService) sendSuccessResponse(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, "✅ "+text)
	if _, err := s.bot.Send(msg); err != nil {
		log.Printf("Failed to send success message: %v", err)
	}
}

// sendErrorResponse sends an error message to the user
func (s *TelegramService) sendErrorResponse(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, "❌ "+text)
	if _, err := s.bot.Send(msg); err != nil {
		log.Printf("Failed to send error message: %v", err)
	}
}

// isImageMimeType checks if the MIME type is an image
func isImageMimeType(mimeType string) bool {
	imageMimeTypes := []string{
		"image/jpeg",
		"image/jpg",
		"image/png",
		"image/gif",
		"image/webp",
		"image/bmp",
	}

	for _, t := range imageMimeTypes {
		if t == mimeType {
			return true
		}
	}
	return false
}

// getExtensionFromMimeType returns the file extension for a given MIME type
func getExtensionFromMimeType(mimeType string) string {
	mimeToExt := map[string]string{
		"image/jpeg": ".jpg",
		"image/jpg":  ".jpg",
		"image/png":  ".png",
		"image/gif":  ".gif",
		"image/webp": ".webp",
		"image/bmp":  ".bmp",
	}

	if ext, ok := mimeToExt[mimeType]; ok {
		return ext
	}
	// Default to .jpg if MIME type is unknown
	return ".jpg"
}
