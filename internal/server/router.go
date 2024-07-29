package router

import (
	"context"
	"net/http"

	"github.com/alexPavlikov/go-wiki/internal/server/locations"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type Router struct {
	handler *locations.Handler
}

func New(handler *locations.Handler) *Router {
	return &Router{
		handler: handler,
	}
}

func (r *Router) Build() http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Post("/v1/wiki", middlewares(r.handler.WikiHandler))

	return router
}

func middlewares(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(context.WithValue(r.Context(), "START_LINK", r.Header.Get("X-WIKILINK-START")))
		r = r.WithContext(context.WithValue(r.Context(), "END_LINK", r.Header.Get("X-WIKILINK-END")))

		h.ServeHTTP(w, r)
	}
}
