package command

import (
	"fmt"

	"github.com/happYness-Project/taskManagementGolang/internal/user/repository"
	userGroupRepo "github.com/happYness-Project/taskManagementGolang/internal/usergroup/repository"
)

// RemoveMemberCommand represents the command to remove a member from a group
type RemoveMemberCommand struct {
	GroupId int
	UserId  string // UUID
}

// RemoveMemberCommandHandler handles removing a member from a group
type RemoveMemberCommandHandler struct {
	groupRepo userGroupRepo.UserGroupRepository
	userRepo  repository.UserRepository
}

func NewRemoveMemberCommandHandler(
	groupRepo userGroupRepo.UserGroupRepository,
	userRepo repository.UserRepository,
) *RemoveMemberCommandHandler {
	return &RemoveMemberCommandHandler{
		groupRepo: groupRepo,
		userRepo:  userRepo,
	}
}

// Handle executes the remove member command
func (h *RemoveMemberCommandHandler) Handle(cmd RemoveMemberCommand) error {
	// Validate user exists
	user, err := h.userRepo.GetUserByUserId(cmd.UserId)
	if err != nil || user == nil {
		return fmt.Errorf("user not found: %s", cmd.UserId)
	}

	// Business rule: Check if user's default group is being removed
	if user.DefaultGroupId == cmd.GroupId {
		user.ClearDefaultGroup()
		err = h.userRepo.UpdateUser(*user)
		if err != nil {
			// Log but continue - this is a side effect
			// In a full event-driven system, this would be handled by an event handler
		}
	}

	// Remove the member
	err = h.groupRepo.RemoveUserFromUserGroup(cmd.GroupId, user.Id)
	if err != nil {
		return fmt.Errorf("failed to remove member: %w", err)
	}

	return nil
}
