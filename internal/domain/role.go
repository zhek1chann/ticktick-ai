package domain

// Role represents user role type
type Role string

// User roles
const (
	RoleUser       Role = "user"
	RoleShopOwner  Role = "shop_owner"
	RoleSuperAdmin Role = "super_admin"
)

// String returns string representation of role
func (r Role) String() string {
	return string(r)
}
