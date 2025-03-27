package handlers

import (
	"address_module/internal/middleware"

	"github.com/go-chi/chi"
	chimiddle "github.com/go-chi/chi/middleware"
)

func Handler(r *chi.Mux) {
	// Global middleware
	r.Use(chimiddle.StripSlashes)

	// Users
	r.Route("/users", func(router chi.Router) {
		router.Use(middleware.Authorization)
		router.Post("/create", AddUser)
		router.Get("/get", GetUserByID)      // expects ?id=
		router.Put("/update", UpdateUser)    // expects full user JSON body with ID
		router.Delete("/delete", DeleteUser) // expects ?id=
	})

	// Roles
	r.Route("/roles", func(router chi.Router) {
		router.Use(middleware.Authorization)
		router.Post("/create", AddRole)
		router.Get("/get", GetRoleByID) // expects ?id=
		router.Put("/update", UpdateRole)
		router.Delete("/delete", DeleteRole) // expects ?id=
	})

	// Permissions
	r.Route("/permissions", func(router chi.Router) {
		router.Use(middleware.Authorization)
		router.Post("/create", AddPermission)
		router.Get("/get", GetPermissionByID) // expects ?id=
		router.Put("/update", UpdatePermission)
		router.Delete("/delete", DeletePermission) // expects ?id=
	})

	// User-Role Assignments
	r.Route("/user_roles", func(router chi.Router) {
		router.Use(middleware.Authorization)
		router.Post("/assign", AssignUserRole)
		router.Delete("/remove", RemoveUserRole) // expects ?user_id=&role_id=
	})

	// Role-Permission Assignments
	r.Route("/role_permissions", func(router chi.Router) {
		router.Use(middleware.Authorization)
		router.Post("/assign", AssignRolePermission)
		router.Delete("/remove", RemoveRolePermission) // expects ?role_id=&permission_id=
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
		//router.Get("/get_by_id", GetContactsByFirmID)
		router.Get("/get", GetAllContacts)
	})
}
