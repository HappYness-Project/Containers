package route

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth"
	"github.com/happYness-Project/taskManagementGolang/internal/user/application"
	"github.com/happYness-Project/taskManagementGolang/internal/user/application/command"
	"github.com/happYness-Project/taskManagementGolang/internal/user/application/query"
	"github.com/happYness-Project/taskManagementGolang/internal/user/domain"
	"github.com/happYness-Project/taskManagementGolang/internal/user/repository"
	userGroupRepo "github.com/happYness-Project/taskManagementGolang/internal/usergroup/repository"
	"github.com/happYness-Project/taskManagementGolang/pkg/constants"
	"github.com/happYness-Project/taskManagementGolang/pkg/loggers"
	"github.com/happYness-Project/taskManagementGolang/pkg/response"
)

type Handler struct {
	logger        *loggers.AppLogger
	commandBus    *application.CommandBus
	queryBus      *application.QueryBus
	userGroupRepo userGroupRepo.UserGroupRepository // Keep for now (used in user detail)
}

func NewHandler(logger *loggers.AppLogger, repo repository.UserRepository, ugRepo userGroupRepo.UserGroupRepository) *Handler {
	return &Handler{
		logger:        logger,
		commandBus:    application.NewCommandBus(repo),
		queryBus:      application.NewQueryBus(repo),
		userGroupRepo: ugRepo,
	}
}

func (h *Handler) RegisterRoutes(router chi.Router) {
	router.Route("/api/users", func(r chi.Router) {
		r.Get("/", h.handleGetUsers)
		r.Post("/", h.handleCreateUser)
		r.Put("/{userID}", h.handleUpdateUser)
		r.Get("/{userID}", h.handleGetUser)
		r.Patch("/{userID}/default-group", h.handleUpdateGroupId)
	})
	router.Get("/api/user-groups/{groupID}/users", h.handleGetUsersByGroupId)

}

