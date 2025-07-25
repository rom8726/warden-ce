package resilience

import (
	"context"
	"time"
)

// NewKafkaCircuitBreaker creates a new circuit breaker for Kafka operations.
func NewKafkaCircuitBreaker() CircuitBreaker {
	return NewCircuitBreaker(Config{
		Name:                     "kafka",
		ErrorThreshold:           0.5, // 50% error threshold
		MinRequests:              10,
		OpenTimeout:              30 * time.Second,
		HalfOpenSuccessThreshold: 2,
		MaxHalfOpenRequests:      5,
		Timeout:                  5 * time.Second,
		IgnoredErrors:            []error{context.Canceled, context.DeadlineExceeded},
	})
}

// NewClickHouseCircuitBreaker creates a new circuit breaker for ClickHouse operations.
func NewClickHouseCircuitBreaker() CircuitBreaker {
	return NewCircuitBreaker(Config{
		Name:                     "clickhouse",
		ErrorThreshold:           0.4, // 40% error threshold
		MinRequests:              5,
		OpenTimeout:              20 * time.Second,
		HalfOpenSuccessThreshold: 3,
		MaxHalfOpenRequests:      3,
		Timeout:                  10 * time.Second,
		IgnoredErrors:            []error{context.Canceled, context.DeadlineExceeded},
	})
}

// NewNotificationCircuitBreaker creates a new circuit breaker for notification services.
func NewNotificationCircuitBreaker() CircuitBreaker {
	return NewCircuitBreaker(Config{
		Name:                     "notification",
		ErrorThreshold:           0.3, // 30% error threshold
		MinRequests:              3,
		OpenTimeout:              15 * time.Second,
		HalfOpenSuccessThreshold: 2,
		MaxHalfOpenRequests:      2,
		Timeout:                  3 * time.Second,
		IgnoredErrors:            []error{context.Canceled, context.DeadlineExceeded},
	})
}

// NewDefaultCircuitBreaker creates a new circuit breaker with default settings.
func NewDefaultCircuitBreaker(name string) CircuitBreaker {
	return NewCircuitBreaker(Config{
		Name:                     name,
		ErrorThreshold:           0.5, // 50% error threshold
		MinRequests:              10,
		OpenTimeout:              30 * time.Second,
		HalfOpenSuccessThreshold: 2,
		MaxHalfOpenRequests:      5,
		Timeout:                  5 * time.Second,
		IgnoredErrors:            []error{context.Canceled, context.DeadlineExceeded},
	})
}
