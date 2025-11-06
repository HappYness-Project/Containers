package domain

import "fmt"

// Role represents a user's role in a group
type Role string

const (
	RoleAdmin  Role = "admin"
	RoleMember Role = "member"
)

// NewRole creates a new Role value object with validation
func NewRole(value string) (Role, error) {
	role := Role(value)
	if !role.IsValid() {
		return "", fmt.Errorf("invalid role: %s. Must be 'admin' or 'member'", value)
	}
	return role, nil
}

// IsValid checks if the role is a valid value
func (r Role) IsValid() bool {
	return r == RoleAdmin || r == RoleMember
}

// String returns the string representation
func (r Role) String() string {
	return string(r)
}

// IsAdmin checks if the role is admin
func (r Role) IsAdmin() bool {
	return r == RoleAdmin
}

// IsMember checks if the role is member
func (r Role) IsMember() bool {
	return r == RoleMember
}
