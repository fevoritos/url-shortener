package transporthttp

import (
	"net/http"

	_ "url-shortener/internal/transport/http/docs"
	httphandlers "url-shortener/internal/transport/http/handlers"

	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func NewRouter(linkHandler *httphandlers.LinkHandler) *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("POST /links", linkHandler.Create)
	router.HandleFunc("GET /{hash}", linkHandler.Goto)

	router.Handle("GET /docs/{any...}", httpSwagger.WrapHandler)

	return router
}
