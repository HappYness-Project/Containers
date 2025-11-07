package route

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth"
	userRepo "github.com/happYness-Project/taskManagementGolang/internal/user/repository"
	"github.com/happYness-Project/taskManagementGolang/internal/usergroup/application"
	"github.com/happYness-Project/taskManagementGolang/internal/usergroup/application/command"
	"github.com/happYness-Project/taskManagementGolang/internal/usergroup/application/query"

	"github.com/happYness-Project/taskManagementGolang/internal/usergroup/repository"
	"github.com/happYness-Project/taskManagementGolang/pkg/constants"
	"github.com/happYness-Project/taskManagementGolang/pkg/loggers"
	"github.com/happYness-Project/taskManagementGolang/pkg/response"
)

type Handler struct {
	logger     *loggers.AppLogger
	commandBus *application.CommandBus
	queryBus   *application.QueryBus
}

func NewHandler(logger *loggers.AppLogger, repo repository.UserGroupRepository, userRepo userRepo.UserRepository) *Handler {
	return &Handler{
		logger:     logger,
		commandBus: application.NewCommandBus(repo, userRepo),
		queryBus:   application.NewQueryBus(repo, userRepo),
	}
}
func (h *Handler) RegisterRoutes(router chi.Router) {
	router.Route("/api/user-groups", func(r chi.Router) {
		r.Get("/", h.handleGetUserGroups)
		r.Get("/{groupID}", h.handleGetUserGroupById)
		r.Delete("/{groupID}", h.handleDeleteUserGroup)
		r.Post("/{groupID}/users", h.handleAddUserToGroup)
		r.Delete("/{groupID}/users/{userID}", h.handleRemoveUserFromGroup)
		r.Patch("/{groupID}/users/{userID}/role", h.handleUpdateUserRoleInGroup)
	})
	router.Post("/api/user-groups", h.handleCreateUserGroup)
	router.Get("/api/users/{userID}/user-groups", h.handleGetUserGroupByUserId)
}

func (h *Handler) handleGetUserGroups(w http.ResponseWriter, r *http.Request) {
	// Use Query Bus
	result, err := h.queryBus.Execute(query.GetAllGroupsQuery{})
	if err != nil {
		h.logger.Error().Err(err).Str("ErrorCode", constants.ServerError).Msg("Error occurred during getting groups.")
		response.InternalServerError(w)
		return
	}
	response.WriteJsonWithEncode(w, http.StatusOK, result)
}
func (h *Handler) handleGetUserGroupById(w http.ResponseWriter, r *http.Request) {
	groupId, err := strconv.Atoi(chi.URLParam(r, "groupID"))
	if err != nil {
		h.logger.Error().Err(err).Str("ErrorCode", constants.InvalidParameter).Msg("Invalid Parameters for GroupID")
		response.ErrorResponse(w, http.StatusBadRequest, *response.New(constants.InvalidParameter, "Invalid Paramters", "Invalid Group ID"))
		return
	}

	// Use Query Bus
	result, err := h.queryBus.Execute(query.GetGroupByIdQuery{GroupId: groupId})
	if err != nil {
		h.logger.Error().Err(err).Str("ErrorCode", UserGroupGetNotFound).Msg(err.Error())
		response.NotFound(w, UserGroupGetNotFound, "group does not exist")
		return
	}

	response.WriteJsonWithEncode(w, http.StatusOK, result)
}

func (h *Handler) handleCreateUserGroup(w http.ResponseWriter, r *http.Request) {
	var createDto CreateUserGroupDto
	if err := response.ParseJson(r, &createDto); err != nil {
		h.logger.Error().Err(err).Str("ErrorCode", constants.RequestBodyError).Msg("Json body for CreateUserGroupRequest is invalid")
		response.InternalServerError(w)
		return
	}

	// Get user ID from JWT
	_, claims, _ := jwtauth.FromContext(r.Context())
	userid := fmt.Sprintf("%v", claims["nameid"])

	// Use Command Bus
	cmd := command.CreateGroupCommand{
		GroupName: createDto.GroupName,
		GroupDesc: createDto.GroupDesc,
		GroupType: createDto.GroupType,
		CreatorId: userid,
	}
	result, err := h.commandBus.Execute(cmd)
	if err != nil {
		h.logger.Error().Err(err).Str("ErrorCode", UserGroupCreationFailure).Msg(err.Error())
		response.ErrorResponse(w, http.StatusBadRequest, *response.New(UserGroupCreationFailure, "Failed to create group", err.Error()))
		return
	}

	groupId := result.(int)
	response.SuccessJson(w, map[string]int{"group_id": groupId}, "User group is created.", http.StatusCreated)
}

func (h *Handler) handleGetUserGroupByUserId(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "userID")
	if userId == "" {
		h.logger.Error().Msg("missing userID")
		response.BadRequestMissingParameters(w)
		return
	}

	// Use Query Bus
	result, err := h.queryBus.Execute(query.GetGroupsByUserIdQuery{UserId: userId})
	if err != nil {
		h.logger.Error().Err(err).Str("ErrorCode", UserGroupGetNotFound).Msg(err.Error())
		response.NotFound(w, UserGroupGetNotFound, "Not able to find user groups")
		return
	}

	response.WriteJsonWithEncode(w, http.StatusOK, result)
}

