package email

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmailConfig(t *testing.T) {
	t.Run("create email config", func(t *testing.T) {
		config := EmailConfig{
			EmailTo: "test@example.com",
		}

		assert.Equal(t, "test@example.com", config.EmailTo)
	})

	t.Run("marshal email config to json", func(t *testing.T) {
		config := EmailConfig{
			EmailTo: "test@example.com",
		}

		data, err := json.Marshal(config)
		assert.NoError(t, err)
		assert.Equal(t, `{"email_to":"test@example.com"}`, string(data))
	})

	t.Run("unmarshal email config from json", func(t *testing.T) {
		jsonData := `{"email_to":"test@example.com"}`
		var config EmailConfig

		err := json.Unmarshal([]byte(jsonData), &config)
		assert.NoError(t, err)
		assert.Equal(t, "test@example.com", config.EmailTo)
	})

	t.Run("empty email config", func(t *testing.T) {
		config := EmailConfig{}

		assert.Equal(t, "", config.EmailTo)
	})
}
