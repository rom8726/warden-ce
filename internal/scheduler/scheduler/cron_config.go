package scheduler

type Config interface {
	Schedule() string
}

type CronConfigDailyAlerts struct{}

func (*CronConfigDailyAlerts) Schedule() string {
	return "0 0 11 * * *"
}

type CronConfigSummaryAlerts struct{}

func (*CronConfigSummaryAlerts) Schedule() string {
	return "0 0 14 * * 1"
}

type CronAnalyticsStats struct{}

func (*CronAnalyticsStats) Schedule() string {
	return "0 0 2 * * *"
}

type CronNotificationsCleaner struct{}

func (*CronNotificationsCleaner) Schedule() string {
	return "0 */10 * * * *"
}

type CronUserNotificationsCleaner struct{}

func (*CronUserNotificationsCleaner) Schedule() string {
	return "0 0 3 * * *" // Every day at 3 AM
}