func (h *Handler) handleGetUsers(w http.ResponseWriter, r *http.Request) {
	// Handle query parameters for email/username search
	if email := r.URL.Query().Get("email"); email != "" {
		h.responseUser(w, "email", email)
		return
	} else if username := r.URL.Query().Get("username"); username != "" {
		h.responseUser(w, "username", username)
		return
	}

	// Use Query Bus
	result, err := h.queryBus.Execute(query.GetAllUsersQuery{})
	if err != nil {
		h.logger.Error().Err(err).Str("ErrorCode", constants.ServerError).Msg("Error occurred during GetAllUsers.")
		response.InternalServerError(w)
		return
	}
	response.SuccessJson(w, result, "success", http.StatusOK)
}
func (h *Handler) handleGetUser(w http.ResponseWriter, r *http.Request) {
	// Use Query Bus
	result, err := h.queryBus.Execute(query.GetUserByIdQuery{UserId: chi.URLParam(r, "userID")})
	if err != nil {
		h.logger.Error().Err(err).Str("ErrorCode", constants.ServerError).Msg("Error occurred during GetUserByUserId")
		response.InternalServerError(w, "Error occurred during retrieving user.")
		return
	}

	user := result.(*domain.User)

	// Get user's groups (this will be moved to a read model later)
	userDetailDto := new(UserDetailDto)
	ugs, err := h.userGroupRepo.GetUserGroupsByUserId(user.Id)
	if err != nil {
		h.logger.Error().Err(err).Msg("not able to get usergroups by user id.")
		response.NotFound(w, UserGetNotFound)
		return
	}
	userDetailDto.Id = user.Id
	userDetailDto.UserId = user.UserId
	userDetailDto.UserName = user.UserName
	userDetailDto.FirstName = user.FirstName
	userDetailDto.LastName = user.LastName
	userDetailDto.CreatedAt = user.CreatedAt
	userDetailDto.UpdatedAt = user.UpdatedAt
	userDetailDto.Email = user.Email
	userDetailDto.IsActive = user.IsActive
	userDetailDto.UserGroup = ugs
	userDetailDto.DefaultGroupId = user.DefaultGroupId

	response.SuccessJson(w, userDetailDto, "success", http.StatusOK)
}
func (h *Handler) handleGetUsersByGroupId(w http.ResponseWriter, r *http.Request) {
	vars := chi.URLParam(r, "groupID")
	if vars == "" {
		h.logger.Error().Msg("Missing Group ID")
		response.BadRequestMissingParameters(w, "Missing Group Id")
		return
	}
	groupId, err := strconv.Atoi(vars)
	if err != nil {
		h.logger.Error().Err(err).Msg("Invalid Group ID")
		response.ErrorResponse(w, http.StatusBadRequest, *(response.New(constants.InvalidParameter, "Invalid parameter", "Invalid Group Id")))
		return
	}

	// Use Query Bus
	result, err := h.queryBus.Execute(query.GetUsersByGroupIdQuery{GroupId: groupId})
	if err != nil {
		h.logger.Error().Err(err).Msg("Error during Get Users by Group ID with Roles")
		response.InternalServerError(w, "Error occurred during retrieving users by group ID")
		return
	}

	usersWithRoles := result.([]*domain.UserWithRole)

	// Convert to DTOs
	var userDtos []*UserWithRoleDto
	for _, userWithRole := range usersWithRoles {
		dto := &UserWithRoleDto{
			Id:             userWithRole.User.Id,
			UserId:         userWithRole.User.UserId,
			UserName:       userWithRole.User.UserName,
			FirstName:      userWithRole.User.FirstName,
			LastName:       userWithRole.User.LastName,
			CreatedAt:      userWithRole.User.CreatedAt,
			UpdatedAt:      userWithRole.User.UpdatedAt,
			Email:          userWithRole.User.Email,
			IsActive:       userWithRole.User.IsActive,
			DefaultGroupId: userWithRole.User.DefaultGroupId,
			Role:           userWithRole.Role,
			JoinedAt:       userWithRole.JoinedAt,
		}
		userDtos = append(userDtos, dto)
	}

	response.SuccessJson(w, userDtos, "success", http.StatusOK)
}
func (h *Handler) responseUser(w http.ResponseWriter, findField string, findVar string) {
	var result interface{}
	var err error

	// Use Query Bus
	if findField == "email" {
		result, err = h.queryBus.Execute(query.GetUserByEmailQuery{Email: findVar})
	} else if findField == "username" {
		result, err = h.queryBus.Execute(query.GetUserByUsernameQuery{Username: findVar})
	}

	if err != nil {
		h.logger.Error().Err(err).Str("ErrorCode", constants.ServerError).Msg("Error occurred during responseUser.")
		response.InternalServerError(w)
		return
	}

	user := result.(*domain.User)

	userDetailDto := new(UserDetailDto)
	ugs, err := h.userGroupRepo.GetUserGroupsByUserId(user.Id)
	if err != nil {
		h.logger.Error().Err(err).Str("ErrorCode", GetUserGroupsNotFound).Msg("Error occurred during GetUserGroupsByUserId")
		response.NotFound(w, GetUserGroupsNotFound, err.Error())
		return
	}

	userDetailDto.Id = user.Id
	userDetailDto.UserId = user.UserId
	userDetailDto.UserName = user.UserName
	userDetailDto.FirstName = user.FirstName
	userDetailDto.LastName = user.LastName
	userDetailDto.CreatedAt = user.CreatedAt
	userDetailDto.UpdatedAt = user.UpdatedAt
	userDetailDto.Email = user.Email
	userDetailDto.IsActive = user.IsActive
	userDetailDto.UserGroup = ugs
	userDetailDto.DefaultGroupId = user.DefaultGroupId

	response.SuccessJson(w, userDetailDto, "success", http.StatusOK)
}
func (h *Handler) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	var createDto CreateUserDto
	if err := response.ParseJson(r, &createDto); err != nil {
		h.logger.Error().Err(err).Str("ErrorCode", constants.RequestBodyError).Msg("Invalid JSON body for CreateUserDto")
		response.InvalidJsonBody(w)
		return
	}
	_, claims, _ := jwtauth.FromContext(r.Context())
	userId := fmt.Sprintf("%v", claims["nameid"])

	// Use Command Bus
	cmd := command.CreateUserCommand{
		UserId:    userId,
		UserName:  createDto.UserName,
		FirstName: createDto.FirstName,
		LastName:  createDto.LastName,
		Email:     createDto.Email,
	}

	_, err := h.commandBus.Execute(cmd)
	if err != nil {
		h.logger.Error().Err(err).Str("ErrorCode", UserCreateServerError).Msg(err.Error())
		response.ErrorResponse(w, http.StatusBadRequest, *response.New(UserCreateServerError, "Bad Request", "cannot create a user"))
		return
	}
	response.SuccessJson(w, map[string]string{"user_id": userId}, "User is created.", http.StatusCreated)
}
func (h *Handler) handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	var updateDto UpdateUserDto
	if err := response.ParseJson(r, &updateDto); err != nil {
		h.logger.Error().Err(err).Str("ErrorCode", constants.RequestBodyError).Msg("Invalid JSON body for UpdateUserDto")
		response.InvalidJsonBody(w, "Json body invalid for update user.")
		return
	}

	// Use Command Bus
	cmd := command.UpdateUserCommand{
		UserId:    chi.URLParam(r, "userID"),
		FirstName: updateDto.FirstName,
		LastName:  updateDto.LastName,
		Email:     updateDto.Email,
	}

	_, err := h.commandBus.Execute(cmd)
	if err != nil {
		h.logger.Error().Err(err).Str("ErrorCode", UserUpdateServerError).Msg("Error occurred during UpdateUser")
		response.ErrorResponse(w, http.StatusBadRequest, *(response.New(UserUpdateServerError, "Bad Request", "Error occurred during Update user.")))
		return
	}
	response.SuccessJson(w, nil, "User is updated.", http.StatusNoContent)
}

func (h *Handler) handleUpdateGroupId(w http.ResponseWriter, r *http.Request) {
	type JsonBody struct {
		DefaultGroupId int `json:"default_group_id"`
	}
	var jsonBody JsonBody
	err := json.NewDecoder(r.Body).Decode(&jsonBody)
	if err != nil {
		h.logger.Error().Err(err).Str("ErrorCode", constants.RequestBodyError).Msg(err.Error())
		response.InvalidJsonBody(w, "Invalid json format for default_group_id")
		return
	}

	// Use Command Bus
	cmd := command.UpdateDefaultGroupCommand{
		UserId:         chi.URLParam(r, "userID"),
		DefaultGroupId: jsonBody.DefaultGroupId,
	}

	_, err = h.commandBus.Execute(cmd)
	if err != nil {
		h.logger.Error().Err(err).Str("ErrorCode", UserDomainError).Msg(err.Error())
		response.BadRequestDomainError(w, UserDomainError, err.Error())
		return
	}

	response.SuccessJson(w, nil, "Default user group ID is updated.", http.StatusOK)
}
