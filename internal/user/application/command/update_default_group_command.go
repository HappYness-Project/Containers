package command

import (
	"fmt"

	"github.com/happYness-Project/taskManagementGolang/internal/user/repository"
)

// UpdateDefaultGroupCommand represents the command to update user's default group
type UpdateDefaultGroupCommand struct {
	UserId         string // UUID
	DefaultGroupId int
}

// UpdateDefaultGroupCommandHandler handles updating a user's default group
type UpdateDefaultGroupCommandHandler struct {
	userRepo repository.UserRepository
}

func NewUpdateDefaultGroupCommandHandler(
	userRepo repository.UserRepository,
) *UpdateDefaultGroupCommandHandler {
	return &UpdateDefaultGroupCommandHandler{
		userRepo: userRepo,
	}
}

// Handle executes the update default group command
func (h *UpdateDefaultGroupCommandHandler) Handle(cmd UpdateDefaultGroupCommand) error {
	// Validate user exists
	user, err := h.userRepo.GetUserByUserId(cmd.UserId)
	if err != nil || user == nil {
		return fmt.Errorf("user not found: %s", cmd.UserId)
	}

	// Update default group using domain logic (includes validation)
	err = user.UpdateDefaultGroupId(cmd.DefaultGroupId)
	if err != nil {
		return fmt.Errorf("domain validation error: %w", err)
	}

	// Persist changes
	err = h.userRepo.UpdateUser(*user)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}
