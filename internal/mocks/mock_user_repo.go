package mocks

import (
	taskDomain "github.com/happYness-Project/taskManagementGolang/internal/task/domain"
	userModel "github.com/happYness-Project/taskManagementGolang/internal/user/domain"
	"github.com/stretchr/testify/mock"
)

type MockUserRepo struct{ mock.Mock }

// CreateUser implements repository.UserRepository.
func (m *MockUserRepo) CreateUser(user userModel.User) error {
	args := m.Called(user)
	return args.Error(0)
}

// GetAllUsers implements repository.UserRepository.
func (m *MockUserRepo) GetAllUsers() ([]*userModel.User, error) {
	args := m.Called()
	return args.Get(0).([]*userModel.User), args.Error(1)
}

// GetUserByEmail implements repository.UserRepository.
func (m *MockUserRepo) GetUserByEmail(email string) (*userModel.User, error) {
	args := m.Called(email)
	return args.Get(0).(*userModel.User), args.Error(1)
}

// GetUserByUserId implements repository.UserRepository.
func (m *MockUserRepo) GetUserByUserId(userId string) (*userModel.User, error) {
	args := m.Called(userId)
	return args.Get(0).(*userModel.User), args.Error(1)
}

// GetUserByUsername implements repository.UserRepository.
func (m *MockUserRepo) GetUserByUsername(username string) (*userModel.User, error) {
	args := m.Called(username)
	return args.Get(0).(*userModel.User), args.Error(1)
}

// GetUsersByGroupId implements repository.UserRepository.
func (m *MockUserRepo) GetUsersByGroupId(groupId int) ([]*userModel.User, error) {
	args := m.Called(groupId)
	return args.Get(0).([]*userModel.User), args.Error(1)
}

// GetUsersByGroupIdWithRoles implements repository.UserRepository.
func (m *MockUserRepo) GetUsersByGroupIdWithRoles(groupId int) ([]*userModel.UserWithRole, error) {
	args := m.Called(groupId)
	return args.Get(0).([]*userModel.UserWithRole), args.Error(1)
}

// GetUserRoleInGroup implements repository.UserRepository.
func (m *MockUserRepo) GetUserRoleInGroup(userId string, groupId int) (string, error) {
	args := m.Called(userId, groupId)
	return args.Get(0).(string), args.Error(1)
}

// UpdateDefaultGroupId implements repository.UserRepository.
func (m *MockUserRepo) UpdateDefaultGroupId(Id int, groupId int) error {
	args := m.Called(Id, groupId)
	return args.Error(0)
}

// UpdateUser implements repository.UserRepository.
func (m *MockUserRepo) UpdateUser(user userModel.User) error {
	args := m.Called(user)
	return args.Error(0)
}

type TaskRepo struct{ mock.Mock }

func (m *TaskRepo) GetAllTasks() ([]taskDomain.Task, error) {
	args := m.Called()
	return args.Get(0).([]taskDomain.Task), args.Error(1)
}

func (m *TaskRepo) GetTaskById(id string) (*taskDomain.Task, error) {
	args := m.Called(id)
	return args.Get(0).(*taskDomain.Task), args.Error(1)
}
