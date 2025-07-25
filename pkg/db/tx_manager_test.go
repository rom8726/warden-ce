package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/mock"
)

//func TestNewTxManager(t *testing.T) {
//	mockPool := &mockPgxPool{}
//
//	txManager := NewTxManager(mockPool)
//
//	assert.NotNil(t, txManager)
//	assert.Equal(t, mockPool, txManager.pool)
//}

//func TestTxManagerImpl_ReadCommitted(t *testing.T) {
//	t.Run("successful transaction", func(t *testing.T) {
//		mockPool := &mockPgxPool{}
//		mockTx := &mockTx{}
//
//		mockPool.On("BeginTx", mock.Anything, mock.Anything).Return(mockTx, nil)
//		mockTx.On("Commit", mock.Anything).Return(nil)
//
//		txManager := NewTxManager(mockPool)
//
//		err := txManager.ReadCommitted(context.Background(), func(ctx context.Context) error {
//			return nil
//		})
//
//		assert.NoError(t, err)
//		mockPool.AssertExpectations(t)
//		mockTx.AssertExpectations(t)
//	})
//
//	t.Run("transaction with error", func(t *testing.T) {
//		mockPool := &mockPgxPool{}
//		mockTx := &mockTx{}
//		expectedErr := errors.New("transaction error")
//
//		mockPool.On("BeginTx", mock.Anything, mock.Anything).Return(mockTx, nil)
//		mockTx.On("Rollback", mock.Anything).Return(nil)
//
//		txManager := NewTxManager(mockPool)
//
//		err := txManager.ReadCommitted(context.Background(), func(ctx context.Context) error {
//			return expectedErr
//		})
//
//		assert.Equal(t, expectedErr, err)
//		mockPool.AssertExpectations(t)
//		mockTx.AssertExpectations(t)
//	})
//
//	t.Run("begin transaction error", func(t *testing.T) {
//		mockPool := &mockPgxPool{}
//		expectedErr := errors.New("begin error")
//
//		mockPool.On("BeginTx", mock.Anything, mock.Anything).Return(nil, expectedErr)
//
//		txManager := NewTxManager(mockPool)
//
//		err := txManager.ReadCommitted(context.Background(), func(ctx context.Context) error {
//			return nil
//		})
//
//		assert.Equal(t, expectedErr, err)
//		mockPool.AssertExpectations(t)
//	})
//}
//
//func TestTxManagerImpl_RepeatableRead(t *testing.T) {
//	t.Run("successful transaction", func(t *testing.T) {
//		mockPool := &mockPgxPool{}
//		mockTx := &mockTx{}
//
//		mockPool.On("BeginTx", mock.Anything, mock.Anything).Return(mockTx, nil)
//		mockTx.On("Commit", mock.Anything).Return(nil)
//
//		txManager := NewTxManager(mockPool)
//
//		err := txManager.RepeatableRead(context.Background(), func(ctx context.Context) error {
//			return nil
//		})
//
//		assert.NoError(t, err)
//		mockPool.AssertExpectations(t)
//		mockTx.AssertExpectations(t)
//	})
//}

// Mock implementations
//type mockPgxPool struct {
//	mock.Mock
//}
//
//func (m *mockPgxPool) BeginTx(ctx context.Context, opts pgx.TxOptions) (Tx, error) {
//	args := m.Called(ctx, opts)
//	return args.Get(0).(Tx), args.Error(1)
//}

type mockTx struct {
	mock.Mock
}

func (m *mockTx) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	mockArgs := m.Called(ctx, sql, args)
	return mockArgs.Get(0).(pgx.Rows), mockArgs.Error(1)
}

func (m *mockTx) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	mockArgs := m.Called(ctx, sql, args)
	return mockArgs.Get(0).(pgx.Row)
}

func (m *mockTx) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	mockArgs := m.Called(ctx, sql, args)
	return mockArgs.Get(0).(pgconn.CommandTag), mockArgs.Error(1)
}

func (m *mockTx) Commit(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *mockTx) Rollback(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}
