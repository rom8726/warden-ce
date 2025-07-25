package telegram

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTelegramConfig(t *testing.T) {
	t.Run("create telegram config", func(t *testing.T) {
		config := TelegramConfig{
			BotToken: "1234567890:ABCdefGHIjklMNOpqrsTUVwxyz",
			ChatID:   "123456789",
		}

		assert.Equal(t, "1234567890:ABCdefGHIjklMNOpqrsTUVwxyz", config.BotToken)
		assert.Equal(t, "123456789", config.ChatID)
	})

	t.Run("marshal telegram config to json", func(t *testing.T) {
		config := TelegramConfig{
			BotToken: "1234567890:ABCdefGHIjklMNOpqrsTUVwxyz",
			ChatID:   "123456789",
		}

		data, err := json.Marshal(config)
		assert.NoError(t, err)
		assert.Equal(t, `{"bot_token":"1234567890:ABCdefGHIjklMNOpqrsTUVwxyz","chat_id":"123456789"}`, string(data))
	})

	t.Run("unmarshal telegram config from json", func(t *testing.T) {
		jsonData := `{"bot_token":"1234567890:ABCdefGHIjklMNOpqrsTUVwxyz","chat_id":"123456789"}`
		var config TelegramConfig

		err := json.Unmarshal([]byte(jsonData), &config)
		assert.NoError(t, err)
		assert.Equal(t, "1234567890:ABCdefGHIjklMNOpqrsTUVwxyz", config.BotToken)
		assert.Equal(t, "123456789", config.ChatID)
	})

	t.Run("empty telegram config", func(t *testing.T) {
		config := TelegramConfig{}

		assert.Equal(t, "", config.BotToken)
		assert.Equal(t, "", config.ChatID)
	})
}
