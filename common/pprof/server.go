package pprof

import (
	"context"
	"net/http"
	"net/http/pprof"
	"time"

	"github.com/sirupsen/logrus"
)

type Server struct {
	httpServer *http.Server
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) ListenAndServe(addr string) error {
	mux := http.NewServeMux()
	mux.Handle("/profile/debug/pprof/", http.StripPrefix("/profile", http.HandlerFunc(pprof.Index)))
	mux.Handle("/profile/debug/pprof/cmdline", http.StripPrefix("/profile", http.HandlerFunc(pprof.Cmdline)))
	mux.Handle("/profile/debug/pprof/profile", http.StripPrefix("/profile", http.HandlerFunc(pprof.Profile)))
	mux.Handle("/profile/debug/pprof/symbol", http.StripPrefix("/profile", http.HandlerFunc(pprof.Symbol)))
	mux.Handle("/profile/debug/pprof/trace", http.StripPrefix("/profile", http.HandlerFunc(pprof.Trace)))
	s.httpServer = &http.Server{
		Addr:    addr,
		Handler: mux,
	}
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown() {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancelFunc()
	}()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		logrus.WithError(err).Error("shutdown web control server failed")
	}
}
