package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type ProjetoKorpHandler struct {
	requestsTotal   *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
	serviceUp       prometheus.Gauge
}

type projetoKorpResponse struct {
	Nome    string `json:"nome"`
	Horario string `json:"horario"`
}

func NewProjetoKorpHandler(reg prometheus.Registerer) *ProjetoKorpHandler {
	requestsTotal := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "korp_requests_total",
			Help: "Total de requisições recebidas pelo serviço http-server-projeto-korp.",
		},
		[]string{"endpoint"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "korp_request_duration_seconds",
			Help:    "Duração das requisições HTTP em segundos.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"endpoint"},
	)

	serviceUp := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "korp_service_up",
			Help: "Indica se o serviço http-server-projeto-korp está disponível. 1 = disponível.",
		},
	)

	reg.MustRegister(requestsTotal, requestDuration, serviceUp)

	serviceUp.Set(1)

	return &ProjetoKorpHandler{
		requestsTotal:   requestsTotal,
		requestDuration: requestDuration,
		serviceUp:       serviceUp,
	}
}

func (h *ProjetoKorpHandler) ServeProjetoKorp(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	endpoint := "/projeto-korp"

	defer func() {
		h.requestsTotal.WithLabelValues(endpoint).Inc()
		h.requestDuration.WithLabelValues(endpoint).Observe(time.Since(start).Seconds())
	}()

	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response := projetoKorpResponse{
		Nome:    "Projeto Korp",
		Horario: time.Now().UTC().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}
