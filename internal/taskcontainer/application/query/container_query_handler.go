package query

import (
	"fmt"

	"github.com/happYness-Project/taskManagementGolang/internal/taskcontainer/domain"
	"github.com/happYness-Project/taskManagementGolang/internal/taskcontainer/repository"
)

// Queries define the read operations
type GetAllContainersQuery struct{}

type GetContainerByIdQuery struct {
	ContainerId string
}

type GetContainersByGroupIdQuery struct {
	GroupId int
}

// ContainerQueryHandler handles all read operations for TaskContainer
type ContainerQueryHandler struct {
	containerRepo repository.ContainerRepository
}

func NewContainerQueryHandler(
	containerRepo repository.ContainerRepository,
) *ContainerQueryHandler {
	return &ContainerQueryHandler{
		containerRepo: containerRepo,
	}
}

// HandleGetAllContainers retrieves all containers
func (h *ContainerQueryHandler) HandleGetAllContainers(query GetAllContainersQuery) ([]*domain.TaskContainer, error) {
	containers, err := h.containerRepo.AllTaskContainers()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve containers: %w", err)
	}
	return containers, nil
}

// HandleGetContainerById retrieves a single container by ID
func (h *ContainerQueryHandler) HandleGetContainerById(query GetContainerByIdQuery) (*domain.TaskContainer, error) {
	container, err := h.containerRepo.GetById(query.ContainerId)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve container: %w", err)
	}
	return container, nil
}

// HandleGetContainersByGroupId retrieves containers for a specific group
func (h *ContainerQueryHandler) HandleGetContainersByGroupId(query GetContainersByGroupIdQuery) ([]domain.TaskContainer, error) {
	containers, err := h.containerRepo.GetContainersByGroupId(query.GroupId)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve containers by group: %w", err)
	}
	return containers, nil
}
