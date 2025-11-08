package application

import (
	"fmt"

	cmd "github.com/happYness-Project/taskManagementGolang/internal/task/application/command"
	"github.com/happYness-Project/taskManagementGolang/internal/task/repository"
)

// CommandBus routes commands to their handlers
type CommandBus struct {
	createTaskHandler         *cmd.CreateTaskCommandHandler
	updateTaskHandler         *cmd.UpdateTaskCommandHandler
	deleteTaskHandler         *cmd.DeleteTaskCommandHandler
	toggleCompletionHandler   *cmd.ToggleCompletionCommandHandler
	toggleImportantHandler    *cmd.ToggleImportantCommandHandler
}

// NewCommandBus creates a new command bus with all handlers registered
func NewCommandBus(taskRepo repository.TaskRepository) *CommandBus {
	return &CommandBus{
		createTaskHandler:       cmd.NewCreateTaskCommandHandler(taskRepo),
		updateTaskHandler:       cmd.NewUpdateTaskCommandHandler(taskRepo),
		deleteTaskHandler:       cmd.NewDeleteTaskCommandHandler(taskRepo),
		toggleCompletionHandler: cmd.NewToggleCompletionCommandHandler(taskRepo),
		toggleImportantHandler:  cmd.NewToggleImportantCommandHandler(taskRepo),
	}
}

// Execute dispatches the command to the appropriate handler
func (bus *CommandBus) Execute(command interface{}) (interface{}, error) {
	switch c := command.(type) {
	case cmd.CreateTaskCommand:
		return bus.createTaskHandler.Handle(c)
	case cmd.UpdateTaskCommand:
		return nil, bus.updateTaskHandler.Handle(c)
	case cmd.DeleteTaskCommand:
		return nil, bus.deleteTaskHandler.Handle(c)
	case cmd.ToggleCompletionCommand:
		return nil, bus.toggleCompletionHandler.Handle(c)
	case cmd.ToggleImportantCommand:
		return nil, bus.toggleImportantHandler.Handle(c)
	default:
		return nil, fmt.Errorf("unknown command type: %T", command)
	}
}
