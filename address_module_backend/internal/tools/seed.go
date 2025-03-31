package tools

import (
	"address_module/internal/model"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

func (db *MySQLDB) SeedInitialData() error {
	log.Info("🔧 Starting seed process for roles, permissions and admin user")

	permissionNames := map[string]string{
		"view_users":           "View users",
		"edit_users":           "Edit users",
		"create_users":         "Create users",
		"delete_users":         "Delete users",
		"view_firms":           "View firms",
		"edit_firms":           "Edit firms",
		"create_firms":         "Create firms",
		"delete_firms":         "Delete firms",
		"view_contacts":        "View contacts",
		"edit_contacts":        "Edit contacts",
		"create_contacts":      "Create contacts",
		"delete_contacts":      "Delete contacts",
		"view_roles":           "View roles",
		"edit_roles":           "Edit roles",
		"create_roles":         "Create roles",
		"delete_roles":         "Delete roles",
		"view_permissions":     "View permissions",
		"edit_permissions":     "Edit permissions",
		"create_permissions":   "Create permissions",
		"delete_permissions":   "Delete permissions",
		"assign_roles":         "Assign roles",
		"unassign_roles":       "Unassign roles",
		"assign_permissions":   "Assign permissions",
		"unassign_permissions": "Unassign permissions",
		"admin_panel":          "Access admin panel",
	}

	for name, desc := range permissionNames {
		if _, err := db.GetPermissionByName(name); err != nil {
			p := model.Permission{Name: name, Description: desc}
			if _, err := db.InsertPermission(p); err != nil {
				log.Errorf("❌ Failed to insert permission '%s': %v", name, err)
			} else {
				log.Infof("✅ Inserted permission: %s", name)
			}
		}
	}

	adminRoleName := "admin"
	adminRole, err := db.GetRoleByName(adminRoleName)
	if err != nil {
		roleID, err := db.InsertRole(model.Role{
			Name:        adminRoleName,
			Description: "Default admin role with full permissions",
		})
		if err != nil {
			log.Errorf("❌ Failed to create admin role: %v", err)
			return err
		}
		log.Infof("✅ Created admin role with ID %d", roleID)
		adminRole = &model.Role{ID: roleID, Name: adminRoleName}
	}

	adminEmail := "admin@system.local"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost) // ⚠️ Use bcrypt in prod
	if err != nil {
		log.Errorf("❌ Failed to create hashed admin password: %v", err)
		return err
	}
	adminUser, err := db.GetUserByEmail(adminEmail)
	if err != nil {
		adminUser = &model.User{
			Username:       "admin",
			Email:          adminEmail,
			HashedPassword: string(hashedPassword),
			CreatedAt:      time.Now(),
		}
		uid, err := db.InsertUser(*adminUser)
		if err != nil {
			log.Errorf("❌ Failed to insert admin user: %v", err)
			return err
		}
		log.Infof("✅ Created admin user with ID %d", uid)
		adminUser.ID = uid
	}

	for name := range permissionNames {
		perm, err := db.GetPermissionByName(name)
		if err != nil {
			log.Warnf("⚠️ Permission %s not found. Skipping assignment to admin role", name)
			continue
		}
		err = db.InsertRolePermission(model.RolePermission{
			RoleID:       adminRole.ID,
			PermissionID: perm.ID,
		})
		if err != nil {
			log.Warnf("🔁 Skipped existing role-permission '%s' → '%s'", adminRole.Name, name)
		}
	}

	err = db.InsertUserRole(model.UserRole{
		UserID: adminUser.ID,
		RoleID: adminRole.ID,
	})
	if err != nil {
		log.Warnf("🔁 Admin role already assigned to admin user or failed: %v", err)
	}

	log.Info("🎉 Seeding complete.")
	return nil
}
