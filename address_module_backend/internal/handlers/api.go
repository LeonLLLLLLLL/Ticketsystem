package handlers

import (
	"address_module/internal/middleware"

	"github.com/go-chi/chi"
	chimiddle "github.com/go-chi/chi/middleware"
)

func Handler(r *chi.Mux) {
	// Global middleware
	r.Use(chimiddle.StripSlashes)

	r.Route("/account", func(router chi.Router) {
		// Middleware for /account route
		router.Use(middleware.Authorization)
		router.Get("/coins", GetCoinBalance)
	})

	r.Route("/firm", func(router chi.Router) {
		router.Use(middleware.Authorization)
		router.Post("/submit", AddFirm)
		router.Get("/get_by_id", GetFirmsByContactID)
		router.Get("/get", GetAllFirms)
	})
	
	r.Route("/contact", func(router chi.Router) {
		router.Use(middleware.Authorization)
		router.Post("/submit", AddContact)
		router.Get("/get_by_id", GetContactsByFirmID)
		router.Get("/get", GetAllContacts)
	})
}
