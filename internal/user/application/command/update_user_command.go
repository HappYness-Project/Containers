package command

import (
	"fmt"

	"github.com/happYness-Project/taskManagementGolang/internal/user/domain"
	"github.com/happYness-Project/taskManagementGolang/internal/user/repository"
)

// UpdateUserCommand represents the command to update a user's profile
type UpdateUserCommand struct {
	UserId    string // UUID
	FirstName string
	LastName  string
	Email     string
}

// UpdateUserCommandHandler handles updating a user
type UpdateUserCommandHandler struct {
	userRepo repository.UserRepository
}

func NewUpdateUserCommandHandler(
	userRepo repository.UserRepository,
) *UpdateUserCommandHandler {
	return &UpdateUserCommandHandler{
		userRepo: userRepo,
	}
}

// Handle executes the update user command
func (h *UpdateUserCommandHandler) Handle(cmd UpdateUserCommand) error {
	// Validate user exists
	user, err := h.userRepo.GetUserByUserId(cmd.UserId)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return domain.ErrUserNotFound
	}

	// Update user using domain logic
	user.UpdateUser(cmd.FirstName, cmd.LastName, cmd.Email)

	// Persist changes
	err = h.userRepo.UpdateUser(*user)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}
