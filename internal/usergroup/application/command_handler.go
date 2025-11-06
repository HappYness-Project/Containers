package application

import (
	"fmt"

	"github.com/happYness-Project/taskManagementGolang/internal/user/repository"
	"github.com/happYness-Project/taskManagementGolang/internal/usergroup/domain"
	userGroupRepo "github.com/happYness-Project/taskManagementGolang/internal/usergroup/repository"
)

// Commands define the write operations
type CreateGroupCommand struct {
	GroupName string
	GroupDesc string
	GroupType string
	CreatorId string // UUID from JWT
}

type AddMemberCommand struct {
	GroupId int
	UserId  string // UUID
}

type RemoveMemberCommand struct {
	GroupId int
	UserId  string // UUID
}

type ChangeMemberRoleCommand struct {
	GroupId int
	UserId  string // UUID
	NewRole string
}

type DeleteGroupCommand struct {
	GroupId int
}

// UserGroupCommandHandler handles all write operations for UserGroup
type UserGroupCommandHandler struct {
	groupRepo userGroupRepo.UserGroupRepository
	userRepo  repository.UserRepository
}

func NewUserGroupCommandHandler(
	groupRepo userGroupRepo.UserGroupRepository,
	userRepo repository.UserRepository,
) *UserGroupCommandHandler {
	return &UserGroupCommandHandler{
		groupRepo: groupRepo,
		userRepo:  userRepo,
	}
}

// HandleCreateGroup handles the creation of a new group
func (h *UserGroupCommandHandler) HandleCreateGroup(cmd CreateGroupCommand) (int, error) {
	// Validate creator exists
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

// HandleAddMember handles adding a member to a group
func (h *UserGroupCommandHandler) HandleAddMember(cmd AddMemberCommand) error {
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

// HandleRemoveMember handles removing a member from a group
func (h *UserGroupCommandHandler) HandleRemoveMember(cmd RemoveMemberCommand) error {
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

// HandleChangeMemberRole handles changing a member's role in a group
func (h *UserGroupCommandHandler) HandleChangeMemberRole(cmd ChangeMemberRoleCommand) error {
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

// HandleDeleteGroup handles deleting a group
func (h *UserGroupCommandHandler) HandleDeleteGroup(cmd DeleteGroupCommand) error {
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
