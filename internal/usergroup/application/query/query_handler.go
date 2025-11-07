package query

import (
	"fmt"

	"github.com/happYness-Project/taskManagementGolang/internal/user/repository"
	"github.com/happYness-Project/taskManagementGolang/internal/usergroup/domain"
	userGroupRepo "github.com/happYness-Project/taskManagementGolang/internal/usergroup/repository"
)

// Queries define the read operations
type GetAllGroupsQuery struct{}

type GetGroupByIdQuery struct {
	GroupId int
}

type GetGroupsByUserIdQuery struct {
	UserId string // UUID
}

type GetGroupMembersQuery struct {
	GroupId int
}

// UserGroupQueryHandler handles all read operations for UserGroup
type UserGroupQueryHandler struct {
	groupRepo userGroupRepo.UserGroupRepository
	userRepo  repository.UserRepository
}

func NewUserGroupQueryHandler(
	groupRepo userGroupRepo.UserGroupRepository,
	userRepo repository.UserRepository,
) *UserGroupQueryHandler {
	return &UserGroupQueryHandler{
		groupRepo: groupRepo,
		userRepo:  userRepo,
	}
}

// HandleGetAllGroups retrieves all groups
func (h *UserGroupQueryHandler) HandleGetAllGroups(query GetAllGroupsQuery) ([]*domain.UserGroup, error) {
	groups, err := h.groupRepo.GetAllUsergroups()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve groups: %w", err)
	}
	return groups, nil
}

// HandleGetGroupById retrieves a single group by ID
func (h *UserGroupQueryHandler) HandleGetGroupById(query GetGroupByIdQuery) (*domain.UserGroup, error) {
	group, err := h.groupRepo.GetById(query.GroupId)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve group: %w", err)
	}

	if group.GroupId == 0 {
		return nil, fmt.Errorf("group not found: %d", query.GroupId)
	}

	return group, nil
}

// HandleGetGroupsByUserId retrieves all groups for a specific user
func (h *UserGroupQueryHandler) HandleGetGroupsByUserId(query GetGroupsByUserIdQuery) ([]*domain.UserGroup, error) {
	// First validate user exists
	user, err := h.userRepo.GetUserByUserId(query.UserId)
	if err != nil || user == nil {
		return nil, fmt.Errorf("user not found: %s", query.UserId)
	}

	groups, err := h.groupRepo.GetUserGroupsByUserId(user.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve user groups: %w", err)
	}

	return groups, nil
}
