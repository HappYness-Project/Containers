package application

import (
	"fmt"

	qry "github.com/happYness-Project/taskManagementGolang/internal/taskcontainer/application/query"
	"github.com/happYness-Project/taskManagementGolang/internal/taskcontainer/repository"
)

// QueryBus routes queries to their handlers
type QueryBus struct {
	queryHandler *qry.ContainerQueryHandler
}

// NewQueryBus creates a new query bus with all handlers registered
func NewQueryBus(
	containerRepo repository.ContainerRepository,
) *QueryBus {
	return &QueryBus{
		queryHandler: qry.NewContainerQueryHandler(containerRepo),
	}
}

// Execute dispatches the query to the appropriate handler
func (bus *QueryBus) Execute(query interface{}) (interface{}, error) {
	switch q := query.(type) {
	case qry.GetAllContainersQuery:
		return bus.queryHandler.HandleGetAllContainers(q)
	case qry.GetContainerByIdQuery:
		return bus.queryHandler.HandleGetContainerById(q)
	case qry.GetContainersByGroupIdQuery:
		return bus.queryHandler.HandleGetContainersByGroupId(q)
	default:
		return nil, fmt.Errorf("unknown query type: %T", query)
	}
}
