package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTxFromContext(t *testing.T) {
	t.Run("with transaction in context", func(t *testing.T) {
		ctx := context.Background()
		mockTx := &mockTx{}
		txCtx := context.WithValue(ctx, txKey{}, mockTx)

		result := TxFromContext(txCtx)
		assert.Equal(t, mockTx, result)
	})

	t.Run("without transaction in context", func(t *testing.T) {
		ctx := context.Background()

		result := TxFromContext(ctx)
		assert.Nil(t, result)
	})

	t.Run("with wrong type in context", func(t *testing.T) {
		ctx := context.Background()
		wrongValue := "not a transaction"
		txCtx := context.WithValue(ctx, txKey{}, wrongValue)

		result := TxFromContext(txCtx)
		assert.Nil(t, result)
	})
}
