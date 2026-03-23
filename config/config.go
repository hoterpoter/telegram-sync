package config

import (
	"errors"
	"os"
)

// Config holds all configuration for the application
type Config struct {
	TelegramBotToken   string
	TelegramWebhookURL string
	GitHubToken        string
	GitHubOwner        string
	GitHubRepo         string
	GitHubBranch       string
	DownloadPath       string
}

// Load reads configuration from environment variables
func Load() (*Config, error) {
	cfg := &Config{
		TelegramBotToken:   os.Getenv("TELEGRAM_BOT_TOKEN"),
		TelegramWebhookURL: os.Getenv("TELEGRAM_WEBHOOK_URL"),
		GitHubToken:        os.Getenv("GITHUB_TOKEN"),
		GitHubOwner:        os.Getenv("GITHUB_OWNER"),
		GitHubRepo:         os.Getenv("GITHUB_REPO"),
		GitHubBranch:       os.Getenv("GITHUB_BRANCH"),
		DownloadPath:       os.Getenv("DOWNLOAD_PATH"),
	}

	// Set defaults
	if cfg.GitHubBranch == "" {
		cfg.GitHubBranch = "main"
	}
	if cfg.DownloadPath == "" {
		cfg.DownloadPath = "./downloads"
	}

	// Validate required fields
	if cfg.TelegramBotToken == "" {
		return nil, errors.New("TELEGRAM_BOT_TOKEN is required")
	}
	if cfg.GitHubToken == "" {
		return nil, errors.New("GITHUB_TOKEN is required")
	}
	if cfg.GitHubOwner == "" {
		return nil, errors.New("GITHUB_OWNER is required")
	}
	if cfg.GitHubRepo == "" {
		return nil, errors.New("GITHUB_REPO is required")
	}

	return cfg, nil
}
