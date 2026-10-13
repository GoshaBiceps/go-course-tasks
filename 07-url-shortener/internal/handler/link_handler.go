package handler

import (
	"encoding/json"
	"net/http"
	"url-shortener/internal/service"
)

type LinkHandler struct {
	service *service.LinkService
	baseURL string
}

func NewLinkHandler(serv *service.LinkService, baseURL string) *LinkHandler {
	return &LinkHandler{
		service: serv,
		baseURL: baseURL,
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
		ShortURL: h.baseURL + "/" + link.ShortCode,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		return
	}

}

func (h *LinkHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	shortCode := r.PathValue("code")

	originalURL, err := h.service.GetOriginalURL(
		r.Context(),
		shortCode,
	)

	if err != nil {
		http.Error(w, "link not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, originalURL, http.StatusFound) // перенавравляем по длинному урлу

}

type StatsResponse struct {
	ShortCode string `json:"short_code"`
	LongURL   string `json:"long_url"`
	Clicks    int    `json:"clicks"`
}

func (h *LinkHandler) Stats(w http.ResponseWriter, r *http.Request) {

	shortCode := r.PathValue("code")

	link, err := h.service.GetStats(r.Context(), shortCode)

	if err != nil {
		http.Error(w, "link not found", http.StatusNotFound)
		return
	}
	response := StatsResponse{
		ShortCode: link.ShortCode,
		LongURL:   link.LongURL,
		Clicks:    link.Clicks,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		return
	}
}
