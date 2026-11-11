package handler

import "net/http"

func NewRouter(linkHand *LinkHandler, userHandler *UserHandler, jwtSecret string) *http.ServeMux {

	mux := http.NewServeMux()

	mux.HandleFunc("POST /shorten", AuthMiddleware(jwtSecret, linkHand.Shorten))
	mux.HandleFunc("GET /{code}", linkHand.Redirect)
	mux.HandleFunc("GET /stats/{code}", linkHand.Stats)

	mux.HandleFunc("POST /register", userHandler.Register)
	mux.HandleFunc("POST /login", userHandler.Login)

	return mux
}
