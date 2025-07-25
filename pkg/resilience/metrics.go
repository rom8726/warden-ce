package resilience

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// CircuitBreakerState tracks the current state of circuit breakers.
	CircuitBreakerState = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "warden_circuit_breaker_state",
			Help: "Current state of the circuit breaker (0 - Closed, 1 - Open, 2 - HalfOpen)",
		},
		[]string{"name"},
	)

	// CircuitBreakerFailures tracks the total number of failures.
	CircuitBreakerFailures = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "warden_circuit_breaker_failures_total",
			Help: "Total number of circuit breaker failures",
		},
		[]string{"name"},
	)

	// CircuitBreakerSuccesses tracks the total number of successes.
	CircuitBreakerSuccesses = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "warden_circuit_breaker_successes_total",
			Help: "Total number of circuit breaker successes",
		},
		[]string{"name"},
	)

	// CircuitBreakerTimeouts tracks the total number of timeouts.
	CircuitBreakerTimeouts = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "warden_circuit_breaker_timeouts_total",
			Help: "Total number of circuit breaker timeouts",
		},
		[]string{"name"},
	)

	// RetryAttempts tracks the total number of retry attempts.
	RetryAttempts = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "warden_retry_attempts_total",
			Help: "Total number of retry attempts",
		},
		[]string{"operation"},
	)
)

// UpdateCircuitBreakerMetrics updates only the state metric (Gauge) for a circuit breaker.
// Счетчики (success, failure, timeout) обновляются только в момент события.
func UpdateCircuitBreakerMetrics(cb CircuitBreaker) {
	name := cb.Name()
	state := cb.State()

	// Update state metric (0 - Closed, 1 - Open, 2 - HalfOpen)
	CircuitBreakerState.WithLabelValues(name).Set(float64(state))
}

// IncrementRetryAttempts increments the retry attempts counter.
func IncrementRetryAttempts(operation string) {
	RetryAttempts.WithLabelValues(operation).Inc()
}
