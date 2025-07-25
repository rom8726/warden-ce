package domain

const (
	EnvelopeTopicHigh   = "envelope.high"   // High priority
	EnvelopeTopicNormal = "envelope.normal" // Normal priority
	EnvelopeTopicLow    = "envelope.low"    // Low priority

	StoreEventTopicHigh   = "store-event.high"   // High priority
	StoreEventTopicNormal = "store-event.normal" // Normal priority
	StoreEventTopicLow    = "store-event.low"    // Low priority
)

const (
	EventsKafkaTopic = "clickhouse.events"
)
