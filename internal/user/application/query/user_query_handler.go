package query

import (
	"fmt"

	"github.com/happYness-Project/taskManagementGolang/internal/user/domain"
	"github.com/happYness-Project/taskManagementGolang/internal/user/repository"
)

// Queries define the read operations
type GetAllUsersQuery struct{}

type GetUserByIdQuery struct {
	UserId string // UUID
}

type GetUserByEmailQuery struct {
	Email string
}

type GetUserByUsernameQuery struct {
	Username string
}

type GetUsersByGroupIdQuery struct {
	GroupId int
}

// UserQueryHandler handles all read operations for User
type UserQueryHandler struct {
	userRepo repository.UserRepository
}

func NewUserQueryHandler(
	userRepo repository.UserRepository,
) *UserQueryHandler {
	return &UserQueryHandler{
		userRepo: userRepo,
	}
}

// HandleGetAllUsers retrieves all users
func (h *UserQueryHandler) HandleGetAllUsers(query GetAllUsersQuery) ([]*domain.User, error) {
	users, err := h.userRepo.GetAllUsers()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve users: %w", err)
	}
	return users, nil
}

// HandleGetUserById retrieves a single user by ID
func (h *UserQueryHandler) HandleGetUserById(query GetUserByIdQuery) (*domain.User, error) {
	user, err := h.userRepo.GetUserByUserId(query.UserId)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve user: %w", err)
	}

	if user == nil {
		return nil, domain.ErrUserNotFound
	}

	return user, nil
}

// HandleGetUserByEmail retrieves a user by email
func (h *UserQueryHandler) HandleGetUserByEmail(query GetUserByEmailQuery) (*domain.User, error) {
	user, err := h.userRepo.GetUserByEmail(query.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve user by email: %w", err)
	}

	if user == nil || user.Id == 0 {
		return nil, fmt.Errorf("user not found with email: %s", query.Email)
	}

	return user, nil
}

// HandleGetUserByUsername retrieves a user by username
func (h *UserQueryHandler) HandleGetUserByUsername(query GetUserByUsernameQuery) (*domain.User, error) {
	user, err := h.userRepo.GetUserByUsername(query.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve user by username: %w", err)
	}

	if user == nil || user.Id == 0 {
		return nil, fmt.Errorf("user not found with username: %s", query.Username)
	}

	return user, nil
}

// HandleGetUsersByGroupId retrieves all users in a group with roles
func (h *UserQueryHandler) HandleGetUsersByGroupId(query GetUsersByGroupIdQuery) ([]*domain.UserWithRole, error) {
	users, err := h.userRepo.GetUsersByGroupIdWithRoles(query.GroupId)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve users by group: %w", err)
	}
	return users, nil
}
