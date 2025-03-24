package handlers

import (
	"address_module/internal/middleware"

	"github.com/go-chi/chi"
	chimiddle "github.com/go-chi/chi/middleware"
)

func Handler(r *chi.Mux) {
	// Global middleware
	r.Use(chimiddle.StripSlashes)

	r.Route("/firm", func(router chi.Router) {
		router.Use(middleware.Authorization)
		router.Post("/submit", AddFirm)
		router.Get("/get_by_id", GetFirmsByContactID)
		router.Get("/get", GetAllFirms)
	})

	r.Route("/contact", func(router chi.Router) {
		router.Use(middleware.Authorization)
		router.Post("/submit", AddContact)
		//router.Get("/get_by_id", GetContactsByFirmID)
		router.Get("/get", GetAllContacts)
	})
}
