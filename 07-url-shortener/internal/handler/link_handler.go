package handler

import (
	"encoding/json"
	"net/http"
	"url-shortener/internal/service"
)

type LinkHandler struct {
	service *service.LinkService
}

func NewLinkHandler(serv *service.LinkService) *LinkHandler {
	return &LinkHandler{
		service: serv,
	}
}

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	ShortURL string `json:"short_url"`
}

func (h *LinkHandler) Shorten(w http.ResponseWriter, r *http.Request) {
	var req ShortenRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	link, err := h.service.CreateShortLink(r.Context(), req.URL)
	if err != nil {
		http.Error(w, "failed to create short link", http.StatusInternalServerError)
		return
	}

	response := ShortenResponse{
		ShortURL: "http://localhost:8080/" + link.ShortCode,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		return
	}

}

func (h *LinkHandler) Redirect(w http.ResponseWriter, r *http.Request) {}
