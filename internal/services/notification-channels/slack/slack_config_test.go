package slack

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSlackConfig(t *testing.T) {
	t.Run("create slack config", func(t *testing.T) {
		config := SlackConfig{
			WebhookURL:  "https://hooks.slack.com/services/test",
			ChannelName: "#general",
		}

		assert.Equal(t, "https://hooks.slack.com/services/test", config.WebhookURL)
		assert.Equal(t, "#general", config.ChannelName)
	})

	t.Run("marshal slack config to json", func(t *testing.T) {
		config := SlackConfig{
			WebhookURL:  "https://hooks.slack.com/services/test",
			ChannelName: "#general",
		}

		data, err := json.Marshal(config)
		assert.NoError(t, err)
		assert.Equal(t, `{"webhook_url":"https://hooks.slack.com/services/test","channel_name":"#general"}`, string(data))
	})

	t.Run("unmarshal slack config from json", func(t *testing.T) {
		jsonData := `{"webhook_url":"https://hooks.slack.com/services/test","channel_name":"#general"}`
		var config SlackConfig

		err := json.Unmarshal([]byte(jsonData), &config)
		assert.NoError(t, err)
		assert.Equal(t, "https://hooks.slack.com/services/test", config.WebhookURL)
		assert.Equal(t, "#general", config.ChannelName)
	})

	t.Run("empty slack config", func(t *testing.T) {
		config := SlackConfig{}

		assert.Equal(t, "", config.WebhookURL)
		assert.Equal(t, "", config.ChannelName)
	})
}
