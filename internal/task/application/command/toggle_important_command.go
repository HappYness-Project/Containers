package command

import (
	"fmt"

	"github.com/happYness-Project/taskManagementGolang/internal/task/repository"
)

type ToggleImportantCommand struct {
	TaskId      string
	IsImportant bool
}

type ToggleImportantCommandHandler struct {
	taskRepo repository.TaskRepository
}

func NewToggleImportantCommandHandler(taskRepo repository.TaskRepository) *ToggleImportantCommandHandler {
	return &ToggleImportantCommandHandler{taskRepo: taskRepo}
}

func (h *ToggleImportantCommandHandler) Handle(cmd ToggleImportantCommand) error {
	// Get existing task
	task, err := h.taskRepo.GetTaskById(cmd.TaskId)
	if err != nil || task == nil {
		return fmt.Errorf("task not found: %w", err)
	}

	// Toggle important using domain method
	task.ToggleImportant(cmd.IsImportant)

	// Persist changes
	err = h.taskRepo.UpdateImportantTask(cmd.TaskId, cmd.IsImportant)
	if err != nil {
		return fmt.Errorf("failed to toggle important: %w", err)
	}

	return nil
}
