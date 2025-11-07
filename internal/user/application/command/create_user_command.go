package command

import (
	"fmt"

	"github.com/happYness-Project/taskManagementGolang/internal/user/model"
	"github.com/happYness-Project/taskManagementGolang/internal/user/repository"
)

// CreateUserCommand represents the command to create a new user
type CreateUserCommand struct {
	UserId    string // UUID from JWT
	UserName  string
	FirstName string
	LastName  string
	Email     string
}

// CreateUserCommandHandler handles the creation of a new user
type CreateUserCommandHandler struct {
	userRepo repository.UserRepository
}

func NewCreateUserCommandHandler(
	userRepo repository.UserRepository,
) *CreateUserCommandHandler {
	return &CreateUserCommandHandler{
		userRepo: userRepo,
	}
}

// Handle executes the create user command
func (h *CreateUserCommandHandler) Handle(cmd CreateUserCommand) error {
	// Create the domain model (includes validation)
	user := model.NewUser(cmd.UserId, cmd.UserName, cmd.FirstName, cmd.LastName, cmd.Email)

	// Persist the user
	err := h.userRepo.CreateUser(*user)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}
