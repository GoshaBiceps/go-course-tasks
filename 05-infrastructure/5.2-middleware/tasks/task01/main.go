// Задание: Chain helper
//
// Реализуй тип Middleware и функцию Chain, которая оборачивает handler
// в цепочку middleware в правильном порядке.
//
// Проверь порядок на двух тестовых middleware, которые печатают "before"/"after".
//
// Ожидаемый результат:
//   $ go run main.go &
//   server started on :8080
//
//   $ curl http://localhost:8080/hello
//   [mw1] before
//   [mw2] before
//   [mw2] after
//   [mw1] after
//   hello, world

package main

import (
	"fmt"
	"net/http"
)

// TODO: определи тип Middleware
type Middleware func(http.Handler) http.Handler

// TODO: реализуй функцию Chain
func Chain(h http.Handler, middlewares ...Middleware) http.Handler {
	//
	// Подсказка: оборачивай handler от последнего middleware к первому:
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}

	return h
}

// TODO: реализуй testMiddleware, которая печатает "[name] before" до вызова next
// и "[name] after" после вызова next.
func testMiddleware(name string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Printf("[%s] before\n", name)
			next.ServeHTTP(w, r)
			fmt.Printf("[%s] after\n", name)
		})
	}
}

func helloHandler(w http.ResponseWriter, _ *http.Request) {
	_, _ = fmt.Fprintln(w, "hello, world")
}

func main() {
	mux := http.NewServeMux()

	// TODO: собери цепочку из двух testMiddleware ("mw1", "mw2") вокруг helloHandler
	handler := Chain(http.HandlerFunc(helloHandler), testMiddleware("mw1"), testMiddleware("mw2"))
	mux.Handle("GET /hello", handler)

	fmt.Println("server started on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		fmt.Printf("server error: %v\n", err)
	}
}
