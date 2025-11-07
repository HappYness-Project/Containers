package route

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/happYness-Project/taskManagementGolang/internal/taskcontainer/application"
	"github.com/happYness-Project/taskManagementGolang/internal/taskcontainer/application/command"
	"github.com/happYness-Project/taskManagementGolang/internal/taskcontainer/application/query"
	container "github.com/happYness-Project/taskManagementGolang/internal/taskcontainer/repository"
	user "github.com/happYness-Project/taskManagementGolang/internal/user/repository"
	"github.com/happYness-Project/taskManagementGolang/pkg/constants"
	"github.com/happYness-Project/taskManagementGolang/pkg/loggers"
	"github.com/happYness-Project/taskManagementGolang/pkg/response"
)

type Handler struct {
	logger     *loggers.AppLogger
	commandBus *application.CommandBus
	queryBus   *application.QueryBus
}

func NewHandler(logger *loggers.AppLogger, repo container.ContainerRepository, userRepo user.UserRepository) *Handler {
	return &Handler{
		logger:     logger,
		commandBus: application.NewCommandBus(repo),
		queryBus:   application.NewQueryBus(repo),
	}
}
func (h *Handler) RegisterRoutes(router chi.Router) {
	router.Route("/api/task-containers", func(r chi.Router) {
		r.Post("/", h.handleCreateTaskContainer)
		r.Get("/", h.handleGetTaskContainers)
		r.Get("/{containerID}", h.handleGetTaskContainerById)
		r.Delete("/{containerID}", h.handleDeleteTaskContainer)
	})
	router.Get("/api/user-groups/{usergroupID}/task-containers", h.handleGetTaskContainersByGroupId)
}
func (h *Handler) handleGetTaskContainers(w http.ResponseWriter, r *http.Request) {
	// Use Query Bus
	result, err := h.queryBus.Execute(query.GetAllContainersQuery{})
	if err != nil {
		h.logger.Error().Err(err).Str("ErrorCode", TaskContainerGetError).Msg(err.Error())
		response.InternalServerError(w, "Error occurred during getting all task containers.")
		return
	}
	response.WriteJsonWithEncode(w, http.StatusOK, result)
}
func (h *Handler) handleGetTaskContainerById(w http.ResponseWriter, r *http.Request) {
	containerId := chi.URLParam(r, "containerID")
	if containerId == "" {
		h.logger.Error().Str("ErrorCode", constants.MissingParameter).Msg("Missing Container ID")
		response.BadRequestMissingParameters(w, "Missing container ID")
		return
	}

	// Use Query Bus
	result, err := h.queryBus.Execute(query.GetContainerByIdQuery{ContainerId: containerId})
	if err != nil {
		h.logger.Error().Err(err).Str("ErrorCode", TaskContainerGetNotFound).Msg(err.Error())
		response.NotFound(w, TaskContainerGetNotFound, "Container does not exist")
		return
	}
	response.WriteJsonWithEncode(w, http.StatusOK, result)
}
func (h *Handler) handleGetTaskContainersByGroupId(w http.ResponseWriter, r *http.Request) {
	groupIdVar := chi.URLParam(r, "usergroupID")
	if groupIdVar == "" {
		h.logger.Error().Str("ErrorCode", constants.MissingParameter).Msg("Missing Group ID")
		response.BadRequestMissingParameters(w, "Missing Group ID")
		return
	}
	groupId, err := strconv.Atoi(groupIdVar)
	if err != nil {
		h.logger.Error().Err(err).Str("ErrorCode", constants.InvalidParameter).Msg(err.Error())
		response.ErrorResponse(w, http.StatusBadRequest, *(response.New(constants.InvalidParameter, "Invalid Parameter", "Invalid Group ID")))
		return
	}

	// Use Query Bus
	result, err := h.queryBus.Execute(query.GetContainersByGroupIdQuery{GroupId: groupId})
	if err != nil {
		h.logger.Error().Err(err).Str("ErrorCode", TaskContainerGetNotFound).Msg(err.Error())
		response.NotFound(w, TaskContainerGetNotFound, "Error occurred during retrieving containers by group id")
		return
	}
	response.WriteJsonWithEncode(w, http.StatusOK, result)
}
func (h *Handler) handleCreateTaskContainer(w http.ResponseWriter, r *http.Request) {
	var createDto CreateContainerDto
	if err := response.ParseJson(r, &createDto); err != nil {
		h.logger.Error().Err(err).Str("ErrorCode", constants.RequestBodyError).Msg("Error occurred during parsing json of CreateContainerDto")
		response.InvalidJsonBody(w, "Error occurred during parsing json of CreateContainerDto")
		return
	}

	// Use Command Bus
	cmd := command.CreateContainerCommand{
		Name:        createDto.Name,
		Description: createDto.Description,
		Type:        createDto.Type,
		UserGroupId: createDto.UserGroupId,
	}

	result, err := h.commandBus.Execute(cmd)
	if err != nil {
		h.logger.Error().Err(err).Str("ErrorCode", TaskContainerServerError).Msg(err.Error())
		response.InternalServerError(w, "Error occurred during creating container")
		return
	}

	containerId := result.(string)
	response.WriteJsonWithEncode(w, http.StatusCreated, containerId)
}
func (h *Handler) handleDeleteTaskContainer(w http.ResponseWriter, r *http.Request) {
	containerId := chi.URLParam(r, "containerID")

	// Use Command Bus
	cmd := command.DeleteContainerCommand{ContainerId: containerId}
	_, err := h.commandBus.Execute(cmd)
	if err != nil {
		h.logger.Error().Err(err).Str("ErrorCode", DeleteTaskContainerError).Msg(err.Error())
		response.NotFound(w, DeleteTaskContainerError, "Error occurred during delete container")
		return
	}

	response.WriteJsonWithEncode(w, http.StatusNoContent, "task container is removed.")
}
