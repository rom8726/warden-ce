package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// EventsReceived counts the number of events received.
	EventsReceived = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "warden_events_received_total",
			Help: "The total number of events received",
		},
		[]string{"project_id"},
	)

	// EventsProcessed counts the number of events processed successfully.
	EventsProcessed = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "warden_events_processed_total",
			Help: "The total number of events processed successfully",
		},
		[]string{"project_id"},
	)

	// ExceptionsReceived counts the number of exceptions received.
	ExceptionsReceived = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "warden_exceptions_received_total",
			Help: "The total number of exceptions received",
		},
		[]string{"project_id"},
	)

	// ExceptionsProcessed counts the number of exceptions processed successfully.
	ExceptionsProcessed = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "warden_exceptions_processed_total",
			Help: "The total number of exceptions processed successfully",
		},
		[]string{"project_id"},
	)

	// ValidationErrors counts the number of validation errors.
	ValidationErrors = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "warden_validation_errors_total",
			Help: "The total number of validation errors",
		},
		[]string{"type"},
	)

	// ProcessingTime measures the time taken to process events and exceptions.
	ProcessingTime = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "warden_processing_time_seconds",
			Help:    "The time taken to process events and exceptions",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"type"},
	)

	// UpsertIssueTotal counts the number of UpsertIssue operations.
	UpsertIssueTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "warden_upsert_issue_total",
			Help: "The total number of UpsertIssue operations",
		},
		[]string{"project_id", "status"},
	)

	// UpsertIssueDuration measures the time taken to perform UpsertIssue operations.
	UpsertIssueDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "warden_upsert_issue_duration_seconds",
			Help:    "The time taken to perform UpsertIssue operations",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"project_id"},
	)

	// UpsertIssueErrors counts the number of UpsertIssue operation errors.
	UpsertIssueErrors = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "warden_upsert_issue_errors_total",
			Help: "The total number of UpsertIssue operation errors",
		},
		[]string{"project_id"},
	)

	// KafkaMessagesProduced counts the number of messages produced to Kafka.
	KafkaMessagesProduced = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "warden_kafka_messages_produced_total",
			Help: "The total number of messages produced to Kafka",
		},
		[]string{"topic"},
	)

	// KafkaMessagesConsumed counts the number of messages consumed from Kafka.
	KafkaMessagesConsumed = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "warden_kafka_messages_consumed_total",
			Help: "The total number of messages consumed from Kafka",
		},
		[]string{"topic"},
	)

	// CacheHits counts the number of cache hits.
	CacheHits = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "warden_cache_hits_total",
			Help: "The total number of cache hits",
		},
		[]string{"type"},
	)

	// CacheMisses counts the number of cache misses.
	CacheMisses = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "warden_cache_misses_total",
			Help: "The total number of cache misses",
		},
		[]string{"type"},
	)

	// CacheSize shows the current size of each cache.
	CacheSize = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "warden_cache_size",
			Help: "The current size of each cache",
		},
		[]string{"type"},
	)

	// CacheCapacity shows the maximum capacity of each cache.
	CacheCapacity = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "warden_cache_capacity",
			Help: "The maximum capacity of each cache",
		},
		[]string{"type"},
	)

	// EnvelopeMessagesSent counts the number of envelope messages sent to Kafka.
	EnvelopeMessagesSent = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "warden_envelope_messages_sent_total",
			Help: "The total number of envelope messages sent to Kafka",
		},
		[]string{"priority"},
	)

	// EnvelopeMessagesReceived counts the number of envelope messages received from Kafka.
	EnvelopeMessagesReceived = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "warden_envelope_messages_received_total",
			Help: "The total number of envelope messages received from Kafka",
		},
		[]string{"priority"},
	)

	// EnvelopeProcessingDuration measures the time taken to process envelope messages.
	EnvelopeProcessingDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "warden_envelope_processing_duration_seconds",
			Help:    "The time taken to process envelope messages",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation"},
	)

	// EnvelopeProcessingErrors counts the number of envelope processing errors.
	EnvelopeProcessingErrors = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "warden_envelope_processing_errors_total",
			Help: "The total number of envelope processing errors",
		},
		[]string{"error_type"},
	)

	// EnvelopeQueueSize shows the current size of envelope processing queues.
	EnvelopeQueueSize = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "warden_envelope_queue_size",
			Help: "The current size of envelope processing queues",
		},
		[]string{"priority"},
	)

	// EnvelopeRetryCount counts the number of envelope retries.
	EnvelopeRetryCount = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "warden_envelope_retry_total",
			Help: "The total number of envelope retries",
		},
		[]string{"retry_count"},
	)

	// StoreEventMessagesSent counts the number of store event messages sent to Kafka.
	StoreEventMessagesSent = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "warden_store_event_messages_sent_total",
			Help: "The total number of store event messages sent to Kafka",
		},
		[]string{"priority"},
	)

	// StoreEventMessagesReceived counts the number of store event messages received from Kafka.
	StoreEventMessagesReceived = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "warden_store_event_messages_received_total",
			Help: "The total number of store event messages received from Kafka",
		},
		[]string{"priority"},
	)

	// StoreEventProcessingDuration measures the time taken to process store event messages.
	StoreEventProcessingDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "warden_store_event_processing_duration_seconds",
			Help:    "The time taken to process store event messages",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation"},
	)

	// StoreEventProcessingErrors counts the number of store event processing errors.
	StoreEventProcessingErrors = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "warden_store_event_processing_errors_total",
			Help: "The total number of store event processing errors",
		},
		[]string{"error_type"},
	)
)
