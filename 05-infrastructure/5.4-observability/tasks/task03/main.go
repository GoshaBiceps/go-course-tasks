// Задание: Нормализация route label для метрик
//
// При использовании сырого r.URL.Path в метриках path-параметры создают
// отдельные временные ряды для каждого ID. Нужен стабильный route-label.
//
// Ожидаемый результат:
//   /api/v1/users/123     → /api/v1/users/{id}
//   /api/v1/orders/abc    → /api/v1/orders/{id}
//   /health               → /health
//   /api/v1/unknown/x/y   → /unknown

package main

import (
	"fmt"
	"net/http"
	"strings"
)

// routePattern возвращает шаблон маршрута для использования в метриках.
// В Go 1.22+ шаблон доступен через r.Pattern.
//
// TODO: реализуй функцию normalizeRoute(r *http.Request) string
// Используй r.Pattern, если он задан.
// Если пустой (маршрут не зарегистрирован) — верни "/unknown".
//
// Подсказка: r.Pattern доступен начиная с Go 1.22.

func normalizeRoute(r *http.Request) string {
	if r.Pattern == "" {
		return "/unknown"
	}

	parts := strings.SplitN(
		r.Pattern,
		" ",
		2,
	)

	return parts[1]
}

func main() {
	mux := http.NewServeMux()

	// Регистрируем маршруты, чтобы проверить, что pattern правильно извлекается
	mux.HandleFunc("GET /api/v1/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL.Path, "→", normalizeRoute(r))
	})
	mux.HandleFunc("GET /api/v1/orders/{id}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL.Path, "→", normalizeRoute(r))
	})
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL.Path, "→", normalizeRoute(r))
	})

	// Симулируем несколько запросов через тестовые объекты
	for _, path := range []string{
		"/api/v1/users/123",
		"/api/v1/orders/abc",
		"/health",
	} {
		req, _ := http.NewRequest("GET", path, nil)
		// Используем mux для маршрутизации, чтобы r.Pattern заполнился
		_, pattern := mux.Handler(req)
		req.Pattern = pattern
		fmt.Printf("%-30s → %s\n", path, normalizeRoute(req))
	}
}
