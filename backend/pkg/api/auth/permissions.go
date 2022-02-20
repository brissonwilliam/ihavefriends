package auth

const (
	PERM_ADD_USER    = "ADD_USER"
	PERM_SUPER_ADMIN = "SUPER_ADMIN"
)

func HasSuperAdminPermission(userPermissions []string) bool {
	for _, p := range userPermissions {
		if p == PERM_SUPER_ADMIN {
			return true
		}
	}
	return false
}

func HasPermission(requiredPermission string, userPermissions []string) bool {
	if HasSuperAdminPermission(userPermissions) {
		return true
	}
	for _, p := range userPermissions {
		if p == requiredPermission {
			return true
		}
	}
	return false
}
