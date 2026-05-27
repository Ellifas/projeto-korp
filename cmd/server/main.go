package main

import (
	"log"
	"net/http"
	"time"

	"github.com/ellifasr/projeto-korp/internal/handler"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	registry := prometheus.NewRegistry()
	projetoKorpHandler := handler.NewProjetoKorpHandler(registry)

	mux := http.NewServeMux()
	mux.HandleFunc("/projeto-korp", projetoKorpHandler.ServeProjetoKorp)
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if _, err := w.Write([]byte(`{"status":"ok"}`)); err != nil {
			log.Printf("failed to write health response: %v", err)
		}
	})
	mux.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))

	server := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Println("http-server-projeto-korp listening on port 8080")

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
