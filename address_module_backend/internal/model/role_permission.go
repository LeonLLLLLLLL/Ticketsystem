package model

type RolePermissionParams struct {
	RoleID       int64 `json:"role_id"`
	PermissionID int64 `json:"permission_id"`
}
