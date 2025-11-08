package command

import (
	"fmt"

	"github.com/happYness-Project/taskManagementGolang/internal/task/repository"
)

type ToggleCompletionCommand struct {
	TaskId      string
	IsCompleted bool
}

type ToggleCompletionCommandHandler struct {
	taskRepo repository.TaskRepository
}

func NewToggleCompletionCommandHandler(taskRepo repository.TaskRepository) *ToggleCompletionCommandHandler {
	return &ToggleCompletionCommandHandler{taskRepo: taskRepo}
}

func (h *ToggleCompletionCommandHandler) Handle(cmd ToggleCompletionCommand) error {
	// Verify task exists
	task, err := h.taskRepo.GetTaskById(cmd.TaskId)
	if err != nil || task == nil {
		return fmt.Errorf("task not found: %w", err)
	}

	// Toggle completion status
	err = h.taskRepo.DoneTask(cmd.TaskId, cmd.IsCompleted)
	if err != nil {
		return fmt.Errorf("failed to toggle completion: %w", err)
	}

	return nil
}
