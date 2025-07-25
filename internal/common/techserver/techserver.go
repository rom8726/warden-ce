package techserver

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/pprof"

	"github.com/julienschmidt/httprouter"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/rom8726/warden/internal/common/config"
	"github.com/rom8726/warden/internal/version"
	"github.com/rom8726/warden/pkg/httpserver"
)

func NewTechServer(cfg *config.Server) (*httpserver.Server, error) {
	lis, err := net.Listen("tcp", cfg.Addr)
	if err != nil {
		return nil, fmt.Errorf("listen %q: %w", cfg.Addr, err)
	}

	router := httprouter.New()
	router.Handle(http.MethodGet, "/health",
		func(writer http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
			writer.WriteHeader(http.StatusOK)
			_, _ = writer.Write([]byte("ok"))
		},
	)

	router.Handle(http.MethodGet, "/version",
		func(writer http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
			type Version struct {
				Version   string `json:"version"`
				BuildTime string `json:"build_time"`
			}

			ver := Version{
				Version:   version.Version,
				BuildTime: version.BuildTime,
			}

			verData, err := json.Marshal(ver)
			if err != nil {
				writer.WriteHeader(http.StatusInternalServerError)
				_, _ = writer.Write([]byte(err.Error()))

				return
			}

			writer.WriteHeader(http.StatusOK)
			_, _ = writer.Write(verData)
		},
	)

	router.Handler(http.MethodGet, "/metrics", promhttp.Handler())

	router.HandlerFunc(http.MethodGet, "/debug/pprof", pprof.Index)
	router.HandlerFunc(http.MethodGet, "/debug/cmdline", pprof.Cmdline)
	router.HandlerFunc(http.MethodGet, "/debug/profile", pprof.Profile)
	router.HandlerFunc(http.MethodGet, "/debug/symbol", pprof.Symbol)
	router.HandlerFunc(http.MethodGet, "/debug/trace", pprof.Trace)
	router.Handler(http.MethodGet, "/debug/allocs", pprof.Handler("allocs"))
	router.Handler(http.MethodGet, "/debug/block", pprof.Handler("block"))
	router.Handler(http.MethodGet, "/debug/goroutine", pprof.Handler("goroutine"))
	router.Handler(http.MethodGet, "/debug/heap", pprof.Handler("heap"))
	router.Handler(http.MethodGet, "/debug/mutex", pprof.Handler("mutex"))
	router.Handler(http.MethodGet, "/debug/threadcreate", pprof.Handler("threadcreate"))

	return &httpserver.Server{
		Listener:     lis,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
		Handler:      router,
	}, nil
}
