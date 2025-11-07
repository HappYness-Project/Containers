package command

import (
	"fmt"

	"github.com/happYness-Project/taskManagementGolang/internal/taskcontainer/repository"
)

// DeleteContainerCommand represents the command to delete a task container
type DeleteContainerCommand struct {
	ContainerId string
}

// DeleteContainerCommandHandler handles deleting a task container
type DeleteContainerCommandHandler struct {
	containerRepo repository.ContainerRepository
}

func NewDeleteContainerCommandHandler(
	containerRepo repository.ContainerRepository,
) *DeleteContainerCommandHandler {
	return &DeleteContainerCommandHandler{
		containerRepo: containerRepo,
	}
}

// Handle executes the delete container command
func (h *DeleteContainerCommandHandler) Handle(cmd DeleteContainerCommand) error {
	err := h.containerRepo.DeleteContainer(cmd.ContainerId)
	if err != nil {
		return fmt.Errorf("failed to delete container: %w", err)
	}

	return nil
}
