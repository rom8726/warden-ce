package resilience

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"time"
)

var (
	// ErrCircuitBreakerOpen is returned when the circuit breaker is open.
	ErrCircuitBreakerOpen = errors.New("circuit breaker is open")
	// ErrTooManyRequests is returned when too many requests are made in half-open state.
	ErrTooManyRequests = errors.New("too many requests in half-open state")
)

// circuitBreaker is an implementation of the CircuitBreaker interface.
type circuitBreaker struct {
	name                     string
	state                    State
	mutex                    sync.RWMutex
	errorThreshold           float64
	minRequests              int64
	openTimeout              time.Duration
	halfOpenSuccessThreshold int64
	maxHalfOpenRequests      int64
	timeout                  time.Duration
	ignoredErrors            []error

	successes         int64
	failures          int64
	timeouts          int64
	halfOpenSuccesses int64
	halfOpenRequests  int64
	stateChangeTime   time.Time
}

// NewCircuitBreaker creates a new circuit breaker with the given configuration.
func NewCircuitBreaker(config Config) CircuitBreaker {
	breaker := &circuitBreaker{
		name:                     config.Name,
		state:                    Closed,
		errorThreshold:           config.ErrorThreshold,
		minRequests:              config.MinRequests,
		openTimeout:              config.OpenTimeout,
		halfOpenSuccessThreshold: config.HalfOpenSuccessThreshold,
		maxHalfOpenRequests:      config.MaxHalfOpenRequests,
		timeout:                  config.Timeout,
		ignoredErrors:            config.IgnoredErrors,
		stateChangeTime:          time.Now(),
	}

	// Initialize state metric
	CircuitBreakerState.WithLabelValues(config.Name).Set(float64(Closed))

	return breaker
}

// Execute executes the given function if the circuit breaker is closed or half-open.
//
//nolint:gocyclo,nestif // need refactoring
func (cb *circuitBreaker) Execute(ctx context.Context, fn func(ctx context.Context) error) error {
	// Check if the context error should be ignored
	if err := ctx.Err(); err != nil && cb.isIgnoredError(err) {
		return err
	}

	// Create a timeout context if timeout is set
	if cb.timeout > 0 {
		ctx, _ = context.WithTimeout(ctx, cb.timeout) //nolint:govet // it's ok here
		// defer cancel()
	}

	// Execute the function and get the error
	err := fn(ctx)

	// Check if the context deadline was exceeded
	if errors.Is(ctx.Err(), context.DeadlineExceeded) && err == nil {
		err = context.DeadlineExceeded
	}

	// If the error should be ignored, return it directly without recording
	if err != nil && cb.isIgnoredError(err) {
		return err
	}

	// Now check the circuit breaker state
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	// Check if the error is a timeout
	if errors.Is(err, context.DeadlineExceeded) {
		cb.timeouts++
		// Update timeout metric
		CircuitBreakerTimeouts.WithLabelValues(cb.name).Inc()
	}

	// If the circuit breaker is open, check if the open timeout has elapsed
	if cb.state == Open {
		if time.Since(cb.stateChangeTime) > cb.openTimeout {
			// Transition to a half-open state
			cb.toHalfOpen()
		} else {
			// If the error is nil, return ErrCircuitBreakerOpen
			if err == nil {
				return ErrCircuitBreakerOpen
			}
			// Otherwise, return the original error
			return err
		}
	}

	// If we're in half-open state, check if we've reached the maximum number of requests
	if cb.state == HalfOpen {
		cb.halfOpenRequests++
		if cb.halfOpenRequests > cb.maxHalfOpenRequests {
			return ErrTooManyRequests
		}
	}

	// Update the circuit breaker state based on the result
	if err != nil {
		// Increment the failure count
		cb.failures++
		// Update failure metric
		CircuitBreakerFailures.WithLabelValues(cb.name).Inc()

		if cb.state == HalfOpen {
			cb.toOpen()
		} else if cb.state == Closed {
			// Check if we should trip the circuit breaker
			total := cb.successes + cb.failures
			if total >= cb.minRequests && float64(cb.failures)/float64(total) >= cb.errorThreshold {
				cb.toOpen()
			}
		}
	} else {
		// Increment the success count
		cb.successes++
		// Update success metric
		CircuitBreakerSuccesses.WithLabelValues(cb.name).Inc()

		if cb.state == HalfOpen {
			cb.halfOpenSuccesses++
			if cb.halfOpenSuccesses >= cb.halfOpenSuccessThreshold {
				cb.toClosed()
			}
		}
	}

	return err
}

