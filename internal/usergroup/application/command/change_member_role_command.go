package command

import (
	"fmt"

	"github.com/happYness-Project/taskManagementGolang/internal/user/repository"
	"github.com/happYness-Project/taskManagementGolang/internal/usergroup/domain"
	userGroupRepo "github.com/happYness-Project/taskManagementGolang/internal/usergroup/repository"
)

// ChangeMemberRoleCommand represents the command to change a member's role
type ChangeMemberRoleCommand struct {
	GroupId int
	UserId  string // UUID
	NewRole string
}

// ChangeMemberRoleCommandHandler handles changing a member's role in a group
type ChangeMemberRoleCommandHandler struct {
	groupRepo userGroupRepo.UserGroupRepository
	userRepo  repository.UserRepository
}

func NewChangeMemberRoleCommandHandler(
	groupRepo userGroupRepo.UserGroupRepository,
	userRepo repository.UserRepository,
) *ChangeMemberRoleCommandHandler {
	return &ChangeMemberRoleCommandHandler{
		groupRepo: groupRepo,
		userRepo:  userRepo,
	}
}

// Handle executes the change member role command
func (h *ChangeMemberRoleCommandHandler) Handle(cmd ChangeMemberRoleCommand) error {
	// Validate role
	role, err := domain.NewRole(cmd.NewRole)
	if err != nil {
		return err
	}

	// Validate user exists
	user, err := h.userRepo.GetUserByUserId(cmd.UserId)
	if err != nil || user == nil {
		return fmt.Errorf("user not found: %s", cmd.UserId)
	}

	// Update the role
	err = h.groupRepo.UpdateUserRoleInGroup(cmd.GroupId, user.Id, role.String())
	if err != nil {
		return fmt.Errorf("failed to update role: %w", err)
	}

	return nil
}
