package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"url-shortener/internal/service"
)

type UserServiceInterface interface {
	Register(ctx context.Context, email string, password string) error
	Login(ctx context.Context, email string, password string) (string, error)
}

type UserHandler struct {
	service UserServiceInterface
}

func NewUserHandler(Userservice UserServiceInterface) *UserHandler {
	return &UserHandler{
		service: Userservice,
	}

}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {

	var req RegisterRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	err = h.service.Register(r.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrInvalidRegistrationData) {
			http.Error(
				w,
				"invalid registration data",
				http.StatusBadRequest,
			)
			return
		}

		if errors.Is(err, service.ErrEmailAlreadyExists) {
			http.Error(w, "email already registered", http.StatusConflict)
			return
		}

		http.Error(
			w,
			"internal server error",
			http.StatusInternalServerError,
		)
		return

	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "user registered successfully")

}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	tokenString, err := h.service.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			http.Error(w, "invalid email or password", http.StatusUnauthorized)
			return
		}

		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}