// ExecuteWithFallback executes the given function with a fallback.
func (cb *circuitBreaker) ExecuteWithFallback(
	ctx context.Context,
	fn func(ctx context.Context) error,
	fallback func(ctx context.Context, err error) error,
) error {
	err := cb.Execute(ctx, fn)
	if err != nil {
		return fallback(ctx, err)
	}

	return nil
}

// State returns the current state of the circuit breaker.
func (cb *circuitBreaker) State() State {
	cb.mutex.RLock()
	defer cb.mutex.RUnlock()

	// Check if we need to transition from Open to HalfOpen
	if cb.state == Open && time.Since(cb.stateChangeTime) > cb.openTimeout {
		// Release read lock and acquire write lock for state transition
		cb.mutex.RUnlock()
		cb.mutex.Lock()

		// Double-check pattern: verify state is still Open and timeout has elapsed
		if cb.state == Open && time.Since(cb.stateChangeTime) > cb.openTimeout {
			cb.toHalfOpen()
		}

		cb.mutex.Unlock()
		cb.mutex.RLock()
	}

	return cb.state
}

// Name returns the name of the circuit breaker.
func (cb *circuitBreaker) Name() string {
	return cb.name
}

// Metrics returns the metrics of the circuit breaker.
func (cb *circuitBreaker) Metrics() Metrics {
	cb.mutex.RLock()
	defer cb.mutex.RUnlock()

	total := cb.successes + cb.failures
	failureRatio := 0.0
	if total > 0 {
		failureRatio = float64(cb.failures) / float64(total)
	}

	return Metrics{
		Successes:       cb.successes,
		Failures:        cb.failures,
		Timeouts:        cb.timeouts,
		FailureRatio:    failureRatio,
		StateChangeTime: cb.stateChangeTime,
	}
}

// toClosed transitions the circuit breaker to the closed state.
func (cb *circuitBreaker) toClosed() {
	prevState := cb.state
	cb.state = Closed
	cb.stateChangeTime = time.Now()
	cb.halfOpenSuccesses = 0
	cb.halfOpenRequests = 0

	// Update state metric
	CircuitBreakerState.WithLabelValues(cb.name).Set(float64(Closed))

	slog.Info("Circuit breaker state changed",
		"name", cb.name,
		"from", prevState,
		"to", cb.state,
	)
}

// toHalfOpen transitions the circuit breaker to the half-open state.
func (cb *circuitBreaker) toHalfOpen() {
	prevState := cb.state
	cb.state = HalfOpen
	cb.stateChangeTime = time.Now()
	cb.halfOpenSuccesses = 0
	cb.halfOpenRequests = 0

	// Update state metric
	CircuitBreakerState.WithLabelValues(cb.name).Set(float64(HalfOpen))

	slog.Info("Circuit breaker state changed",
		"name", cb.name,
		"from", prevState,
		"to", cb.state,
	)
}

// isIgnoredError checks if the given error should be ignored by the circuit breaker.
func (cb *circuitBreaker) isIgnoredError(err error) bool {
	if err == nil {
		return false
	}

	// Check using errors.Is for wrapped errors first (this handles retry.Error and other wrapped errors)
	for _, ignoredErr := range cb.ignoredErrors {
		if errors.Is(err, ignoredErr) {
			slog.Debug("Error ignored by circuit breaker (wrapped match)",
				"name", cb.name,
				"error", err,
				"ignored_error", ignoredErr,
			)

			return true
		}
	}

	return false
}

// toOpen transitions the circuit breaker to the open state.
func (cb *circuitBreaker) toOpen() {
	prevState := cb.state
	cb.state = Open
	cb.stateChangeTime = time.Now()
	cb.halfOpenSuccesses = 0
	cb.halfOpenRequests = 0

	// Update state metric
	CircuitBreakerState.WithLabelValues(cb.name).Set(float64(Open))

	slog.Info("Circuit breaker state changed",
		"name", cb.name,
		"from", prevState,
		"to", cb.state,
	)
}
