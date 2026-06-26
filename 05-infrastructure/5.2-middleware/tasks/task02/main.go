// Задание: Logging middleware
//
// Реализуй LoggingMiddleware на базе log/slog, который логирует:
// метод, путь, HTTP-статус ответа и длительность обработки.
//
// Для перехвата статуса используй обёртку statusRecorder.
//
// Ожидаемый результат:
//   $ go run main.go &
//   server started on :8080
//
//   $ curl http://localhost:8080/health
//   ok
//
//   (в логах сервера появится JSON-строка вида)
//   {"time":"...","level":"INFO","msg":"http request","method":"GET","path":"/health","status":200,"duration_ms":0}

package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

type Middleware func(http.Handler) http.Handler

func Chain(h http.Handler, middlewares ...Middleware) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}

// TODO: определи структуру statusRecorder, которая оборачивает http.ResponseWriter
// и запоминает HTTP-статус ответа.
type statusRecorder struct {
	http.ResponseWriter
	status int
}

// TODO: переопредели метод WriteHeader у statusRecorder
func (r *statusRecorder) WriteHeader(status int) {
	r.status = status                    // cохранили себе
	r.ResponseWriter.WriteHeader(status) // отправили статус клиенту
}

// TODO: реализуй LoggingMiddleware
func LoggingMiddleware(logger *slog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			//             1. Запомни время начала: start := time.Now()
			//             2. Создай statusRecorder с начальным статусом 200
			//             3. Вызови next.ServeHTTP(rec, r)
			//             4. Залогируй: method, path, status, duration_ms
			start := time.Now()

			rec := &statusRecorder{
				status:         http.StatusOK,
				ResponseWriter: w,
			}

			next.ServeHTTP(rec, r)

			duration := time.Since(start)

			logger.Info(
				"http request",
				"method", r.Method,
				"path", r.URL.Path,
				"status", rec.status,
				"duration_ms", duration.Milliseconds(),
			)
		})
	}
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	mux := http.NewServeMux()

	// TODO: оберни healthHandler в LoggingMiddleware и зарегистрируй на "GET /health"
	health := Chain(http.HandlerFunc(healthHandler), LoggingMiddleware(logger))
	mux.Handle("GET /health", health)

	fmt.Println("server started on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		fmt.Printf("server error: %v\n", err)
	}
}
