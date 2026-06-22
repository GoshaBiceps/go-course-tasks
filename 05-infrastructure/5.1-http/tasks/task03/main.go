// Задание: JSON endpoint с валидацией
//
// Реализуй POST /api/v1/orders, который принимает JSON-тело,
// валидирует входные данные и возвращает JSON-ответ.
//
// Ожидаемый результат:
//   $ go run main.go
//   2025/01/01 12:00:00 server started on :8080
//
//   $ curl -X POST http://localhost:8080/api/v1/orders -d '{"item":"laptop","quantity":2}'
//   {"order_id":"...","item":"laptop","quantity":2,"status":"created"}
//
//   $ curl -X POST http://localhost:8080/api/v1/orders -d '{"item":"","quantity":1}'
//   {"error":"item is required"}
//
//   $ curl -X POST http://localhost:8080/api/v1/orders -d 'not json'
//   {"error":"invalid json"}

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type orderRequest struct {
	Item     string `json:"item"`
	Quantity int    `json:"quantity"`
}

type orderResponse struct {
	OrderID  string `json:"order_id"`
	Item     string `json:"item"`
	Quantity int    `json:"quantity"`
	Status   string `json:"status"`
}

type apiError struct {
	Error string `json:"error"`
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/v1/orders", handleCreateOrder)

	log.Println("server started on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}

func handleCreateOrder(w http.ResponseWriter, r *http.Request) {
	// TODO: декодируй JSON-тело запроса в orderRequest
	// Подсказка: json.NewDecoder(r.Body).Decode(&req)
	// При ошибке декодирования верни 400 с apiError{Error: "invalid json"}
	req := orderRequest{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, apiError{Error: "invalid json"})
		return
	}

	// TODO: проверь, что поле Item не пустое
	// Если пустое — верни 400 с apiError{Error: "item is required"}
	if req.Item == "" {
		writeJSON(w, http.StatusBadRequest, apiError{Error: "item is required"})
		return
	}

	// TODO: сформируй orderResponse:
	//   OrderID:  fmt.Sprintf("ord-%d", time.Now().UnixNano())
	//   Item:     req.Item
	//   Quantity: req.Quantity
	//   Status:   "created"
	resp := orderResponse{
		OrderID:  fmt.Sprintf("ord-%d", time.Now().UnixNano()),
		Item:     req.Item,
		Quantity: req.Quantity,
		Status:   "created",
	}

	// TODO: верни ответ со статусом 201 (http.StatusCreated)
	// Подсказка: writeJSON(w, http.StatusCreated, resp)
	writeJSON(w, http.StatusCreated, resp)

}
