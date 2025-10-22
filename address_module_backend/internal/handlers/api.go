package handlers

import (
	"address_module/internal/middleware"

	"github.com/go-chi/chi"
	chimiddle "github.com/go-chi/chi/middleware"
)

func Handler(r *chi.Mux) {
	// Global middleware
	r.Use(chimiddle.StripSlashes)

	// Login & Register
	r.Route("/auth", func(router chi.Router) {
		router.Post("/login", LoginHandler) // points to middleware package now
		router.Post("/register", RegisterHandler)
	})

	// Users
	r.Route("/users", func(router chi.Router) {
		router.Use(middleware.Authorization)
		router.With(middleware.RequirePermission("create_users")).Post("/create", AddUser)
		router.With(middleware.RequirePermission("view_users")).Get("/get", GetUserByID)        // expects ?id=
		router.With(middleware.RequirePermission("edit_users")).Put("/update", UpdateUser)      // expects full user JSON body with ID
		router.With(middleware.RequirePermission("delete_users")).Delete("/delete", DeleteUser) // expects ?id=
	})

	// Roles
	r.Route("/roles", func(router chi.Router) {
		router.Use(middleware.Authorization)
		router.With(middleware.RequirePermission("create_roles")).Post("/create", AddRole)
		router.With(middleware.RequirePermission("view_roles")).Get("/get", GetRoleByID) // expects ?id=
		router.With(middleware.RequirePermission("edit_roles")).Put("/update", UpdateRole)
		router.With(middleware.RequirePermission("delete_roles")).Delete("/delete", DeleteRole) // expects ?id=
	})

	// Permissions
	r.Route("/permissions", func(router chi.Router) {
		router.Use(middleware.Authorization)
		router.With(middleware.RequirePermission("create_permissions")).Post("/create", AddPermission)
		router.With(middleware.RequirePermission("view_permissions")).Get("/get", GetPermissionByID) // expects ?id=
		router.With(middleware.RequirePermission("edit_permissions")).Put("/update", UpdatePermission)
		router.With(middleware.RequirePermission("delete_permissions")).Delete("/delete", DeletePermission) // expects ?id=
	})

	// User-Role Assignments
	r.Route("/user_roles", func(router chi.Router) {
		router.Use(middleware.Authorization)
		router.With(middleware.RequirePermission("assign_roles")).Post("/assign", AssignUserRole)
		router.With(middleware.RequirePermission("unassign_roles")).Delete("/remove", RemoveUserRole) // expects ?user_id=&role_id=
	})

	// Role-Permission Assignments
	r.Route("/role_permissions", func(router chi.Router) {
		router.Use(middleware.Authorization)
		router.With(middleware.RequirePermission("assign_permissions")).Post("/assign", AssignRolePermission)
		router.With(middleware.RequirePermission("unassign_permissions")).Delete("/remove", RemoveRolePermission) // expects ?role_id=&permission_id=
	})

	r.Route("/firm", func(router chi.Router) {
		router.Use(middleware.Authorization)
		router.With(middleware.RequirePermission("create_firms")).Post("/submit", AddFirm)
		router.With(middleware.RequirePermission("view_firms")).Get("/get", GetAllFirms)
	})

	r.Route("/contact", func(router chi.Router) {
		router.Use(middleware.Authorization)
		router.With(middleware.RequirePermission("create_contacts")).Post("/submit", AddContact)
		router.With(middleware.RequirePermission("view_contacts")).Get("/get", GetAllContacts)
	})

	// ✅ Devices (no permission middleware)
	r.Route("/devices", func(router chi.Router) {
		router.Post("/create", AddDevice)
		router.Get("/get", GetDeviceByID) // expects ?id=
		router.Get("/list", ListDevices)
		router.Put("/update", UpdateDevice)
		router.Delete("/delete", DeleteDevice) // expects ?id=
	})

	// ✅ Device Links (no permission middleware)
	r.Route("/device_links", func(router chi.Router) {
		router.Post("/create", AddDeviceLink)
		router.Put("/update", UpdateDeviceLink)
		router.Get("/list", ListDeviceLinks)
		router.Get("/get", GetDeviceLinksForDevice) // expects ?device_id=
		router.Delete("/delete", DeleteDeviceLink)  // expects ?id=
	})
}
