package query

import (
	"fmt"

	"github.com/happYness-Project/taskManagementGolang/internal/task/domain"
	"github.com/happYness-Project/taskManagementGolang/internal/task/repository"
)

// Queries define the read operations
type GetAllTasksQuery struct{}

type GetTaskByIdQuery struct {
	TaskId string
}

type GetTasksByContainerIdQuery struct {
	ContainerId string
}

type GetAllTasksByGroupIdQuery struct {
	GroupId       int
	OnlyImportant bool
}

// TaskQueryHandler handles all read operations for Task
type TaskQueryHandler struct {
	taskRepo repository.TaskRepository
}

func NewTaskQueryHandler(taskRepo repository.TaskRepository) *TaskQueryHandler {
	return &TaskQueryHandler{taskRepo: taskRepo}
}

// HandleGetAllTasks retrieves all tasks
func (h *TaskQueryHandler) HandleGetAllTasks(query GetAllTasksQuery) ([]domain.Task, error) {
	tasks, err := h.taskRepo.GetAllTasks()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve tasks: %w", err)
	}
	return tasks, nil
}

// HandleGetTaskById retrieves a single task by ID
func (h *TaskQueryHandler) HandleGetTaskById(query GetTaskByIdQuery) (*domain.Task, error) {
	task, err := h.taskRepo.GetTaskById(query.TaskId)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve task: %w", err)
	}
	return task, nil
}

// HandleGetTasksByContainerId retrieves tasks for a specific container
func (h *TaskQueryHandler) HandleGetTasksByContainerId(query GetTasksByContainerIdQuery) ([]domain.Task, error) {
	tasks, err := h.taskRepo.GetTasksByContainerId(query.ContainerId)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve tasks by container: %w", err)
	}
	return tasks, nil
}

// HandleGetAllTasksByGroupId retrieves tasks for a specific group
func (h *TaskQueryHandler) HandleGetAllTasksByGroupId(query GetAllTasksByGroupIdQuery) ([]domain.Task, error) {
	var tasks []domain.Task
	var err error

	if query.OnlyImportant {
		tasks, err = h.taskRepo.GetAllTasksByGroupIdOnlyImportant(query.GroupId)
	} else {
		tasks, err = h.taskRepo.GetAllTasksByGroupId(query.GroupId)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to retrieve tasks by group: %w", err)
	}
	return tasks, nil
}
