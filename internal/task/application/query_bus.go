package application

import (
	"fmt"

	qry "github.com/happYness-Project/taskManagementGolang/internal/task/application/query"
	"github.com/happYness-Project/taskManagementGolang/internal/task/repository"
)

// QueryBus routes queries to their handlers
type QueryBus struct {
	queryHandler *qry.TaskQueryHandler
}

// NewQueryBus creates a new query bus with all handlers registered
func NewQueryBus(taskRepo repository.TaskRepository) *QueryBus {
	return &QueryBus{
		queryHandler: qry.NewTaskQueryHandler(taskRepo),
	}
}

// Execute dispatches the query to the appropriate handler
func (bus *QueryBus) Execute(query interface{}) (interface{}, error) {
	switch q := query.(type) {
	case qry.GetAllTasksQuery:
		return bus.queryHandler.HandleGetAllTasks(q)
	case qry.GetTaskByIdQuery:
		return bus.queryHandler.HandleGetTaskById(q)
	case qry.GetTasksByContainerIdQuery:
		return bus.queryHandler.HandleGetTasksByContainerId(q)
	case qry.GetAllTasksByGroupIdQuery:
		return bus.queryHandler.HandleGetAllTasksByGroupId(q)
	default:
		return nil, fmt.Errorf("unknown query type: %T", query)
	}
}
