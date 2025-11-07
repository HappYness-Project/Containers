package domain

import (
	"errors"
	"time"
)

type User struct {
	Id             int       `json:"id"`
	UserId         string    `json:"user_id"`
	UserName       string    `json:"username"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	Email          string    `json:"email"`
	IsActive       bool      `json:"is_active"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	DefaultGroupId int       `json:"default_group_id"`
}

func NewUser(userId string, userName string, firstName string, lastName string, email string) *User {
	user := User{
		UserId:         userId,
		UserName:       userName,
		FirstName:      firstName,
		LastName:       lastName,
		Email:          email,
		IsActive:       true,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		DefaultGroupId: 0,
	}

	return &user
}

func (u *User) UpdateDefaultGroupId(groupId int) error {
	if groupId < 0 {
		return errors.New("group ID cannot be negative")
	}

	if u.DefaultGroupId == groupId {
		return errors.New("group ID is already set to the specified value")
	}

	u.DefaultGroupId = groupId
	u.UpdatedAt = time.Now()

	return nil
}

func (u *User) ClearDefaultGroup() {
	u.DefaultGroupId = 0
	u.UpdatedAt = time.Now()
}

func (u *User) UpdateUser(firstName, lastName, email string) {
	u.FirstName = firstName
	u.LastName = lastName
	u.Email = email
	u.UpdatedAt = time.Now()
}

type UserWithRole struct {
	*User
	Role     string    `json:"role"`
	JoinedAt time.Time `json:"joined_at"`
}
