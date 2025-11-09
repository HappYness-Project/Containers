package builders

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/happYness-Project/taskManagementGolang/internal/user/domain"
)

type UserBuilder struct {
	user *domain.User
}

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

func (b *UserBuilder) WithUserId(userId string) *UserBuilder {
	b.user.UserId = userId
	return b
}

func (b *UserBuilder) WithUserName(userName string) *UserBuilder {
	b.user.UserName = userName
	return b
}

func (b *UserBuilder) WithFirstName(firstName string) *UserBuilder {
	b.user.FirstName = firstName
	return b
}

func (b *UserBuilder) WithLastName(lastName string) *UserBuilder {
	b.user.LastName = lastName
	return b
}

func (b *UserBuilder) WithEmail(email string) *UserBuilder {
	b.user.Email = email
	return b
}

func (b *UserBuilder) WithDefaultGroupId(groupId int) *UserBuilder {
	b.user.UpdateDefaultGroupId(groupId)
	return b
}

func (b *UserBuilder) Inactive() *UserBuilder {
	b.user.IsActive = false
	return b
}

func (b *UserBuilder) WithFullName(firstName, lastName string) *UserBuilder {
	b.user.FirstName = firstName
	b.user.LastName = lastName
	return b
}

func (b *UserBuilder) Build() *domain.User {
	return b.user
}
