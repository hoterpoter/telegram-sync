package services

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/google/go-github/v57/github"
	"github.com/hoterpoter/telegram-sync/config"
	"golang.org/x/oauth2"
)

// GitHubService handles GitHub API operations
type GitHubService struct {
	client *github.Client
	config *config.Config
}

// NewGitHubService creates a new GitHub service
func NewGitHubService(cfg *config.Config) *GitHubService {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: cfg.GitHubToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	return &GitHubService{
		client: client,
		config: cfg,
	}
}

// CommitFile commits a file to the GitHub repository
func (s *GitHubService) CommitFile(ctx context.Context, path string, content []byte, message string) error {
	opts := &github.RepositoryContentFileOptions{
		Message: github.String(message),
		Content: content,
		Branch:  github.String(s.config.GitHubBranch),
	}

	// Try to get existing file to update it
	fileContent, _, resp, err := s.client.Repositories.GetContents(
		ctx,
		s.config.GitHubOwner,
		s.config.GitHubRepo,
		path,
		&github.RepositoryContentGetOptions{Ref: s.config.GitHubBranch},
	)

	if err == nil && resp.StatusCode == 200 {
		// File exists, update it
		opts.SHA = fileContent.SHA
	}

	// Create or update the file
	_, _, err = s.client.Repositories.CreateFile(
		ctx,
		s.config.GitHubOwner,
		s.config.GitHubRepo,
		path,
		opts,
	)

	if err != nil {
		return fmt.Errorf("failed to commit file: %w", err)
	}

	log.Printf("Successfully committed file: %s", path)
	return nil
}

// SaveMessage saves a text message to GitHub
func (s *GitHubService) SaveMessage(ctx context.Context, chatID int64, messageID int, username, text string, date time.Time) error {
	// Create a structured path
	dateStr := date.Format("2006-01-02")
	dirPath := filepath.Join("messages", fmt.Sprintf("chat_%d", chatID), dateStr)
	filename := fmt.Sprintf("msg_%d.txt", messageID)
	fullPath := filepath.Join(dirPath, filename)

	// Format message content
	content := fmt.Sprintf("From: %s\nDate: %s\nMessage ID: %d\nChat ID: %d\n\n%s\n",
		username, date.Format(time.RFC3339), messageID, chatID, text)

	commitMsg := fmt.Sprintf("Add message from %s at %s", username, date.Format(time.RFC3339))

	return s.CommitFile(ctx, fullPath, []byte(content), commitMsg)
}

// SavePhoto saves a photo to GitHub
func (s *GitHubService) SavePhoto(ctx context.Context, chatID int64, messageID int, username string, photoData []byte, date time.Time, caption string) error {
	// Create a structured path
	dateStr := date.Format("2006-01-02")
	dirPath := filepath.Join("photos", fmt.Sprintf("chat_%d", chatID), dateStr)
	filename := fmt.Sprintf("photo_%d.jpg", messageID)
	fullPath := filepath.Join(dirPath, filename)

	commitMsg := fmt.Sprintf("Add photo from %s at %s", username, date.Format(time.RFC3339))

	if err := s.CommitFile(ctx, fullPath, photoData, commitMsg); err != nil {
		return err
	}

	// If there's a caption, save it as well
	if caption != "" {
		captionPath := filepath.Join(dirPath, fmt.Sprintf("photo_%d_caption.txt", messageID))
		captionContent := fmt.Sprintf("From: %s\nDate: %s\nMessage ID: %d\nChat ID: %d\n\n%s\n",
			username, date.Format(time.RFC3339), messageID, chatID, caption)

		captionCommitMsg := fmt.Sprintf("Add photo caption from %s at %s", username, date.Format(time.RFC3339))
		if err := s.CommitFile(ctx, captionPath, []byte(captionContent), captionCommitMsg); err != nil {
			log.Printf("Failed to save caption: %v", err)
		}
	}

	return nil
}
