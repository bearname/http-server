package server

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
)

type HttpServer struct {
	Logger log.Logger
}

func (s *HttpServer) StartServer(port string, handler http.Handler) *http.Server {
	srv := &http.Server{Addr: ":" + port, Handler: handler}
	s.Logger.Error(srv.ListenAndServe())
	return srv
}

func (s *HttpServer) WaitForKillSignal() {
	killSignalChan := s.getKillSignalChan()
	killSignal := <-killSignalChan
	switch killSignal {
	case os.Interrupt:
		s.Logger.Info("got SIGINT...")
	case syscall.SIGTERM:
		s.Logger.Info("got SIGTERM...")
	}
}

func (s *HttpServer) getKillSignalChan() chan os.Signal {
	osKillSignalChan := make(chan os.Signal, 1)
	signal.Notify(osKillSignalChan, os.Interrupt, syscall.SIGTERM)

	return osKillSignalChan
}

func HealthCheckHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprint(w, "{\"status\": \"OK\"}")
}

func ReadyCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprintf(w, "{\"host\": \"%v\"}", r.Host)
}

func LogMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(log.Fields{
			"method":     r.Method,
			"url":        r.URL,
			"remoteAddr": r.RemoteAddr,
			"userAgent":  r.UserAgent(),
		}).Info("got a new request")
		h.ServeHTTP(w, r)
	})
}
