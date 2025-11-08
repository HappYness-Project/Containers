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

	// Update fields
	task.TaskName = cmd.TaskName
	task.TaskDesc = cmd.TaskDesc
	task.TargetDate = cmd.TargetDate
	task.Priority = cmd.Priority
	task.Category = cmd.Category
	task.UpdatedAt = time.Now().UTC()

	// Persist changes
	err = h.taskRepo.UpdateTask(*task)
	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	return nil
}
