package transporthttp

import (
	"net/http"
	httphandlers "url-shortener/internal/transport/http/handlers"
)

func NewRouter(linkHandler *httphandlers.LinkHandler) *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("POST /link", linkHandler.Create)
	router.HandleFunc("GET /{hash}", linkHandler.Goto)

	return router
}
