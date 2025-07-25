package versions

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/rom8726/warden/internal/backend/contract"
)

type Service struct {
	componentURLs map[string]string
	httpClient    *http.Client
}

// New creates a new versions service.
func New() *Service {
	return &Service{
		componentURLs: map[string]string{
			"backend":           versionEndpoint(backendHost()),
			"envelope-consumer": versionEndpoint(envelopeConsumerHost()),
			"ingest-server":     versionEndpoint(ingestServerHost()),
			"issue-notificator": versionEndpoint(issueNotificatorHost()),
			"user-notificator":  versionEndpoint(userNotificatorHost()),
			"scheduler":         versionEndpoint(schedulerHost()),
		},
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// GetVersions collects versions from all system components.
func (s *Service) GetVersions(ctx context.Context) ([]contract.ComponentVersion, error) {
	components := make([]contract.ComponentVersion, 0, len(s.componentURLs))

	for name, componentURL := range s.componentURLs {
		version, err := s.getComponentVersion(ctx, componentURL)
		if err != nil {
			// If we can't get a version, mark as unavailable but continue
			components = append(components, contract.ComponentVersion{
				Name:      name,
				Version:   "unknown",
				BuildTime: "unknown",
				Status:    "unavailable",
			})

			continue
		}

		components = append(components, contract.ComponentVersion{
			Name:      name,
			Version:   version.Version,
			BuildTime: version.BuildTime,
			Status:    "available",
		})
	}

	return components, nil
}

type versionResponse struct {
	Version   string `json:"version"`
	BuildTime string `json:"build_time"`
}

func (s *Service) getComponentVersion(ctx context.Context, componentURL string) (*versionResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, componentURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var version versionResponse
	if err := json.NewDecoder(resp.Body).Decode(&version); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return &version, nil
}

func backendHost() string {
	if host := os.Getenv("WARDEN_BACKEND_HOST"); host != "" {
		return host
	}

	return "warden-backend"
}

func envelopeConsumerHost() string {
	if host := os.Getenv("WARDEN_ENVELOPE_CONSUMER_HOST"); host != "" {
		return host
	}

	return "warden-envelope-consumer"
}

func ingestServerHost() string {
	if host := os.Getenv("WARDEN_INGEST_SERVER_HOST"); host != "" {
		return host
	}

	return "warden-ingest-server"
}

func issueNotificatorHost() string {
	if host := os.Getenv("WARDEN_ISSUE_NOTIFICATOR_HOST"); host != "" {
		return host
	}

	return "warden-issue-notificator"
}

func userNotificatorHost() string {
	if host := os.Getenv("WARDEN_USER_NOTIFICATOR_HOST"); host != "" {
		return host
	}

	return "warden-user-notificator"
}

func schedulerHost() string {
	if host := os.Getenv("WARDEN_SCHEDULER_HOST"); host != "" {
		return host
	}

	return "warden-scheduler"
}

func versionEndpoint(host string) string {
	return fmt.Sprintf("http://%s:8081/version", host)
}
