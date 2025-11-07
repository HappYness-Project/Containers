package application

import (
	"fmt"

	qry "github.com/happYness-Project/taskManagementGolang/internal/user/application/query"
	"github.com/happYness-Project/taskManagementGolang/internal/user/repository"
)

// QueryBus routes queries to their handlers
type QueryBus struct {
	queryHandler *qry.UserQueryHandler
}

// NewQueryBus creates a new query bus with all handlers registered
func NewQueryBus(
	userRepo repository.UserRepository,
) *QueryBus {
	return &QueryBus{
		queryHandler: qry.NewUserQueryHandler(userRepo),
	}
}

// Execute dispatches the query to the appropriate handler
func (bus *QueryBus) Execute(query interface{}) (interface{}, error) {
	switch q := query.(type) {
	case qry.GetAllUsersQuery:
		return bus.queryHandler.HandleGetAllUsers(q)
	case qry.GetUserByIdQuery:
		return bus.queryHandler.HandleGetUserById(q)
	case qry.GetUserByEmailQuery:
		return bus.queryHandler.HandleGetUserByEmail(q)
	case qry.GetUserByUsernameQuery:
		return bus.queryHandler.HandleGetUserByUsername(q)
	case qry.GetUsersByGroupIdQuery:
		return bus.queryHandler.HandleGetUsersByGroupId(q)
	default:
		return nil, fmt.Errorf("unknown query type: %T", query)
	}
}
