package mocks

import (
	"github.com/happYness-Project/taskManagementGolang/internal/usergroup/domain"
	"github.com/stretchr/testify/mock"
)

type MockUserGroupRepo struct{ mock.Mock }

// DeleteUserGroup implements repository.UserGroupRepository.
func (m *MockUserGroupRepo) DeleteUserGroup(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

// GetAllUsergroups implements repository.UserGroupRepository.
func (m *MockUserGroupRepo) GetAllUsergroups() ([]*domain.UserGroup, error) {
	args := m.Called()
	return args.Get(0).([]*domain.UserGroup), args.Error(1)
}

// GetById implements repository.UserGroupRepository.
func (m *MockUserGroupRepo) GetById(id int) (*domain.UserGroup, error) {
	args := m.Called(id)
	return args.Get(0).(*domain.UserGroup), args.Error(1)
}

// GetUserGroupsByUserId implements repository.UserGroupRepository.
func (m *MockUserGroupRepo) GetUserGroupsByUserId(userId int) ([]*domain.UserGroup, error) {
	args := m.Called(userId)
	return args.Get(0).([]*domain.UserGroup), args.Error(1)
}

// InsertUserGroupUserTable implements repository.UserGroupRepository.
func (m *MockUserGroupRepo) InsertUserGroupUserTable(groupId int, userId int) error {
	args := m.Called(groupId, userId)
	return args.Error(0)
}

func (m *MockUserGroupRepo) RemoveUserFromUserGroup(groupId int, userId int) error {
	args := m.Called(groupId, userId)
	return args.Error(0)
}
func (m *MockUserGroupRepo) CreateGroupWithUsers(ug domain.UserGroup, userId int) (int, error) {
	args := m.Called(ug, userId)
	return args.Get(0).(int), args.Error(0)
}

// UpdateUserRoleInGroup implements repository.UserGroupRepository.
func (m *MockUserGroupRepo) UpdateUserRoleInGroup(groupId int, userId int, role string) error {
	args := m.Called(groupId, userId, role)
	return args.Error(0)
}
