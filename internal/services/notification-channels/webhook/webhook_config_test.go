package webhook

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWebhookConfig(t *testing.T) {
	t.Run("create webhook config", func(t *testing.T) {
		config := WebhookConfig{
			WebhookURL: "https://webhook.example.com/hook",
		}

		assert.Equal(t, "https://webhook.example.com/hook", config.WebhookURL)
	})

	t.Run("marshal webhook config to json", func(t *testing.T) {
		config := WebhookConfig{
			WebhookURL: "https://webhook.example.com/hook",
		}

		data, err := json.Marshal(config)
		assert.NoError(t, err)
		assert.Equal(t, `{"webhook_url":"https://webhook.example.com/hook"}`, string(data))
	})

	t.Run("unmarshal webhook config from json", func(t *testing.T) {
		jsonData := `{"webhook_url":"https://webhook.example.com/hook"}`
		var config WebhookConfig

		err := json.Unmarshal([]byte(jsonData), &config)
		assert.NoError(t, err)
		assert.Equal(t, "https://webhook.example.com/hook", config.WebhookURL)
	})

	t.Run("empty webhook config", func(t *testing.T) {
		config := WebhookConfig{}

		assert.Equal(t, "", config.WebhookURL)
	})
}
