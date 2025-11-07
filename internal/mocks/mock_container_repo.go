package mocks

import (
	containerDomain "github.com/happYness-Project/taskManagementGolang/internal/taskcontainer/domain"
	"github.com/stretchr/testify/mock"
)

type MockContainerRepo struct{ mock.Mock }

// AllTaskContainers implements repository.ContainerRepository.
func (m *MockContainerRepo) AllTaskContainers() ([]*containerDomain.TaskContainer, error) {
	args := m.Called()
	return args.Get(0).([]*containerDomain.TaskContainer), args.Error(1)
}

// GetById implements repository.ContainerRepository.
func (m *MockContainerRepo) GetById(id string) (*containerDomain.TaskContainer, error) {
	args := m.Called(id)
	return args.Get(0).(*containerDomain.TaskContainer), args.Error(1)
}

// GetContainersByGroupId implements repository.ContainerRepository.
func (m *MockContainerRepo) GetContainersByGroupId(groupId int) ([]containerDomain.TaskContainer, error) {
	args := m.Called(groupId)
	return args.Get(0).([]containerDomain.TaskContainer), args.Error(1)
}

// CreateContainer implements repository.ContainerRepository.
func (m *MockContainerRepo) CreateContainer(container containerDomain.TaskContainer) error {
	args := m.Called(container)
	return args.Error(0)
}

// DeleteContainer implements repository.ContainerRepository.
func (m *MockContainerRepo) DeleteContainer(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

// RemoveContainerByUsergroupId implements repository.ContainerRepository.
func (m *MockContainerRepo) RemoveContainerByUsergroupId(groupId int) error {
	args := m.Called(groupId)
	return args.Error(0)
}
