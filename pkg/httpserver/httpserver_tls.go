package httpserver

import (
	"context"
	"net"
	"net/http"
	"time"
)

type ServerTLS struct {
	Listener net.Listener

	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration

	Handler http.Handler

	CertFile string
	KeyFile  string
}

func (s *ServerTLS) ListenAndServe(ctx context.Context) error {
	srv := http.Server{
		ReadTimeout:  s.ReadTimeout,
		WriteTimeout: s.WriteTimeout,
		IdleTimeout:  s.IdleTimeout,
		Handler:      s.Handler,
	}

	errc := make(chan error, 1)

	go func() {
		errc <- srv.ServeTLS(s.Listener, s.CertFile, s.KeyFile)
		close(errc)
	}()

	select {
	case <-ctx.Done():
	case err := <-errc:
		return err
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return srv.Shutdown(shutdownCtx)
}
