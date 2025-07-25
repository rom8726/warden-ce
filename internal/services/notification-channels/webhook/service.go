package webhook

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rom8726/warden/internal/domain"
)

type Service struct {
	httpClient *http.Client
	baseURL    string
}

func New(baseURL string) *Service {
	return &Service{
		httpClient: &http.Client{},
		baseURL:    baseURL,
	}
}

func (s *Service) Type() domain.NotificationType {
	return domain.NotificationTypeWebhook
}

func (s *Service) Send(
	ctx context.Context,
	issue *domain.Issue,
	project *domain.Project,
	configData json.RawMessage,
	isRegress bool,
) error {
	var cfg WebhookConfig
	if err := json.Unmarshal(configData, &cfg); err != nil {
		return fmt.Errorf("unmarshal config: %w", err)
	}

	payload := map[string]interface{}{
		"issue_id":          issue.ID.Uint(),
		"project_id":        issue.ProjectID.Uint(),
		"project_name":      project.Name,
		"issue_level":       string(issue.Level),
		"issue_title":       issue.Title,
		"issue_status":      string(issue.Status),
		"issue_first_seen":  issue.FirstSeen,
		"issue_platform":    issue.Platform,
		"issue_occurrences": issue.TotalEvents,
		"issue_url":         fmt.Sprintf("%s/projects/%d/issues/%d", s.baseURL, issue.ProjectID, issue.ID),
		"is_regress":        isRegress,
	}

	reqBody, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, cfg.WebhookURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
