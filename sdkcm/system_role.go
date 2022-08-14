package sdkcm

type SystemRole int

const (
	SysRoleRoot SystemRole = iota
	SysRoleAdmin
	SysRoleModerator
	SysRoleUser
	SysRoleGuest
)

func AllSysRoles() []string {
	return []string{"root", "admin", "moderator", "user", "guest"}
}

func (sr SystemRole) String() string {
	return AllSysRoles()[sr]
}

func ParseSystemRole(role string) SystemRole {
	for i, v := range AllSysRoles() {
		if role == v {
			return SystemRole(i)
		}
	}
	return SysRoleUser
}
