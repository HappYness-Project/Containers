package command

import (
	"fmt"
	"time"

	"github.com/happYness-Project/taskManagementGolang/internal/task/domain"
	"github.com/happYness-Project/taskManagementGolang/internal/task/repository"
)

type CreateTaskCommand struct {
	ContainerId string
	TaskName    string
	TaskDesc    string
	TargetDate  time.Time
	Priority    string
	Category    string
}

type CreateTaskCommandHandler struct {
	taskRepo repository.TaskRepository
}

func NewCreateTaskCommandHandler(taskRepo repository.TaskRepository) *CreateTaskCommandHandler {
	return &CreateTaskCommandHandler{taskRepo: taskRepo}
}

func (h *CreateTaskCommandHandler) Handle(cmd CreateTaskCommand) (domain.Task, error) {
	// Create task using domain logic
	task, err := domain.CreateTask(
		cmd.TaskName,
		cmd.TaskDesc,
		cmd.TargetDate,
		cmd.Priority,
		cmd.Category,
	)
	if err != nil {
		return domain.Task{}, fmt.Errorf("failed to create task: %w", err)
	}

	// Persist task
	newTask, err := h.taskRepo.CreateTask(cmd.ContainerId, *task)
	if err != nil {
		return domain.Task{}, fmt.Errorf("failed to persist task: %w", err)
	}

	return newTask, nil
}
