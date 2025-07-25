package mattermost

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMattermostConfig(t *testing.T) {
	t.Run("create mattermost config", func(t *testing.T) {
		config := MattermostConfig{
			WebhookURL:  "https://mattermost.example.com/hooks/test",
			ChannelName: "general",
		}

		assert.Equal(t, "https://mattermost.example.com/hooks/test", config.WebhookURL)
		assert.Equal(t, "general", config.ChannelName)
	})

	t.Run("marshal mattermost config to json", func(t *testing.T) {
		config := MattermostConfig{
			WebhookURL:  "https://mattermost.example.com/hooks/test",
			ChannelName: "general",
		}

		data, err := json.Marshal(config)
		assert.NoError(t, err)
		assert.Equal(t, `{"webhook_url":"https://mattermost.example.com/hooks/test","channel_name":"general"}`, string(data))
	})

	t.Run("unmarshal mattermost config from json", func(t *testing.T) {
		jsonData := `{"webhook_url":"https://mattermost.example.com/hooks/test","channel_name":"general"}`
		var config MattermostConfig

		err := json.Unmarshal([]byte(jsonData), &config)
		assert.NoError(t, err)
		assert.Equal(t, "https://mattermost.example.com/hooks/test", config.WebhookURL)
		assert.Equal(t, "general", config.ChannelName)
	})

	t.Run("empty mattermost config", func(t *testing.T) {
		config := MattermostConfig{}

		assert.Equal(t, "", config.WebhookURL)
		assert.Equal(t, "", config.ChannelName)
	})
}
