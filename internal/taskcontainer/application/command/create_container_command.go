package command

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/happYness-Project/taskManagementGolang/internal/taskcontainer/domain"
	"github.com/happYness-Project/taskManagementGolang/internal/taskcontainer/repository"
)

// CreateContainerCommand represents the command to create a new task container
type CreateContainerCommand struct {
	Name        string
	Description string
	Type        string
	UserGroupId int
}

// CreateContainerCommandHandler handles creating a new task container
type CreateContainerCommandHandler struct {
	containerRepo repository.ContainerRepository
}

func NewCreateContainerCommandHandler(
	containerRepo repository.ContainerRepository,
) *CreateContainerCommandHandler {
	return &CreateContainerCommandHandler{
		containerRepo: containerRepo,
	}
}

// Handle executes the create container command
func (h *CreateContainerCommandHandler) Handle(cmd CreateContainerCommand) (string, error) {
	// Create domain model
	container := domain.TaskContainer{
		Id:             uuid.New().String(),
		Name:           cmd.Name,
		Description:    cmd.Description,
		Type:           cmd.Type,
		IsActive:       true,
		Activity_level: 0,
		UsergroupId:    cmd.UserGroupId,
	}

	// Persist
	err := h.containerRepo.CreateContainer(container)
	if err != nil {
		return "", fmt.Errorf("failed to create container: %w", err)
	}

	return container.Id, nil
}
