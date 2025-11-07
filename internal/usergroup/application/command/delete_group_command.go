package command

import (
	"fmt"

	userGroupRepo "github.com/happYness-Project/taskManagementGolang/internal/usergroup/repository"
)

// DeleteGroupCommand represents the command to delete a group
type DeleteGroupCommand struct {
	GroupId int
}

// DeleteGroupCommandHandler handles deleting a group
type DeleteGroupCommandHandler struct {
	groupRepo userGroupRepo.UserGroupRepository
}

func NewDeleteGroupCommandHandler(
	groupRepo userGroupRepo.UserGroupRepository,
) *DeleteGroupCommandHandler {
	return &DeleteGroupCommandHandler{
		groupRepo: groupRepo,
	}
}

// Handle executes the delete group command
func (h *DeleteGroupCommandHandler) Handle(cmd DeleteGroupCommand) error {
	// Validate group exists
	group, err := h.groupRepo.GetById(cmd.GroupId)
	if err != nil || group.GroupId == 0 {
		return fmt.Errorf("group not found: %d", cmd.GroupId)
	}

	// Delete the group (cascade will handle related records)
	err = h.groupRepo.DeleteUserGroup(cmd.GroupId)
	if err != nil {
		return fmt.Errorf("failed to delete group: %w", err)
	}

	return nil
}
