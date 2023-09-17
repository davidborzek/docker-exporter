package handler_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/davidborzek/docker-exporter/internal/handler"
	"github.com/stretchr/testify/assert"
)

const (
	authToken = "someToken"
)

func TestMetricsHandlerReturnsOK(t *testing.T) {
	req, err := http.NewRequest("GET", "/metrics", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	h := handler.New("")

	h.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestMetricsHandlerReturnsUnauthorizedForEmptyAuthorizationHeader(t *testing.T) {
	req, err := http.NewRequest("GET", "/metrics", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	h := handler.New(authToken)

	h.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestMetricsHandlerReturnsUnauthorizedForInvalidToken(t *testing.T) {
	req, err := http.NewRequest("GET", "/metrics", nil)
	req.Header.Add("Authorization", "Bearer invalidToken")
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	h := handler.New(authToken)

	h.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestMetricsHandlerReturnsOKForValidToken(t *testing.T) {
	req, err := http.NewRequest("GET", "/metrics", nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", authToken))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	h := handler.New(authToken)

	h.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}
