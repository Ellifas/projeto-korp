package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

func TestServeProjetoKorpReturnsExpectedJSON(t *testing.T) {
	registry := prometheus.NewRegistry()
	projetoKorpHandler := NewProjetoKorpHandler(registry)

	request := httptest.NewRequest(http.MethodGet, "/projeto-korp", nil)
	responseRecorder := httptest.NewRecorder()

	projetoKorpHandler.ServeProjetoKorp(responseRecorder, request)

	if responseRecorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, responseRecorder.Code)
	}

	contentType := responseRecorder.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Fatalf("expected Content-Type application/json, got %s", contentType)
	}

	var response projetoKorpResponse
	if err := json.NewDecoder(responseRecorder.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Nome != "Projeto Korp" {
		t.Fatalf("expected nome Projeto Korp, got %s", response.Nome)
	}

	parsedTime, err := time.Parse(time.RFC3339, response.Horario)
	if err != nil {
		t.Fatalf("expected horario in RFC3339 format, got %s", response.Horario)
	}

	if parsedTime.Location() != time.UTC {
		t.Fatalf("expected horario in UTC, got %s", parsedTime.Location())
	}
}

func TestServeProjetoKorpRejectsInvalidMethod(t *testing.T) {
	registry := prometheus.NewRegistry()
	projetoKorpHandler := NewProjetoKorpHandler(registry)

	request := httptest.NewRequest(http.MethodPost, "/projeto-korp", nil)
	responseRecorder := httptest.NewRecorder()

	projetoKorpHandler.ServeProjetoKorp(responseRecorder, request)

	if responseRecorder.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status %d, got %d", http.StatusMethodNotAllowed, responseRecorder.Code)
	}

	allowHeader := responseRecorder.Header().Get("Allow")
	if allowHeader != http.MethodGet {
		t.Fatalf("expected Allow header GET, got %s", allowHeader)
	}
}
