package mattermost

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"text/template"
	"time"

	"github.com/rom8726/warden/internal/domain"
)

type ServiceParams struct {
	BaseURL string
}

type Service struct {
	httpClient *http.Client
	cfg        *ServiceParams
}

func New(cfg *ServiceParams) *Service {
	return &Service{
		httpClient: &http.Client{},
		cfg:        cfg,
	}
}

func (s *Service) Type() domain.NotificationType {
	return domain.NotificationTypeMattermost
}

func (s *Service) Send(
	ctx context.Context,
	issue *domain.Issue,
	project *domain.Project,
	configData json.RawMessage,
	isRegress bool,
) error {
	var cfg MattermostConfig
	if err := json.Unmarshal(configData, &cfg); err != nil {
		return fmt.Errorf("unmarshal config: %w", err)
	}

	message, err := renderMessage(issue, project, s.cfg.BaseURL, isRegress)
	if err != nil {
		return fmt.Errorf("render message: %w", err)
	}

	reqBody, err := json.Marshal(map[string]interface{}{
		"channel": cfg.ChannelName,
		"text":    message,
	})
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

func renderMessage(issue *domain.Issue, project *domain.Project, baseURL string, isRegress bool) (string, error) {
	const msgTemplate = `#### [{{.ProjectName}}] {{.IssueTitle}}
**Is Regress:** {{if .IsRegress}}Yes{{else}}No{{end}}
**Level:** {{.IssueLevel}}
**First Seen:** {{.FirstSeen}}  
**Last Seen:** {{.LastSeen}}  
**URL:** {{.IssueURL}}`

	tmpl, err := template.New("mattermost").Parse(msgTemplate)
	if err != nil {
		return "", fmt.Errorf("parse template: %w", err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, map[string]interface{}{
		"ProjectName": project.Name,
		"IsRegress":   isRegress,
		"IssueTitle":  issue.Title,
		"IssueLevel":  issue.Level,
		//"Environment": issue.Environment,
		"FirstSeen": issue.FirstSeen.Format(time.RFC3339),
		"LastSeen":  issue.LastSeen.Format(time.RFC3339),
		"IssueURL":  fmt.Sprintf("%s/projects/%d/issues/%d", baseURL, issue.ProjectID, issue.ID),
	})
	if err != nil {
		return "", fmt.Errorf("execute template: %w", err)
	}

	return buf.String(), nil
}
