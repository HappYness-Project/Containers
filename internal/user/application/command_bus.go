package application

import (
	"fmt"

	cmd "github.com/happYness-Project/taskManagementGolang/internal/user/application/command"
	"github.com/happYness-Project/taskManagementGolang/internal/user/repository"
)

// CommandBus routes commands to their handlers
type CommandBus struct {
	createUserHandler        *cmd.CreateUserCommandHandler
	updateUserHandler        *cmd.UpdateUserCommandHandler
	updateDefaultGroupHandler *cmd.UpdateDefaultGroupCommandHandler
}

// NewCommandBus creates a new command bus with all handlers registered
func NewCommandBus(
	userRepo repository.UserRepository,
) *CommandBus {
	return &CommandBus{
		createUserHandler:        cmd.NewCreateUserCommandHandler(userRepo),
		updateUserHandler:        cmd.NewUpdateUserCommandHandler(userRepo),
		updateDefaultGroupHandler: cmd.NewUpdateDefaultGroupCommandHandler(userRepo),
	}
}

// Execute dispatches the command to the appropriate handler
func (bus *CommandBus) Execute(command interface{}) (interface{}, error) {
	switch c := command.(type) {
	case cmd.CreateUserCommand:
		return nil, bus.createUserHandler.Handle(c)
	case cmd.UpdateUserCommand:
		return nil, bus.updateUserHandler.Handle(c)
	case cmd.UpdateDefaultGroupCommand:
		return nil, bus.updateDefaultGroupHandler.Handle(c)
	default:
		return nil, fmt.Errorf("unknown command type: %T", command)
	}
}
