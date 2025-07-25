package infra

import (
	"context"
	"log/slog"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"

	"github.com/rom8726/warden/pkg/resilience"
)

type ClickHouseConn interface {
	clickhouse.Conn

	QueryWithRetries(ctx context.Context, query string, args ...any) (driver.Rows, error)
	ExecWithRetries(ctx context.Context, query string, args ...any) error
}

var _ ClickHouseConn = (*ClickHouseConnImpl)(nil)

// ClickHouseConnImpl is a wrapper around clickhouse.Conn that adds circuit breaker and retry functionality.
type ClickHouseConnImpl struct {
	clickhouse.Conn
	circuitBreaker resilience.CircuitBreaker
}

// NewClickHouseConn creates a new ClickHouseConn with circuit breaker.
func NewClickHouseConn(conn clickhouse.Conn) *ClickHouseConnImpl {
	return &ClickHouseConnImpl{
		Conn:           conn,
		circuitBreaker: resilience.NewClickHouseCircuitBreaker(),
	}
}

// QueryWithRetries executes a query with circuit breaker and retry patterns.
func (c *ClickHouseConnImpl) QueryWithRetries(ctx context.Context, query string, args ...any) (driver.Rows, error) {
	var rows driver.Rows

	err := resilience.WithCircuitBreakerAndRetry(
		ctx,
		c.circuitBreaker,
		func(ctx context.Context) error {
			var err error
			rows, err = c.Conn.Query(ctx, query, args...)
			if err != nil {
				slog.Debug("ClickHouse query failed, will retry", "error", err)

				return err
			}

			return nil
		},
		resilience.ClickHouseRetryOptions()...,
	)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

// ExecWithRetries executes a statement with circuit breaker and retry patterns.
func (c *ClickHouseConnImpl) ExecWithRetries(ctx context.Context, query string, args ...any) error {
	return resilience.WithCircuitBreakerAndRetry(
		ctx,
		c.circuitBreaker,
		func(ctx context.Context) error {
			err := c.Conn.Exec(ctx, query, args...)
			if err != nil {
				slog.Debug("ClickHouse exec failed, will retry", "error", err)

				return err
			}

			return nil
		},
		resilience.ClickHouseRetryOptions()...,
	)
}
