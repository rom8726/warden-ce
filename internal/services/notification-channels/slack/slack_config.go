package slack

type SlackConfig struct {
	WebhookURL  string `json:"webhook_url"`
	ChannelName string `json:"channel_name"`
}
