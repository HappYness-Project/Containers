package application

import (
	"fmt"

	"github.com/happYness-Project/taskManagementGolang/internal/user/repository"
	qry "github.com/happYness-Project/taskManagementGolang/internal/usergroup/application/query"
	userGroupRepo "github.com/happYness-Project/taskManagementGolang/internal/usergroup/repository"
)

// QueryBus routes queries to their handlers
type QueryBus struct {
	queryHandler *qry.UserGroupQueryHandler
}

// NewQueryBus creates a new query bus with all handlers registered
func NewQueryBus(
	groupRepo userGroupRepo.UserGroupRepository,
	userRepo repository.UserRepository,
) *QueryBus {
	return &QueryBus{
		queryHandler: qry.NewUserGroupQueryHandler(groupRepo, userRepo),
	}
}

// Execute dispatches the query to the appropriate handler
func (bus *QueryBus) Execute(query interface{}) (interface{}, error) {
	switch q := query.(type) {
	case qry.GetAllGroupsQuery:
		return bus.queryHandler.HandleGetAllGroups(q)
	case qry.GetGroupByIdQuery:
		return bus.queryHandler.HandleGetGroupById(q)
	case qry.GetGroupsByUserIdQuery:
		return bus.queryHandler.HandleGetGroupsByUserId(q)
	default:
		return nil, fmt.Errorf("unknown query type: %T", query)
	}
}
