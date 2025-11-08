package command

import (
	"fmt"
	"time"

	"github.com/happYness-Project/taskManagementGolang/internal/task/repository"
)

type UpdateTaskCommand struct {
	TaskId     string
	TaskName   string
	TaskDesc   string
	TargetDate time.Time
	Priority   string
	Category   string
}

type UpdateTaskCommandHandler struct {
	taskRepo repository.TaskRepository
}

func NewUpdateTaskCommandHandler(taskRepo repository.TaskRepository) *UpdateTaskCommandHandler {
	return &UpdateTaskCommandHandler{taskRepo: taskRepo}
}

func (h *UpdateTaskCommandHandler) Handle(cmd UpdateTaskCommand) error {
	// Get existing task
	task, err := h.taskRepo.GetTaskById(cmd.TaskId)
	if err != nil || task == nil {
		return fmt.Errorf("task not found: %w", err)
	}

	// Update task using domain method (enforces validation)
	err = task.UpdateTask(cmd.TaskName, cmd.TaskDesc, cmd.TargetDate, cmd.Priority, cmd.Category)
	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	// Persist changes
	err = h.taskRepo.UpdateTask(*task)
	if err != nil {
		return fmt.Errorf("failed to persist task update: %w", err)
	}

	return nil
}
