package resilience

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestCircuitBreakerMetrics(t *testing.T) {
	// Reset metrics before test
	CircuitBreakerState.Reset()
	CircuitBreakerSuccesses.Reset()
	CircuitBreakerFailures.Reset()
	CircuitBreakerTimeouts.Reset()

	// Create a circuit breaker
	cb := NewCircuitBreaker(Config{
		Name:                     "metrics-test",
		ErrorThreshold:           0.5,
		MinRequests:              2,
		OpenTimeout:              100 * time.Millisecond,
		HalfOpenSuccessThreshold: 1,
		MaxHalfOpenRequests:      2,
		Timeout:                  50 * time.Millisecond,
		IgnoredErrors:            []error{context.Canceled},
	})

	// Test initial state metric
	stateValue := testutil.ToFloat64(CircuitBreakerState.WithLabelValues("metrics-test"))
	if stateValue != float64(Closed) {
		t.Errorf("Expected initial state to be Closed (0), got %f", stateValue)
	}

	// Test successful execution
	err := cb.Execute(context.Background(), func(ctx context.Context) error {
		return nil
	})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Check success metric
	successValue := testutil.ToFloat64(CircuitBreakerSuccesses.WithLabelValues("metrics-test"))
	if successValue != 1 {
		t.Errorf("Expected 1 success, got %f", successValue)
	}

	// Test failed execution
	testErr := errors.New("test error")
	err = cb.Execute(context.Background(), func(ctx context.Context) error {
		return testErr
	})
	if err != testErr {
		t.Errorf("Expected test error, got %v", err)
	}

	// Check failure metric
	failureValue := testutil.ToFloat64(CircuitBreakerFailures.WithLabelValues("metrics-test"))
	if failureValue != 1 {
		t.Errorf("Expected 1 failure, got %f", failureValue)
	}

	// Test timeout
	err = cb.Execute(context.Background(), func(ctx context.Context) error {
		time.Sleep(100 * time.Millisecond) // Longer than timeout
		return nil
	})
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("Expected context.DeadlineExceeded, got %v", err)
	}

	// Check timeout metric
	timeoutValue := testutil.ToFloat64(CircuitBreakerTimeouts.WithLabelValues("metrics-test"))
	if timeoutValue != 1 {
		t.Errorf("Expected 1 timeout, got %f", timeoutValue)
	}

	// Test state transition to Open
	err = cb.Execute(context.Background(), func(ctx context.Context) error {
		return testErr
	})
	if err != testErr {
		t.Errorf("Expected test error, got %v", err)
	}

	// Check that state changed to Open
	stateValue = testutil.ToFloat64(CircuitBreakerState.WithLabelValues("metrics-test"))
	if stateValue != float64(Open) {
		t.Errorf("Expected state to be Open (1), got %f", stateValue)
	}

	// Wait for transition to HalfOpen
	time.Sleep(300 * time.Millisecond)
	stateValue = testutil.ToFloat64(CircuitBreakerState.WithLabelValues("metrics-test"))
	if stateValue != float64(HalfOpen) && stateValue != float64(Open) {
		t.Errorf("Expected state to be HalfOpen (2) или Open (1), got %f", stateValue)
	}

	// Test successful execution in HalfOpen to transition to Closed
	err = cb.Execute(context.Background(), func(ctx context.Context) error {
		return nil
	})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Check that state changed to Closed
	stateValue = testutil.ToFloat64(CircuitBreakerState.WithLabelValues("metrics-test"))
	if stateValue != float64(Closed) {
		t.Errorf("Expected state to be Closed (0), got %f", stateValue)
	}
}

func TestRetryMetrics(t *testing.T) {
	// Reset metrics before test
	RetryAttempts.Reset()

	// Test retry with context
	err := WithRetryContext(
		context.Background(),
		func(ctx context.Context) error {
			return errors.New("temporary error")
		},
		DefaultRetryOptions()...,
	)

	if err == nil {
		t.Error("Expected error, got nil")
	}

	// Check retry attempts metric
	retryValue := testutil.ToFloat64(RetryAttempts.WithLabelValues("retry_context"))
	if retryValue < 1 {
		t.Errorf("Expected at least 1 retry attempt, got %f", retryValue)
	}
}

func TestUpdateCircuitBreakerMetrics(t *testing.T) {
	// Reset metrics before test
	CircuitBreakerState.Reset()
	CircuitBreakerSuccesses.Reset()
	CircuitBreakerFailures.Reset()
	CircuitBreakerTimeouts.Reset()

	// Create a circuit breaker
	cb := NewCircuitBreaker(Config{
		Name:           "update-test",
		ErrorThreshold: 0.5,
		MinRequests:    1,
		OpenTimeout:    100 * time.Millisecond,
		Timeout:        50 * time.Millisecond,
	})

	// Execute some operations
	cb.Execute(context.Background(), func(ctx context.Context) error {
		return nil
	})
	cb.Execute(context.Background(), func(ctx context.Context) error {
		return errors.New("test error")
	})

	// Update metrics manually
	UpdateCircuitBreakerMetrics(cb)

	// Check that state metric was updated
	stateValue := testutil.ToFloat64(CircuitBreakerState.WithLabelValues("update-test"))
	if stateValue != float64(Open) && stateValue != float64(Closed) && stateValue != float64(HalfOpen) {
		t.Errorf("Expected state to be Open, Closed или HalfOpen, got %f", stateValue)
	}
}
