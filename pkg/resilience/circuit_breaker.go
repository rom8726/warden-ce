package resilience

import (
	"context"
	"time"
)

// State represents the state of a circuit breaker.
type State int

const (
	// Closed is the normal state of the circuit breaker, allowing requests to pass through.
	Closed State = iota
	// Open is the state where the circuit breaker blocks all requests.
	Open
	// HalfOpen is the state where the circuit breaker allows a limited number of requests to pass through.
	HalfOpen
)

// String returns the string representation of the state.
func (s State) String() string {
	switch s {
	case Closed:
		return "Closed"
	case Open:
		return "Open"
	case HalfOpen:
		return "HalfOpen"
	default:
		return "Unknown"
	}
}

// Metrics represents the metrics of a circuit breaker.
type Metrics struct {
	Successes       int64
	Failures        int64
	Timeouts        int64
	FailureRatio    float64
	StateChangeTime time.Time
}

// CircuitBreaker defines the interface for a circuit breaker.
type CircuitBreaker interface {
	// Execute executes the given function if the circuit breaker is closed or half-open
	// Returns an error if the circuit breaker is open or if the function returns an error
	Execute(ctx context.Context, fn func(ctx context.Context) error) error

	// ExecuteWithFallback executes the given function if the circuit breaker is closed or half-open
	// If the circuit breaker is open or if the function returns an error, the fallback function is executed
	ExecuteWithFallback(
		ctx context.Context,
		fn func(ctx context.Context) error,
		fallback func(ctx context.Context, err error) error,
	) error

	// State returns the current state of the circuit breaker
	State() State

	// Name returns the name of the circuit breaker
	Name() string

	// Metrics returns the metrics of the circuit breaker
	Metrics() Metrics
}

// Config represents the configuration for a circuit breaker.
type Config struct {
	// Name is the name of the circuit breaker
	Name string

	// ErrorThreshold is the threshold of errors that will trip the circuit breaker (0.0-1.0)
	ErrorThreshold float64

	// MinRequests is the minimum number of requests needed before the error threshold is calculated
	MinRequests int64

	// OpenTimeout is the time the circuit breaker stays open before transitioning to half-open
	OpenTimeout time.Duration

	// HalfOpenSuccessThreshold is the number of successful requests in half-open state to transition to closed
	HalfOpenSuccessThreshold int64

	// MaxHalfOpenRequests is the maximum number of requests allowed in half-open state
	MaxHalfOpenRequests int64

	// Timeout is the timeout for requests
	Timeout time.Duration

	// IgnoredErrors is a list of errors that are not counted as failures
	IgnoredErrors []error
}
