package builders

import (
	"github.com/google/uuid"
	"github.com/happYness-Project/taskManagementGolang/internal/taskcontainer/domain"
)

// TaskContainerBuilder provides a fluent interface for building test task containers
type TaskContainerBuilder struct {
	container *domain.TaskContainer
}

// NewTaskContainerBuilder creates a new task container builder with sensible defaults
func NewTaskContainerBuilder() *TaskContainerBuilder {
	return &TaskContainerBuilder{
		container: &domain.TaskContainer{
			Id:             uuid.New().String(),
			Name:           "Test Container",
			Description:    "Test Description",
			Type:           "default",
			IsActive:       true,
			Activity_level: 0,
			UsergroupId:    0, // Must be set
		},
	}
}

// WithId sets a specific container ID
func (b *TaskContainerBuilder) WithId(id string) *TaskContainerBuilder {
	b.container.Id = id
	return b
}

// WithName sets the container name
func (b *TaskContainerBuilder) WithName(name string) *TaskContainerBuilder {
	b.container.Name = name
	return b
}

// WithDescription sets the container description
func (b *TaskContainerBuilder) WithDescription(desc string) *TaskContainerBuilder {
	b.container.Description = desc
	return b
}

// WithType sets the container type
func (b *TaskContainerBuilder) WithType(containerType string) *TaskContainerBuilder {
	b.container.Type = containerType
	return b
}

// WithUsergroupId sets the usergroup ID
func (b *TaskContainerBuilder) WithUsergroupId(groupId int) *TaskContainerBuilder {
	b.container.UsergroupId = groupId
	return b
}

// WithActivityLevel sets the activity level
func (b *TaskContainerBuilder) WithActivityLevel(level int) *TaskContainerBuilder {
	b.container.Activity_level = level
	return b
}

// Inactive marks the container as inactive
func (b *TaskContainerBuilder) Inactive() *TaskContainerBuilder {
	b.container.IsActive = false
	return b
}

// Build returns the built container
func (b *TaskContainerBuilder) Build() *domain.TaskContainer {
	return b.container
}
