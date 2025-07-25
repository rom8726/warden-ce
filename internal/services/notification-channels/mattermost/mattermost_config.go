package mattermost

type MattermostConfig struct {
	WebhookURL  string `json:"webhook_url"`
	ChannelName string `json:"channel_name"`
}
