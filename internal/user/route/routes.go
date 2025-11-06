package route

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth"
	"github.com/happYness-Project/taskManagementGolang/internal/user/model"
	"github.com/happYness-Project/taskManagementGolang/internal/user/repository"
	userRepo "github.com/happYness-Project/taskManagementGolang/internal/user/repository"
	userGroupRepo "github.com/happYness-Project/taskManagementGolang/internal/usergroup/repository"
	"github.com/happYness-Project/taskManagementGolang/pkg/constants"
	"github.com/happYness-Project/taskManagementGolang/pkg/errors"
	"github.com/happYness-Project/taskManagementGolang/pkg/loggers"
	"github.com/happYness-Project/taskManagementGolang/pkg/response"
)

type Handler struct {
	logger        *loggers.AppLogger
	userRepo      userRepo.UserRepository
	userGroupRepo userGroupRepo.UserGroupRepository
}

func NewHandler(logger *loggers.AppLogger, repo repository.UserRepository, ugRepo userGroupRepo.UserGroupRepository) *Handler {
	return &Handler{logger: logger, userRepo: repo, userGroupRepo: ugRepo}
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

	if r.URL.Query().Get("email") != "" {
		h.responseUser(w, "email", r.URL.Query().Get("email"))
		return
	} else if r.URL.Query().Get("username") != "" {
		h.responseUser(w, "username", r.URL.Query().Get("username"))
		return
	}
	users, err := h.userRepo.GetAllUsers()
	if err != nil {
		h.logger.Error().Err(err).Str("ErrorCode", constants.ServerError).Msg("Error occurred during GetAllUsers.")
		response.InternalServerError(w)
		return
	}
	response.SuccessJson(w, users, "success", http.StatusOK)
}
func (h *Handler) handleGetUser(w http.ResponseWriter, r *http.Request) {
	user, err := h.userRepo.GetUserByUserId(chi.URLParam(r, "userID"))
	if err != nil {
		h.logger.Error().Err(err).Str("ErrorCode", constants.ServerError).Msg("Error occurred during GetUserByUserId")
		response.InternalServerError(w, "Error occurred during retrieving user.")
		return
	}
	if user == nil {
		h.logger.Error().Str("ErrorCode", UserGetNotFound)
		response.NotFound(w, UserGetNotFound, "user does not exist")
		return
	}

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

	usersWithRoles, err := h.userRepo.GetUsersByGroupIdWithRoles(groupId)
	if err != nil {
		h.logger.Error().Err(err).Msg("Error during Get Users by Group ID with Roles")
		response.InternalServerError(w, "Error occurred during retrieving users by group ID")
		return
	}

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
	var user *model.User
	var err error
	if findField == "email" {
		user, err = h.userRepo.GetUserByEmail(findVar)
	} else if findField == "username" {
		user, err = h.userRepo.GetUserByUsername(findVar)
	}
	if err != nil {
		h.logger.Error().Err(err).Str("ErrorCode", constants.ServerError).Msg("Error occurred during responseUser.")
		response.InternalServerError(w)
		return
	}
	if user.Id == 0 {
		h.logger.Error().Err(err).Str("ErrorCode", UserGetNotFound)
		response.NotFound(w, UserGetNotFound, "cannot find an user")
		return
	}

	userDetailDto := new(UserDetailDto)
	ugs, err := h.userGroupRepo.GetUserGroupsByUserId(user.Id) // TODO Let's double check this - what if there are server errors.
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

	user := model.NewUser(userId, createDto.UserName, createDto.FirstName, createDto.LastName, createDto.Email)

	err := h.userRepo.CreateUser(*user)
	if err != nil {
		h.logger.Error().Err(err).Str("ErrorCode", UserCreateServerError).Msg(err.Error())
		response.ErrorResponse(w, http.StatusBadRequest, *response.New(UserCreateServerError, "Bad Request", "cannot create a user"))
		// different error message is required for this
		// if user already exists, conflict error should return
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

	user, err := h.userRepo.GetUserByUserId(chi.URLParam(r, "userID"))
	if err != nil || user == nil {
		h.logger.Error().Err(err).Str("ErrorCode", UserGetNotFound).Msg("Error occurred during GetUserByUserId")
		response.NotFound(w, UserGetNotFound, "cannot find a user")
		return
	}

	user.UpdateUser(updateDto.FirstName, updateDto.LastName, updateDto.Email) // Todo : Domain Validation error code.

	err = h.userRepo.UpdateUser(*user)
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

	user, err := h.userRepo.GetUserByUserId(chi.URLParam(r, "userID"))
	if err != nil || user == nil {
		h.logger.Error().Err(err).Str("ErrorCode", UserGetNotFound) // what if there was an server error.
		response.NotFound(w, UserGetNotFound, "Not able to find a user")
		return
	}

	err = user.UpdateDefaultGroupId(jsonBody.DefaultGroupId)
	if err != nil {
		h.logger.Error().Err(err).Str("ErrorCode", UserDomainError).Msg(err.Error())
		response.BadRequestDomainError(w, UserDomainError, err.Error())
		return
	}

	err = h.userRepo.UpdateUser(*user)
	if err != nil {
		h.logger.Error().Err(err).Str("ErrorCode", UserUpdateServerError).Msg(err.Error())
		response.ErrorResponse(w, http.StatusBadRequest, *(response.New(UserUpdateServerError, errors.Badrequest, "Error ocurred during update an user")))
		return
	}

	response.SuccessJson(w, nil, "Default user group ID is updated.", http.StatusOK)
}
