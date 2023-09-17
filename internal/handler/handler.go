package handler

import (
	"net/http"
)

type handler struct {
	expectedToken string
	mux           *http.ServeMux
}

func New(authToken string) *handler {
	s := &handler{
		expectedToken: authToken,
		mux:           http.NewServeMux(),
	}

	s.mux.HandleFunc("/metrics", s.handleMetrics())
	s.mux.HandleFunc("/health", s.handleHealth())

	return s
}

func (s *handler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(rw, r)
}
