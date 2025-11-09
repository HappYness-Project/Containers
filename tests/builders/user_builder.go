package builders

import (
	"fmt"

	"github.com/happYness-Project/taskManagementGolang/internal/user/domain"
	"github.com/google/uuid"
)

// UserBuilder provides a fluent interface for building test users
type UserBuilder struct {
	user *domain.User
}

// NewUserBuilder creates a new user builder with sensible defaults
func NewUserBuilder() *UserBuilder {
	userId := uuid.New().String()
	user := domain.NewUser(
		userId,
		fmt.Sprintf("testuser_%s", userId[:8]),
		"Test",
		"User",
		fmt.Sprintf("test_%s@example.com", userId[:8]),
	)

	return &UserBuilder{user: user}
}

// WithUserId sets a specific user ID
func (b *UserBuilder) WithUserId(userId string) *UserBuilder {
	b.user.UserId = userId
	return b
}

// WithUserName sets the username
func (b *UserBuilder) WithUserName(userName string) *UserBuilder {
	b.user.UserName = userName
	return b
}

// WithFirstName sets the first name
func (b *UserBuilder) WithFirstName(firstName string) *UserBuilder {
	b.user.FirstName = firstName
	return b
}

// WithLastName sets the last name
func (b *UserBuilder) WithLastName(lastName string) *UserBuilder {
	b.user.LastName = lastName
	return b
}

// WithEmail sets the email
func (b *UserBuilder) WithEmail(email string) *UserBuilder {
	b.user.Email = email
	return b
}

// WithDefaultGroupId sets the default group ID
func (b *UserBuilder) WithDefaultGroupId(groupId int) *UserBuilder {
	b.user.UpdateDefaultGroupId(groupId)
	return b
}

// Inactive marks the user as inactive
func (b *UserBuilder) Inactive() *UserBuilder {
	b.user.IsActive = false
	return b
}

// WithFullName sets both first and last name
func (b *UserBuilder) WithFullName(firstName, lastName string) *UserBuilder {
	b.user.FirstName = firstName
	b.user.LastName = lastName
	return b
}

// Build returns the built user
func (b *UserBuilder) Build() *domain.User {
	return b.user
}
