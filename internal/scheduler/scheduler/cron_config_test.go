package scheduler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCronConfigDailyAlerts_Schedule(t *testing.T) {
	config := &CronConfigDailyAlerts{}

	schedule := config.Schedule()

	assert.Equal(t, "0 0 11 * * *", schedule)
}

func TestCronConfigSummaryAlerts_Schedule(t *testing.T) {
	config := &CronConfigSummaryAlerts{}

	schedule := config.Schedule()

	assert.Equal(t, "0 0 14 * * 1", schedule)
}

func TestConfigInterface(t *testing.T) {
	t.Run("daily alerts implements Config interface", func(t *testing.T) {
		var config Config = &CronConfigDailyAlerts{}

		schedule := config.Schedule()
		assert.Equal(t, "0 0 11 * * *", schedule)
	})

	t.Run("summary alerts implements Config interface", func(t *testing.T) {
		var config Config = &CronConfigSummaryAlerts{}

		schedule := config.Schedule()
		assert.Equal(t, "0 0 14 * * 1", schedule)
	})
}
