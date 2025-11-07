package application

import (
	"fmt"

	cmd "github.com/happYness-Project/taskManagementGolang/internal/taskcontainer/application/command"
	"github.com/happYness-Project/taskManagementGolang/internal/taskcontainer/repository"
)

// CommandBus routes commands to their handlers
type CommandBus struct {
	createContainerHandler *cmd.CreateContainerCommandHandler
	deleteContainerHandler *cmd.DeleteContainerCommandHandler
}

// NewCommandBus creates a new command bus with all handlers registered
func NewCommandBus(
	containerRepo repository.ContainerRepository,
) *CommandBus {
	return &CommandBus{
		createContainerHandler: cmd.NewCreateContainerCommandHandler(containerRepo),
		deleteContainerHandler: cmd.NewDeleteContainerCommandHandler(containerRepo),
	}
}

// Execute dispatches the command to the appropriate handler
func (bus *CommandBus) Execute(command interface{}) (interface{}, error) {
	switch c := command.(type) {
	case cmd.CreateContainerCommand:
		return bus.createContainerHandler.Handle(c)
	case cmd.DeleteContainerCommand:
		return nil, bus.deleteContainerHandler.Handle(c)
	default:
		return nil, fmt.Errorf("unknown command type: %T", command)
	}
}
