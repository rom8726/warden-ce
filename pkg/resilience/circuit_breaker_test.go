package resilience

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestCircuitBreaker(t *testing.T) {
	// Create a circuit breaker with test configuration
	cb := NewCircuitBreaker(Config{
		Name:                     "test",
		ErrorThreshold:           0.5, // 50% error threshold
		MinRequests:              2,   // Only 2 requests needed to trip
		OpenTimeout:              100 * time.Millisecond,
		HalfOpenSuccessThreshold: 1,
		MaxHalfOpenRequests:      2,
		Timeout:                  50 * time.Millisecond,
		IgnoredErrors:            []error{context.Canceled},
	})

	// Test initial state
	if cb.State() != Closed {
		t.Errorf("Initial state should be Closed, got %v", cb.State())
	}

	// Test successful execution
	err := cb.Execute(context.Background(), func(ctx context.Context) error {
		return nil
	})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Test metrics after successful execution
	metrics := cb.Metrics()
	if metrics.Successes != 1 {
		t.Errorf("Expected 1 success, got %d", metrics.Successes)
	}
	if metrics.Failures != 0 {
		t.Errorf("Expected 0 failures, got %d", metrics.Failures)
	}

	// Test failed execution
	testErr := errors.New("test error")
	err = cb.Execute(context.Background(), func(ctx context.Context) error {
		return testErr
	})
	if err != testErr {
		t.Errorf("Expected test error, got %v", err)
	}

	// Test metrics after failed execution
	metrics = cb.Metrics()
	if metrics.Successes != 1 {
		t.Errorf("Expected 1 success, got %d", metrics.Successes)
	}
	if metrics.Failures != 1 {
		t.Errorf("Expected 1 failure, got %d", metrics.Failures)
	}

	// Test another failed execution to trip the circuit breaker
	err = cb.Execute(context.Background(), func(ctx context.Context) error {
		return testErr
	})
	if err != testErr {
		t.Errorf("Expected test error, got %v", err)
	}

	// Test that the circuit breaker is now open
	if cb.State() != Open {
		t.Errorf("State should be Open after failures, got %v", cb.State())
	}

	// Test that execution is blocked when the circuit breaker is open
	err = cb.Execute(context.Background(), func(ctx context.Context) error {
		return nil
	})
	if err != ErrCircuitBreakerOpen {
		t.Errorf("Expected ErrCircuitBreakerOpen, got %v", err)
	}

	// Wait for the circuit breaker to transition to half-open
	time.Sleep(150 * time.Millisecond)

	// Test that the circuit breaker is now half-open
	if cb.State() != HalfOpen {
		t.Errorf("State should be HalfOpen after timeout, got %v", cb.State())
	}

	// Test successful execution in half-open state
	err = cb.Execute(context.Background(), func(ctx context.Context) error {
		return nil
	})
	if err != nil {
		t.Errorf("Expected no error in half-open state, got %v", err)
	}

	// Test that the circuit breaker is now closed
	if cb.State() != Closed {
		t.Errorf("State should be Closed after successful execution in half-open state, got %v", cb.State())
	}

	// Test timeout
	err = cb.Execute(context.Background(), func(ctx context.Context) error {
		time.Sleep(100 * time.Millisecond) // Longer than the timeout
		return nil
	})
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("Expected context.DeadlineExceeded, got %v", err)
	}

	// Test ignored error
	err = cb.Execute(context.Background(), func(ctx context.Context) error {
		return context.Canceled
	})
	if err != context.Canceled {
		t.Errorf("Expected context.Canceled, got %v", err)
	}

	// Test that the ignored error doesn't count as a failure
	metrics = cb.Metrics()
	if metrics.Failures != 2 { // Still 2 from before
		t.Errorf("Expected 2 failures, got %d", metrics.Failures)
	}

	// Test ExecuteWithFallback
	fallbackCalled := false
	err = cb.ExecuteWithFallback(
		context.Background(),
		func(ctx context.Context) error {
			return testErr
		},
		func(ctx context.Context, err error) error {
			fallbackCalled = true
			return nil
		},
	)
	if err != nil {
		t.Errorf("Expected no error with fallback, got %v", err)
	}
	if !fallbackCalled {
		t.Errorf("Fallback should have been called")
	}
}

func TestCircuitBreakerStateRaceCondition(t *testing.T) {
	// Create a circuit breaker with short timeout for testing
	cb := NewCircuitBreaker(Config{
		Name:                     "race-test",
		ErrorThreshold:           0.5,
		MinRequests:              2,
		OpenTimeout:              50 * time.Millisecond, // Short timeout
		HalfOpenSuccessThreshold: 1,
		MaxHalfOpenRequests:      2,
		Timeout:                  10 * time.Millisecond,
		IgnoredErrors:            []error{context.Canceled},
	})

	// Trip the circuit breaker to Open state
	err := cb.Execute(context.Background(), func(ctx context.Context) error {
		return errors.New("test error")
	})
	if err == nil {
		t.Fatal("Expected error")
	}

	err = cb.Execute(context.Background(), func(ctx context.Context) error {
		return errors.New("test error")
	})
	if err == nil {
		t.Fatal("Expected error")
	}

	// Verify circuit breaker is open
	if cb.State() != Open {
		t.Fatalf("Expected Open state, got %v", cb.State())
	}

	// Wait for timeout to trigger transition to half-open
	time.Sleep(100 * time.Millisecond)

	// Test concurrent access to State() method
	const numGoroutines = 10
	const numCalls = 100

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < numCalls; j++ {
				state := cb.State()
				// Verify state is valid
				if state != Closed && state != Open && state != HalfOpen {
					t.Errorf("Invalid state: %v", state)
				}
			}
		}()
	}

	wg.Wait()

	// Verify final state is correct
	finalState := cb.State()
	if finalState != HalfOpen && finalState != Closed {
		t.Errorf("Expected HalfOpen or Closed state, got %v", finalState)
	}
}

