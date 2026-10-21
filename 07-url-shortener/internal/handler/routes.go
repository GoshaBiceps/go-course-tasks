package handler

import "net/http"

func NewRouter(hand *LinkHandler) *http.ServeMux {

	mux := http.NewServeMux()

	mux.HandleFunc("POST /shorten", hand.Shorten)
	mux.HandleFunc("GET /{code}", hand.Redirect)
	mux.HandleFunc("GET /stats/{code}", hand.Stats)

	return mux
}