func (h *Handler) handleAddUserToGroup(w http.ResponseWriter, r *http.Request) {
	groupId, err := strconv.Atoi(chi.URLParam(r, "groupID"))
	if err != nil {
		h.logger.Error().Err(err).Str("ErrorCode", constants.InvalidParameter).Msg("invalid Group Id")
		response.BadRequestMissingParameters(w)
		return
	}

	type JsonBody struct {
		UserId string `json:"user_id"`
	}
	var jsonBody JsonBody
	err = json.NewDecoder(r.Body).Decode(&jsonBody)
	if err != nil {
		h.logger.Error().Err(err).Str("ErrorCode", constants.RequestBodyError)
		response.InvalidJsonBody(w, err.Error())
		return
	}

	// Use Command Bus
	cmd := command.AddMemberCommand{
		GroupId: groupId,
		UserId:  jsonBody.UserId,
	}

	_, err = h.commandBus.Execute(cmd)
	if err != nil {
		h.logger.Error().Err(err).Str("ErrorCode", UserGroupAddUserError).Msg(err.Error())
		response.ErrorResponse(w, http.StatusBadRequest, *response.New(UserGroupAddUserError, "Bad Request", err.Error()))
		return
	}

	response.WriteJsonWithEncode(w, http.StatusCreated, fmt.Sprintf("User is added to the user group ID: %d", groupId))
}

func (h *Handler) handleDeleteUserGroup(w http.ResponseWriter, r *http.Request) {
	groupId, err := strconv.Atoi(chi.URLParam(r, "groupID"))
	if err != nil {
		h.logger.Error().Err(err).Str("ErrorCode", constants.InvalidParameter).Msg("Invalid groupId")
		response.BadRequestMissingParameters(w, "invalid groupId")
		return
	}

	// Use Command Bus
	cmd := command.DeleteGroupCommand{GroupId: groupId}
	_, err = h.commandBus.Execute(cmd)
	if err != nil {
		h.logger.Error().Err(err).Str("ErrorCode", DeleteUserGroupError).Msg(err.Error())
		response.ErrorResponse(w, http.StatusBadRequest, *response.New(DeleteUserGroupError, "Bad Request", err.Error()))
		return
	}

	response.SuccessJson(w, nil, fmt.Sprintf("User group ID %d deleted successfully", groupId), 204)
}

func (h *Handler) handleRemoveUserFromGroup(w http.ResponseWriter, r *http.Request) {
	vars := chi.URLParam(r, "groupID")
	if vars == "" {
		h.logger.Error().Str("ErrorCode", constants.MissingParameter).Msg("Missing GroupID")
		response.BadRequestMissingParameters(w, "Missing Group ID")
		return
	}
	groupId, err := strconv.Atoi(vars)
	if err != nil {
		h.logger.Error().Err(err).Str("ErrorCode", constants.InvalidParameter).Msg(err.Error())
		response.ErrorResponse(w, http.StatusBadRequest, *response.New(constants.InvalidParameter, "Invalid Parameter", "Invalid Group ID"))
		return
	}

	userId := chi.URLParam(r, "userID")

	// Use Command Bus (it handles the default group clearing logic)
	cmd := command.RemoveMemberCommand{
		GroupId: groupId,
		UserId:  userId,
	}

	_, err = h.commandBus.Execute(cmd)
	if err != nil {
		h.logger.Error().Err(err).Str("ErrorCode", RemoveUserFromUserGroupError).Msg(err.Error())
		response.ErrorResponse(w, http.StatusBadRequest, *response.New(RemoveUserFromUserGroupError, "Bad Request", err.Error()))
		return
	}

	response.SuccessJson(w, nil, fmt.Sprintf("User is removed from user group ID: %d", groupId), 204)
}

func (h *Handler) handleUpdateUserRoleInGroup(w http.ResponseWriter, r *http.Request) {
	groupId, err := strconv.Atoi(chi.URLParam(r, "groupID"))
	if err != nil {
		h.logger.Error().Err(err).Str("ErrorCode", constants.InvalidParameter).Msg("Invalid Group ID")
		response.ErrorResponse(w, http.StatusBadRequest, *response.New(constants.InvalidParameter, "Invalid Parameter", "Invalid Group ID"))
		return
	}

	userId := chi.URLParam(r, "userID")

	var updateDto UpdateUserRoleDto
	if err := response.ParseJson(r, &updateDto); err != nil {
		h.logger.Error().Err(err).Str("ErrorCode", constants.RequestBodyError).Msg("Invalid JSON body for UpdateUserRoleRequest")
		response.InvalidJsonBody(w, "Invalid JSON body")
		return
	}

	// Use Command Bus (it validates the role using the Role value object)
	cmd := command.ChangeMemberRoleCommand{
		GroupId: groupId,
		UserId:  userId,
		NewRole: updateDto.Role,
	}

	_, err = h.commandBus.Execute(cmd)
	if err != nil {
		h.logger.Error().Err(err).Str("ErrorCode", UpdateUserRoleError).Msg(err.Error())
		response.ErrorResponse(w, http.StatusBadRequest, *response.New(UpdateUserRoleError, "Bad Request", err.Error()))
		return
	}

	response.SuccessJson(w, nil, fmt.Sprintf("User role updated to '%s' in group ID: %d", updateDto.Role, groupId), http.StatusOK)
}