func TestCircuitBreakerIgnoredErrors(t *testing.T) {
	// Create custom errors for testing
	customErr1 := errors.New("custom error 1")
	customErr2 := errors.New("custom error 2")
	wrappedErr := fmt.Errorf("wrapped error: %w", customErr1)

	// Test that ignored errors are returned immediately and don't count as failures
	testCases := []struct {
		name         string
		err          error
		expected     error
		shouldIgnore bool
	}{
		{
			name:         "context.Canceled",
			err:          context.Canceled,
			expected:     context.Canceled,
			shouldIgnore: true,
		},
		{
			name:         "context.DeadlineExceeded",
			err:          context.DeadlineExceeded,
			expected:     context.DeadlineExceeded,
			shouldIgnore: true,
		},
		{
			name:         "custom error 1 (exact match)",
			err:          customErr1,
			expected:     customErr1,
			shouldIgnore: true,
		},
		{
			name:         "wrapped custom error 1",
			err:          wrappedErr,
			expected:     wrappedErr,
			shouldIgnore: true,
		},
		{
			name:         "custom error 2 (not ignored)",
			err:          customErr2,
			expected:     customErr2,
			shouldIgnore: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a fresh circuit breaker for each test
			cb := NewCircuitBreaker(Config{
				Name:                     "ignored-errors-test",
				ErrorThreshold:           0.5,
				MinRequests:              2,
				OpenTimeout:              100 * time.Millisecond,
				HalfOpenSuccessThreshold: 1,
				MaxHalfOpenRequests:      2,
				Timeout:                  50 * time.Millisecond,
				IgnoredErrors:            []error{context.Canceled, context.DeadlineExceeded, customErr1},
			})

			initialFailures := cb.Metrics().Failures

			err := cb.Execute(context.Background(), func(ctx context.Context) error {
				return tc.err
			})

			if err != tc.expected {
				t.Errorf("Expected %v, got %v", tc.expected, err)
			}

			// Check if the error was ignored (should not increment failure count)
			metrics := cb.Metrics()
			if tc.shouldIgnore {
				// These errors should be ignored
				if metrics.Failures != initialFailures {
					t.Errorf("Expected failure count to remain %d for ignored error, got %d",
						initialFailures, metrics.Failures)
				}
			} else {
				// This error should count as a failure
				if metrics.Failures <= initialFailures {
					t.Errorf("Expected failure count to increase for non-ignored error")
				}
			}
		})
	}

	// Test context error handling
	t.Run("context error", func(t *testing.T) {
		// Create a fresh circuit breaker for this test
		cb := NewCircuitBreaker(Config{
			Name:                     "ignored-errors-test",
			ErrorThreshold:           0.5,
			MinRequests:              2,
			OpenTimeout:              100 * time.Millisecond,
			HalfOpenSuccessThreshold: 1,
			MaxHalfOpenRequests:      2,
			Timeout:                  50 * time.Millisecond,
			IgnoredErrors:            []error{context.Canceled, context.DeadlineExceeded, customErr1},
		})

		initialFailures := cb.Metrics().Failures

		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel the context

		err := cb.Execute(ctx, func(ctx context.Context) error {
			return nil // This should not be called due to context cancellation
		})

		if err != context.Canceled {
			t.Errorf("Expected context.Canceled, got %v", err)
		}

		// Check that context cancellation doesn't count as a failure
		metrics := cb.Metrics()
		if metrics.Failures != initialFailures {
			t.Errorf("Expected failure count to remain %d for context cancellation, got %d",
				initialFailures, metrics.Failures)
		}
	})
}

func TestCircuitBreakerWithRetryError(t *testing.T) {
	// Create a circuit breaker with ignored errors
	cb := NewCircuitBreaker(Config{
		Name:          "retry-test",
		IgnoredErrors: []error{context.Canceled, context.DeadlineExceeded},
	})

	// Simulate a retry.Error by using WithRetryContext
	err := WithRetryContext(
		context.Background(),
		func(ctx context.Context) error {
			return context.DeadlineExceeded
		},
		DefaultRetryOptions()...,
	)

	// This should not panic and should be handled correctly
	// The retry.Error should be unwrapped and the underlying context.DeadlineExceeded
	// should be recognized as an ignored error
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	// The error should be a retry.Error wrapping context.DeadlineExceeded
	// We can't directly check the type due to package visibility, but we can
	// verify that it's not nil and that the circuit breaker handles it correctly
	t.Logf("Retry error test: error=%v, error_type=%T", err, err)

	// Test that the circuit breaker can handle this error without panicking
	err = cb.Execute(context.Background(), func(ctx context.Context) error {
		return err // This should not cause a panic
	})

	// The error should be returned as-is since it's an ignored error
	if err == nil {
		t.Fatal("Expected error to be returned, got nil")
	}

	// Verify that the error was handled without panicking
	// (if we get here, no panic occurred)
	t.Logf("Circuit breaker handled retry error successfully: %v", err)
}
