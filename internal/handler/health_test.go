package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/davidborzek/docker-exporter/internal/handler"
	"github.com/stretchr/testify/assert"
)

func TestHealthHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	h := handler.New("")

	h.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}
