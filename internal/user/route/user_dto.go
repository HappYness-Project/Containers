package route

import (
	"time"

	"github.com/happYness-Project/taskManagementGolang/internal/usergroup/domain"
)

type CreateUserDto struct {
	UserName  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}
type UpdateUserDto struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}
type UserDetailDto struct {
	Id             int                 `json:"id"`
	UserId         string              `json:"user_id"`
	UserName       string              `json:"username"`
	FirstName      string              `json:"first_name"`
	LastName       string              `json:"last_name"`
	CreatedAt      time.Time           `json:"created_at"`
	UpdatedAt      time.Time           `json:"updated_at"`
	Email          string              `json:"email"`
	IsActive       bool                `json:"is_active"`
	DefaultGroupId int                 `json:"default_group_id"`
	UserGroup      []*domain.UserGroup `json:"user_groups"`
}

type UserWithRoleDto struct {
	Id             int       `json:"id"`
	UserId         string    `json:"user_id"`
	UserName       string    `json:"username"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Email          string    `json:"email"`
	IsActive       bool      `json:"is_active"`
	DefaultGroupId int       `json:"default_group_id"`
	Role           string    `json:"role"`
	JoinedAt       time.Time `json:"joined_at"`
}
