package command

import (
	"fmt"

	"github.com/happYness-Project/taskManagementGolang/internal/task/repository"
)

type DeleteTaskCommand struct {
	TaskId string
}

type DeleteTaskCommandHandler struct {
	taskRepo repository.TaskRepository
}

func NewDeleteTaskCommandHandler(taskRepo repository.TaskRepository) *DeleteTaskCommandHandler {
	return &DeleteTaskCommandHandler{taskRepo: taskRepo}
}

func (h *DeleteTaskCommandHandler) Handle(cmd DeleteTaskCommand) error {
	err := h.taskRepo.DeleteTask(cmd.TaskId)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	return nil
}
