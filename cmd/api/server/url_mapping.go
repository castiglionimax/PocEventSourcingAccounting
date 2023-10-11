package server

import (
	"net/http"

	"github.com/go-chi/chi"

	"github.com/castiglionimax/challengeXepelin/internal/controller"
)

type mapping struct {
	controller controller.Controller
}

func newMapping() *mapping {
	return &mapping{
		controller: resolveController(),
	}
}

func (m mapping) mapUrlsToControllers(route *chi.Mux) {
	route.Get("/ping", alive())
	route.Post("/accounts", m.controller.CreateAccount)
	route.With(loggerHVTransaction).Post("/transactions", m.controller.Transaction)

	route.Get("/accounts/{id}/balance", m.controller.Balance)

}

func alive() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("pong"))
	}
}
