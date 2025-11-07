package command

import (
	"fmt"

	"github.com/happYness-Project/taskManagementGolang/internal/user/repository"
	"github.com/happYness-Project/taskManagementGolang/internal/usergroup/domain"
	userGroupRepo "github.com/happYness-Project/taskManagementGolang/internal/usergroup/repository"
)

// CreateGroupCommand represents the command to create a new group
type CreateGroupCommand struct {
	GroupName string
	GroupDesc string
	GroupType string
	CreatorId string // UUID from JWT
}

// CreateGroupCommandHandler handles the creation of a new group
type CreateGroupCommandHandler struct {
	groupRepo userGroupRepo.UserGroupRepository
	userRepo  repository.UserRepository
}

func NewCreateGroupCommandHandler(
	groupRepo userGroupRepo.UserGroupRepository,
	userRepo repository.UserRepository,
) *CreateGroupCommandHandler {
	return &CreateGroupCommandHandler{
		groupRepo: groupRepo,
		userRepo:  userRepo,
	}
}

// Handle executes the create group command
func (h *CreateGroupCommandHandler) Handle(cmd CreateGroupCommand) (int, error) {
	creator, err := h.userRepo.GetUserByUserId(cmd.CreatorId)
	if err != nil || creator == nil {
		return 0, fmt.Errorf("creator user not found: %s", cmd.CreatorId)
	}

	// Create the domain model
	group, err := domain.NewUserGroup(cmd.GroupName, cmd.GroupDesc, cmd.GroupType)
	if err != nil {
		return 0, fmt.Errorf("invalid group data: %w", err)
	}

	// Persist the group with creator as admin (transactional)
	groupId, err := h.groupRepo.CreateGroupWithUsers(*group, creator.Id)
	if err != nil {
		return 0, fmt.Errorf("failed to create group: %w", err)
	}

	return groupId, nil
}
