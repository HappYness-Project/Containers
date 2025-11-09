package builders

import (
	"github.com/happYness-Project/taskManagementGolang/internal/usergroup/domain"
)

// UserGroupBuilder provides a fluent interface for building test user groups
type UserGroupBuilder struct {
	group *domain.UserGroup
	err   error
}

// NewUserGroupBuilder creates a new user group builder with sensible defaults
func NewUserGroupBuilder() *UserGroupBuilder {
	group, err := domain.NewUserGroup(
		"Test Group",
		"Test group description",
		"team",
	)

	return &UserGroupBuilder{
		group: group,
		err:   err,
	}
}

// WithName sets the group name
func (b *UserGroupBuilder) WithName(name string) *UserGroupBuilder {
	if b.group != nil {
		b.group.GroupName = name
	}
	return b
}

// WithDescription sets the group description
func (b *UserGroupBuilder) WithDescription(desc string) *UserGroupBuilder {
	if b.group != nil {
		b.group.GroupDesc = desc
	}
	return b
}

// WithType sets the group type
func (b *UserGroupBuilder) WithType(groupType string) *UserGroupBuilder {
	if b.group != nil {
		b.group.Type = groupType
	}
	return b
}

// WithThumbnail sets the thumbnail URL
func (b *UserGroupBuilder) WithThumbnail(thumbnail string) *UserGroupBuilder {
	if b.group != nil {
		b.group.Thumbnail = thumbnail
	}
	return b
}

// WithGroupId sets a specific group ID (useful for testing)
func (b *UserGroupBuilder) WithGroupId(id int) *UserGroupBuilder {
	if b.group != nil {
		b.group.GroupId = id
	}
	return b
}

// Inactive marks the group as inactive
func (b *UserGroupBuilder) Inactive() *UserGroupBuilder {
	if b.group != nil {
		b.group.IsActive = false
	}
	return b
}

// TeamType sets the type to team
func (b *UserGroupBuilder) TeamType() *UserGroupBuilder {
	return b.WithType("team")
}

// ProjectType sets the type to project
func (b *UserGroupBuilder) ProjectType() *UserGroupBuilder {
	return b.WithType("project")
}

// PersonalType sets the type to personal
func (b *UserGroupBuilder) PersonalType() *UserGroupBuilder {
	return b.WithType("personal")
}

// Build returns the built group or an error
func (b *UserGroupBuilder) Build() (*domain.UserGroup, error) {
	return b.group, b.err
}

// MustBuild returns the built group or panics (useful in test setup)
func (b *UserGroupBuilder) MustBuild() *domain.UserGroup {
	if b.err != nil {
		panic(b.err)
	}
	return b.group
}
