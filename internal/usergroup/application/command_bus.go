package application

import (
	"fmt"

	"github.com/happYness-Project/taskManagementGolang/internal/user/repository"
	cmd "github.com/happYness-Project/taskManagementGolang/internal/usergroup/application/command"
	userGroupRepo "github.com/happYness-Project/taskManagementGolang/internal/usergroup/repository"
)

// CommandBus routes commands to their handlers
type CommandBus struct {
	createGroupHandler      *cmd.CreateGroupCommandHandler
	addMemberHandler        *cmd.AddMemberCommandHandler
	removeMemberHandler     *cmd.RemoveMemberCommandHandler
	changeMemberRoleHandler *cmd.ChangeMemberRoleCommandHandler
	deleteGroupHandler      *cmd.DeleteGroupCommandHandler
}

// NewCommandBus creates a new command bus with all handlers registered
func NewCommandBus(
	groupRepo userGroupRepo.UserGroupRepository,
	userRepo repository.UserRepository,
) *CommandBus {
	return &CommandBus{
		createGroupHandler:      cmd.NewCreateGroupCommandHandler(groupRepo, userRepo),
		addMemberHandler:        cmd.NewAddMemberCommandHandler(groupRepo, userRepo),
		removeMemberHandler:     cmd.NewRemoveMemberCommandHandler(groupRepo, userRepo),
		changeMemberRoleHandler: cmd.NewChangeMemberRoleCommandHandler(groupRepo, userRepo),
		deleteGroupHandler:      cmd.NewDeleteGroupCommandHandler(groupRepo),
	}
}

// Execute dispatches the command to the appropriate handler
func (bus *CommandBus) Execute(command interface{}) (interface{}, error) {
	switch c := command.(type) {
	case cmd.CreateGroupCommand:
		return bus.createGroupHandler.Handle(c)
	case cmd.AddMemberCommand:
		return nil, bus.addMemberHandler.Handle(c)
	case cmd.RemoveMemberCommand:
		return nil, bus.removeMemberHandler.Handle(c)
	case cmd.ChangeMemberRoleCommand:
		return nil, bus.changeMemberRoleHandler.Handle(c)
	case cmd.DeleteGroupCommand:
		return nil, bus.deleteGroupHandler.Handle(c)
	default:
		return nil, fmt.Errorf("unknown command type: %T", command)
	}
}
