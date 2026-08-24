// Задание: Интеграционная мини-задача — observability для Token Service
//
// Собери минимальную observability-конфигурацию:
//   - структурированные логи slog (JSON, с service/env/version)
//   - Prometheus-метрики на /metrics
//   - OpenTelemetry-трейсы для сценария "выдача токена"
//
// Требования:
//   - минимум 3 полезные метрики
//   - минимум 2 span на один запрос
//   - лог ошибки содержит trace_id
//
// Зависимости:
//   go get github.com/prometheus/client_golang/prometheus
//   go get github.com/prometheus/client_golang/prometheus/promhttp
//   go get go.opentelemetry.io/otel
//   go get go.opentelemetry.io/otel/sdk/trace
//
// Ожидаемый результат:
//   $ go run main.go
//   {"level":"INFO","msg":"server started","service":"token-service","addr":":8080"}
//
//   $ curl -X POST http://localhost:8080/api/v1/tokens -d '{"user_id":"u-1"}'
//   {"token_id":"tok-...","user_id":"u-1"}
//
//   $ curl http://localhost:8080/metrics
//   token_service_http_requests_total{method="POST",route="/api/v1/tokens",status="201"} 1
//   token_service_tokens_issued_total 1

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// Считает количество HTTP-запросов по method, route и status.
var httpRequestsTotal = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "token_service_http_requests_total",
		Help: "Total number of HTTP requests",
	},
	[]string{"method", "route", "status"},
)

// Измеряет длительность HTTP-запросов по method и route.
var httpRequestDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name: "token_service_http_request_duration_seconds",
		Help: "HTTP request duration in seconds",
	},
	[]string{"method", "route"},
)

// Считает общее количество выданных токенов.
var tokensIssued = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "token_service_tokens_issued_total",
		Help: "Total number of issued tokens",
	},
)

// Обёртка над ResponseWriter, которая запоминает HTTP-статус.
type statusRecorder struct {
	http.ResponseWriter
	status int
}

// Перехватываем статус, сохраняем его и передаём настоящему ResponseWriter.
func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func writeJSON(
	w http.ResponseWriter,
	status int,
	payload any,
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(payload)
}

// Собирает HTTP-метрики для каждого запроса.
func metricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// Запоминаем время начала запроса.
			start := time.Now()

			// Оборачиваем настоящий ResponseWriter.
			rec := &statusRecorder{
				ResponseWriter: w,
				status:         http.StatusOK,
			}

			// Запускаем настоящий handler.
			next.ServeHTTP(rec, r)

			// Считаем время выполнения запроса.
			duration := time.Since(start).Seconds()

			// Увеличиваем счётчик запросов.
			httpRequestsTotal.
				WithLabelValues(
					r.Method,
					r.URL.Path,
					strconv.Itoa(rec.status),
				).
				Inc()

			// Записываем длительность запроса.
			httpRequestDuration.
				WithLabelValues(
					r.Method,
					r.URL.Path,
				).
				Observe(duration)
		},
	)
}

// Проверка работоспособности сервиса.
func healthHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	writeJSON(
		w,
		http.StatusOK,
		map[string]string{
			"status": "ok",
		},
	)
}

// Выдача токена.
func handleIssueToken(
	w http.ResponseWriter,
	r *http.Request,
	logger *slog.Logger,
) {
	// Сюда декодируем входящий JSON.
	var req struct {
		UserID string `json:"user_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(
			w,
			http.StatusBadRequest,
			map[string]string{
				"error": "invalid request",
			},
		)
		return
	}

	// Получаем tracer.
	tracer := otel.Tracer("token-service")

	// Родительский span всего сценария выдачи токена.
	ctx, span := tracer.Start(
		r.Context(),
		"issue-token",
	)
	defer span.End()

	// Дочерний span проверки пользователя.
	_, validateSpan := tracer.Start(
		ctx,
		"validate-user",
	)

	// Здесь в реальном сервисе была бы проверка пользователя.
	validateSpan.End()

	// Пример ошибки, чтобы error-лог тоже содержал trace_id.
	if req.UserID == "" {
		logger.Error(
			"token issue failed",
			"trace_id",
			span.SpanContext().TraceID().String(),
			"error",
			"user_id is required",
		)

		writeJSON(
			w,
			http.StatusBadRequest,
			map[string]string{
				"error": "user_id is required",
			},
		)
		return
	}

	// Успешно выдаём токен → увеличиваем счётчик.
	tokensIssued.Inc()

	// Создаём простой учебный ID токена.
	tokenID := fmt.Sprintf(
		"tok-%d",
		time.Now().UnixNano(),
	)

	// Лог связываем с trace через trace_id.
	logger.Info(
		"token issued",
		"trace_id",
		span.SpanContext().TraceID().String(),
		"span_id",
		span.SpanContext().SpanID().String(),
		"user_id",
		req.UserID,
	)

	// Возвращаем созданный токен.
	writeJSON(
		w,
		http.StatusCreated,
		map[string]string{
			"token_id": tokenID,
			"user_id":  req.UserID,
		},
	)
}

func main() {
	// Структурированный JSON-логгер.
	logger := slog.New(
		slog.NewJSONHandler(os.Stdout, nil),
	).With(
		"service", "token-service",
		"env", "development",
		"version", "0.1.0",
	)

	ctx := context.Background()

	// Exporter отправляет готовые span'ы в терминал.
	exporter, err := stdouttrace.New(
		stdouttrace.WithPrettyPrint(),
	)
	if err != nil {
		logger.Error(
			"failed to create trace exporter",
			"error",
			err,
		)
		return
	}

	// Создаём TracerProvider и подключаем exporter через batcher.
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
	)

	// Делаем наш provider глобальным для OpenTelemetry.
	otel.SetTracerProvider(tp)

	// Корректное завершение TracerProvider.
	defer tp.Shutdown(ctx)

	// Регистрируем наши Prometheus-метрики.
	prometheus.MustRegister(
		httpRequestsTotal,
		httpRequestDuration,
		tokensIssued,
	)

	mux := http.NewServeMux()

	mux.Handle(
		"GET /health",
		metricsMiddleware(
			http.HandlerFunc(healthHandler),
		),
	)

	// Здесь обёртка нужна потому, что handleIssueToken
	// принимает дополнительный параметр logger.
	mux.Handle(
		"POST /api/v1/tokens",
		metricsMiddleware(
			http.HandlerFunc(
				func(
					w http.ResponseWriter,
					r *http.Request,
				) {
					handleIssueToken(
						w,
						r,
						logger,
					)
				},
			),
		),
	)

	// Prometheus отдаёт зарегистрированные метрики здесь.
	mux.Handle(
		"GET /metrics",
		promhttp.Handler(),
	)

	srv := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      10 * time.Second,
	}

	logger.Info(
		"server started",
		"addr",
		":8080",
	)

	if err := srv.ListenAndServe(); err != nil &&
		err != http.ErrServerClosed {

		logger.Error(
			"server failed",
			"error",
			err,
		)
	}
}
