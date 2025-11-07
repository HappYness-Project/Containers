package command

import (
	"fmt"

	"github.com/happYness-Project/taskManagementGolang/internal/user/repository"
	userGroupRepo "github.com/happYness-Project/taskManagementGolang/internal/usergroup/repository"
)

// AddMemberCommand represents the command to add a member to a group
type AddMemberCommand struct {
	GroupId int
	UserId  string // UUID
}

// AddMemberCommandHandler handles adding a member to a group
type AddMemberCommandHandler struct {
	groupRepo userGroupRepo.UserGroupRepository
	userRepo  repository.UserRepository
}

func NewAddMemberCommandHandler(
	groupRepo userGroupRepo.UserGroupRepository,
	userRepo repository.UserRepository,
) *AddMemberCommandHandler {
	return &AddMemberCommandHandler{
		groupRepo: groupRepo,
		userRepo:  userRepo,
	}
}

// Handle executes the add member command
func (h *AddMemberCommandHandler) Handle(cmd AddMemberCommand) error {
	// Validate user exists
	user, err := h.userRepo.GetUserByUserId(cmd.UserId)
	if err != nil || user == nil {
		return fmt.Errorf("user not found: %s", cmd.UserId)
	}

	// Validate group exists
	group, err := h.groupRepo.GetById(cmd.GroupId)
	if err != nil || group.GroupId == 0 {
		return fmt.Errorf("group not found: %d", cmd.GroupId)
	}

	// Add member with default role (member)
	err = h.groupRepo.InsertUserGroupUserTable(cmd.GroupId, user.Id)
	if err != nil {
		return fmt.Errorf("failed to add member: %w", err)
	}

	return nil
}
