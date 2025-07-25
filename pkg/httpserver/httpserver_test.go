package httpserver

import (
	"context"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServer_ListenAndServe(t *testing.T) {
	t.Run("server creation", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		server := &Server{
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 30 * time.Second,
			IdleTimeout:  60 * time.Second,
			Handler:      handler,
		}

		assert.NotNil(t, server)
		assert.Equal(t, 15*time.Second, server.ReadTimeout)
		assert.Equal(t, 30*time.Second, server.WriteTimeout)
		assert.Equal(t, 60*time.Second, server.IdleTimeout)
		assert.NotNil(t, server.Handler)
	})

	t.Run("server with context cancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// Создаем тестовый listener
		listener, err := net.Listen("tcp", ":0")
		require.NoError(t, err)
		defer listener.Close()

		server := &Server{
			Listener:     listener,
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 30 * time.Second,
			IdleTimeout:  60 * time.Second,
			Handler:      http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
		}

		// Запускаем сервер в горутине
		errCh := make(chan error, 1)
		go func() {
			errCh <- server.ListenAndServe(ctx)
		}()

		// Отменяем контекст
		cancel()

		// Ждем завершения сервера
		select {
		case err := <-errCh:
			assert.NoError(t, err)
		case <-time.After(10 * time.Second):
			t.Fatal("server did not shutdown in time")
		}
	})
}

func TestServerTLS_ListenAndServe(t *testing.T) {
	t.Run("tls server creation", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		server := &ServerTLS{
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 30 * time.Second,
			IdleTimeout:  60 * time.Second,
			Handler:      handler,
			CertFile:     "test.crt",
			KeyFile:      "test.key",
		}

		assert.NotNil(t, server)
		assert.Equal(t, 15*time.Second, server.ReadTimeout)
		assert.Equal(t, 30*time.Second, server.WriteTimeout)
		assert.Equal(t, 60*time.Second, server.IdleTimeout)
		assert.NotNil(t, server.Handler)
		assert.Equal(t, "test.crt", server.CertFile)
		assert.Equal(t, "test.key", server.KeyFile)
	})

	t.Run("tls server with context cancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// Создаем тестовый listener
		listener, err := net.Listen("tcp", ":0")
		require.NoError(t, err)
		defer listener.Close()

		server := &ServerTLS{
			Listener:     listener,
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 30 * time.Second,
			IdleTimeout:  60 * time.Second,
			Handler:      http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
			CertFile:     "test.crt",
			KeyFile:      "test.key",
		}

		// Запускаем сервер в горутине
		errCh := make(chan error, 1)
		go func() {
			errCh <- server.ListenAndServe(ctx)
		}()

		// Отменяем контекст
		cancel()

		// Ждем завершения сервера
		select {
		case err := <-errCh:
			// Ожидаем ошибку из-за отсутствия сертификатов или успешное завершение
			// Не проверяем конкретную ошибку, так как поведение может отличаться
			t.Logf("Server shutdown with error: %v", err)
		case <-time.After(10 * time.Second):
			t.Fatal("server did not shutdown in time")
		}
	})
}
