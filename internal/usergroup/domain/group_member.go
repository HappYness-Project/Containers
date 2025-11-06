package domain

import "time"

// GroupMember represents a member of a user group (entity within the aggregate)
type GroupMember struct {
	UserId   string
	Role     Role
	JoinedAt time.Time
}

// NewGroupMember creates a new group member with the specified role
func NewGroupMember(userId string, role Role) *GroupMember {
	return &GroupMember{
		UserId:   userId,
		Role:     role,
		JoinedAt: time.Now(),
	}
}

// PromoteToAdmin promotes the member to admin role
func (m *GroupMember) PromoteToAdmin() {
	m.Role = RoleAdmin
}

// DemoteToMember demotes the member to regular member role
func (m *GroupMember) DemoteToMember() {
	m.Role = RoleMember
}

// ChangeRole changes the member's role
func (m *GroupMember) ChangeRole(newRole Role) {
	m.Role = newRole
}
