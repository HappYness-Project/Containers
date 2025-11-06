package route

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/happYness-Project/taskManagementGolang/internal/mocks"
	userModel "github.com/happYness-Project/taskManagementGolang/internal/user/model"
	userGroupModel "github.com/happYness-Project/taskManagementGolang/internal/usergroup/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/happYness-Project/taskManagementGolang/pkg/configs"
	"github.com/happYness-Project/taskManagementGolang/pkg/constants"
	"github.com/happYness-Project/taskManagementGolang/pkg/errors"
	"github.com/happYness-Project/taskManagementGolang/pkg/loggers"
	"github.com/happYness-Project/taskManagementGolang/pkg/response"
	"github.com/stretchr/testify/mock"
)

func TestUserHandler(t *testing.T) {

	env := configs.Env{}
	logger := loggers.Setup(env)
	mockUserRepo := new(mocks.MockUserRepo)
	mockUserGroupRepo := new(mocks.MockUserGroupRepo)
	handler := NewHandler(logger, mockUserRepo, mockUserGroupRepo)

	// Reset mocks before each test to prevent interference
	t.Cleanup(func() {
		mockUserRepo.ExpectedCalls = nil
		mockUserGroupRepo.ExpectedCalls = nil
	})

	t.Run("when get all users, Then return status code 200 with users", func(t *testing.T) {
		expectedUsers := []*userModel.User{{Id: 1, UserId: "1", UserName: "Alice"}}

		mockUserRepo.On("GetAllUsers").Return(expectedUsers, nil)
		req := httptest.NewRequest(http.MethodGet, "/api/users/", nil)
		rr := httptest.NewRecorder()
		router := chi.NewRouter()
		router.Get("/api/users/", handler.handleGetUsers)
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		// Parse response body
		var response map[string]interface{}
		err := json.Unmarshal(rr.Body.Bytes(), &response)
		require.NoError(t, err)

		// Verify response structure
		assert.Equal(t, "success", response["message"])
		assert.NotNil(t, response["data"])

		// Parse data as array of users
		dataBytes, _ := json.Marshal(response["data"])
		var users []userModel.User
		err = json.Unmarshal(dataBytes, &users)
		require.NoError(t, err)

		// Verify user data
		assert.Len(t, users, 1)
		assert.Equal(t, expectedUsers[0].Id, users[0].Id)
		assert.Equal(t, expectedUsers[0].UserId, users[0].UserId)
		assert.Equal(t, expectedUsers[0].UserName, users[0].UserName)

		mockUserRepo.AssertExpectations(t)
	})

	t.Run("when get user by ID, Then return status code 200 with user details", func(t *testing.T) {
		userID := "test-user-id"
		expectedUser := &userModel.User{
			Id:        1,
			UserId:    userID,
			UserName:  "testuser",
			FirstName: "Test",
			LastName:  "User",
			Email:     "test@example.com",
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		expectedUserGroups := []*userGroupModel.UserGroup{
			{GroupId: 1, GroupName: "Admin", GroupDesc: "Administrators"},
		}

		mockUserRepo.On("GetUserByUserId", userID).Return(expectedUser, nil)
		mockUserGroupRepo.On("GetUserGroupsByUserId", expectedUser.Id).Return(expectedUserGroups, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/users/"+userID, nil)
		rr := httptest.NewRecorder()
		router := chi.NewRouter()

		// Set up chi router with URL parameter
		router.Route("/api/users/{userID}", func(r chi.Router) {
			r.Get("/", handler.handleGetUser)
		})

		router.ServeHTTP(rr, req)

		// Parse response body
		var response map[string]interface{}
		err := json.Unmarshal(rr.Body.Bytes(), &response)
		require.NoError(t, err)

		dataBytes, _ := json.Marshal(response["data"])
		var userDetail UserDetailDto
		_ = json.Unmarshal(dataBytes, &userDetail)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, expectedUser.Id, userDetail.Id)
		assert.Equal(t, expectedUser.UserId, userDetail.UserId)
		assert.Equal(t, expectedUser.UserName, userDetail.UserName)
		assert.Equal(t, expectedUser.FirstName, userDetail.FirstName)
		assert.Equal(t, expectedUser.LastName, userDetail.LastName)
		assert.Equal(t, expectedUser.Email, userDetail.Email)
		assert.Equal(t, expectedUser.IsActive, userDetail.IsActive)
		assert.Equal(t, expectedUser.DefaultGroupId, userDetail.DefaultGroupId)

		assert.Len(t, userDetail.UserGroup, 1)
		assert.Equal(t, expectedUserGroups[0].GroupId, userDetail.UserGroup[0].GroupId)
		assert.Equal(t, expectedUserGroups[0].GroupName, userDetail.UserGroup[0].GroupName)
		assert.Equal(t, expectedUserGroups[0].GroupDesc, userDetail.UserGroup[0].GroupDesc)
	})

	t.Run("when get user by ID but error from repository layer, Then return status code 500", func(t *testing.T) {
		userID := "non-existent-user"

		mockUserRepo.On("GetUserByUserId", userID).Return((*userModel.User)(nil), fmt.Errorf("database error"))

		req := httptest.NewRequest(http.MethodGet, "/api/users/"+userID, nil)
		rr := httptest.NewRecorder()
		router := chi.NewRouter()

		router.Route("/api/users/{userID}", func(r chi.Router) {
			r.Get("/", handler.handleGetUser)
		})

		router.ServeHTTP(rr, req)

		var details response.ProblemDetails
		json.Unmarshal(rr.Body.Bytes(), &details)
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Equal(t, "Internal server error", details.Title)
		assert.Equal(t, "Error occurred during retrieving user.", details.Detail)
		assert.Equal(t, constants.ServerError, details.ErrorCode)
	})

	t.Run("when get user by ID but user not found, Then return status code 404", func(t *testing.T) {
		mockUserRepo.ExpectedCalls = nil
		mockUserGroupRepo.ExpectedCalls = nil
		userID := "non-existent-user"
		mockUserRepo.On("GetUserByUserId", userID).Return((*userModel.User)(nil), nil)

		req := httptest.NewRequest(http.MethodGet, "/api/users/"+userID, nil)
		rr := httptest.NewRecorder()
		router := chi.NewRouter()

		router.Route("/api/users/{userID}", func(r chi.Router) {
			r.Get("/", handler.handleGetUser)
		})

		router.ServeHTTP(rr, req)

		var details response.ProblemDetails
		json.Unmarshal(rr.Body.Bytes(), &details)
		assert.Equal(t, http.StatusNotFound, rr.Code)
		assert.Equal(t, "Not found", details.Title)
		assert.Equal(t, "user does not exist", details.Detail)
		assert.Equal(t, UserGetNotFound, details.ErrorCode)
	})

	t.Run("when get users by group ID, Then return status code 200 with users", func(t *testing.T) {
		groupID := "1"
		expectedUsersWithRoles := []*userModel.UserWithRole{
			{
				User: &userModel.User{Id: 1, UserId: "user1", UserName: "Alice"},
				Role: "admin",
			},
			{
				User: &userModel.User{Id: 2, UserId: "user2", UserName: "Bob"},
				Role: "member",
			},
		}

		mockUserRepo.On("GetUsersByGroupIdWithRoles", 1).Return(expectedUsersWithRoles, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/user-groups/"+groupID+"/users", nil)
		rr := httptest.NewRecorder()
		router := chi.NewRouter()

		router.Route("/api/user-groups/{groupID}", func(r chi.Router) {
			r.Get("/users", handler.handleGetUsersByGroupId)
		})

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		// Parse response body
		var response map[string]interface{}
		err := json.Unmarshal(rr.Body.Bytes(), &response)
		require.NoError(t, err)

		dataBytes, _ := json.Marshal(response["data"])
		var users []UserWithRoleDto
		err = json.Unmarshal(dataBytes, &users)
		require.NoError(t, err)

		assert.Len(t, users, 2)
		assert.Equal(t, expectedUsersWithRoles[0].User.Id, users[0].Id)
		assert.Equal(t, expectedUsersWithRoles[0].User.UserId, users[0].UserId)
		assert.Equal(t, expectedUsersWithRoles[0].User.UserName, users[0].UserName)
		assert.Equal(t, expectedUsersWithRoles[0].Role, users[0].Role)
		assert.Equal(t, expectedUsersWithRoles[1].User.Id, users[1].Id)
		assert.Equal(t, expectedUsersWithRoles[1].User.UserId, users[1].UserId)
		assert.Equal(t, expectedUsersWithRoles[1].User.UserName, users[1].UserName)
		assert.Equal(t, expectedUsersWithRoles[1].Role, users[1].Role)
	})

	t.Run("when get users by group ID with invalid group ID, Then return status code 400", func(t *testing.T) {
		groupID := "invalid"

		req := httptest.NewRequest(http.MethodGet, "/api/user-groups/"+groupID+"/users", nil)
		rr := httptest.NewRecorder()
		router := chi.NewRouter()

		router.Route("/api/user-groups/{groupID}", func(r chi.Router) {
			r.Get("/users", handler.handleGetUsersByGroupId)
		})

		router.ServeHTTP(rr, req)

		var details response.ProblemDetails
		json.Unmarshal(rr.Body.Bytes(), &details)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, "Invalid parameter", details.Title)
		assert.Equal(t, "Invalid Group Id", details.Detail)
		assert.Equal(t, constants.InvalidParameter, details.ErrorCode)
	})

	t.Run("when get users by group ID with empty group ID, Then return status code 400", func(t *testing.T) {
		mockUserRepo.ExpectedCalls = nil
		mockUserGroupRepo.ExpectedCalls = nil

		req := httptest.NewRequest(http.MethodGet, "/api/user-groups//users", nil)
		rr := httptest.NewRecorder()
		router := chi.NewRouter()

		router.Route("/api/user-groups/{groupID}", func(r chi.Router) {
			r.Get("/users", handler.handleGetUsersByGroupId)
		})

		router.ServeHTTP(rr, req)

		var details response.ProblemDetails
		json.Unmarshal(rr.Body.Bytes(), &details)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, errors.Badrequest, details.Title)
		assert.Equal(t, "Missing Group Id", details.Detail)
		assert.Equal(t, constants.MissingParameter, details.ErrorCode)
	})

	t.Run("when update user default group ID with valid data, Then return status code 200", func(t *testing.T) {
		userID := "test-user-id"
		existingUser := &userModel.User{
			Id:             1,
			UserId:         userID,
			UserName:       "testuser",
			FirstName:      "Test",
			LastName:       "User",
			Email:          "test@example.com",
			IsActive:       true,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
			DefaultGroupId: 0,
		}

		mockUserRepo.On("GetUserByUserId", userID).Return(existingUser, nil)
		mockUserRepo.On("UpdateUser", mock.MatchedBy(func(u userModel.User) bool {
			return u.Id == existingUser.Id &&
				u.UserName == existingUser.UserName &&
				u.FirstName == existingUser.FirstName &&
				u.LastName == existingUser.LastName &&
				u.Email == existingUser.Email &&
				u.IsActive == existingUser.IsActive &&
				u.CreatedAt.Equal(existingUser.CreatedAt)
		})).Return(nil)

		requestBody := `{"default_group_id": 5}`
		req := httptest.NewRequest(http.MethodPatch, "/api/users/"+userID+"/default-group", strings.NewReader(requestBody))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		router := chi.NewRouter()

		router.Route("/api/users/{userID}", func(r chi.Router) {
			r.Patch("/default-group", handler.handleUpdateGroupId)
		})

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		var response map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &response)
		assert.Contains(t, response["message"], "Default user group ID is updated.")
	})

	t.Run("when update user default group ID with negative group ID, Then return status code 400", func(t *testing.T) {
		userID := "test-user-id"
		existingUser := &userModel.User{
			Id:             1,
			UserId:         userID,
			UserName:       "testuser",
			FirstName:      "Test",
			LastName:       "User",
			Email:          "test@example.com",
			IsActive:       true,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
			DefaultGroupId: 0,
		}
		mockUserRepo.On("GetUserByUserId", userID).Return(existingUser, nil)
		requestBody := `{"default_group_id": -1}`
		req := httptest.NewRequest(http.MethodPatch, "/api/users/"+userID+"/default-group", strings.NewReader(requestBody))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		router := chi.NewRouter()

		router.Route("/api/users/{userID}", func(r chi.Router) {
			r.Patch("/default-group", handler.handleUpdateGroupId)
		})

		router.ServeHTTP(rr, req)

		var details response.ProblemDetails
		json.Unmarshal(rr.Body.Bytes(), &details)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, "Domain validation error", details.Title)
		assert.Equal(t, "group ID cannot be negative", details.Detail)
		assert.Equal(t, UserDomainError, details.ErrorCode)
	})

	t.Run("when update user default group ID with same group ID, Then return status code 400", func(t *testing.T) {
		mockUserRepo.ExpectedCalls = nil
		mockUserGroupRepo.ExpectedCalls = nil
		userID := "test-user-id"
		existingUser := &userModel.User{
			Id:             1,
			UserId:         userID,
			UserName:       "testuser",
			FirstName:      "Test",
			LastName:       "User",
			Email:          "test@example.com",
			IsActive:       true,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
			DefaultGroupId: 5, // Already has group ID 5
		}
		mockUserRepo.On("GetUserByUserId", userID).Return(existingUser, nil)

		requestBody := `{"default_group_id": 5}`
		req := httptest.NewRequest(http.MethodPatch, "/api/users/"+userID+"/default-group", strings.NewReader(requestBody))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		router := chi.NewRouter()

		router.Route("/api/users/{userID}", func(r chi.Router) {
			r.Patch("/default-group", handler.handleUpdateGroupId)
		})

		router.ServeHTTP(rr, req)

		var details response.ProblemDetails
		json.Unmarshal(rr.Body.Bytes(), &details)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, "Domain validation error", details.Title)
		assert.Equal(t, "group ID is already set to the specified value", details.Detail)
		assert.Equal(t, UserDomainError, details.ErrorCode)
	})

	t.Run("when update user default group ID but user not found, Then return status code 404", func(t *testing.T) {
		mockUserRepo.ExpectedCalls = nil
		mockUserGroupRepo.ExpectedCalls = nil

		userID := "non-existent-user"
		mockUserRepo.On("GetUserByUserId", userID).Return((*userModel.User)(nil), nil)
		requestBody := `{"default_group_id": 5}`
		req := httptest.NewRequest(http.MethodPatch, "/api/users/"+userID+"/default-group", strings.NewReader(requestBody))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		router := chi.NewRouter()

		router.Route("/api/users/{userID}", func(r chi.Router) {
			r.Patch("/default-group", handler.handleUpdateGroupId)
		})

		router.ServeHTTP(rr, req)

		var details response.ProblemDetails
		json.Unmarshal(rr.Body.Bytes(), &details)
		assert.Equal(t, http.StatusNotFound, rr.Code)
		assert.Equal(t, "Not found", details.Title)
		assert.Equal(t, "Not able to find a user", details.Detail)
		assert.Equal(t, UserGetNotFound, details.ErrorCode)
	})

	t.Run("when update user default group ID with invalid JSON, Then return status code 400", func(t *testing.T) {
		userID := "test-user-id"
		existingUser := &userModel.User{
			Id:             1,
			UserId:         userID,
			UserName:       "testuser",
			FirstName:      "Test",
			LastName:       "User",
			Email:          "test@example.com",
			IsActive:       true,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
			DefaultGroupId: 0,
		}
		mockUserRepo.On("GetUserByUserId", userID).Return(existingUser, nil)

		requestBody := `{"default_group_id": "invalid"}`
		req := httptest.NewRequest(http.MethodPatch, "/api/users/"+userID+"/default-group", strings.NewReader(requestBody))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		router := chi.NewRouter()

		router.Route("/api/users/{userID}", func(r chi.Router) {
			r.Patch("/default-group", handler.handleUpdateGroupId)
		})

		router.ServeHTTP(rr, req)

		var details response.ProblemDetails
		json.Unmarshal(rr.Body.Bytes(), &details)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, "Invalid json body format", details.Title)
		assert.Equal(t, "Invalid json format for default_group_id", details.Detail)
		assert.Equal(t, constants.RequestBodyError, details.ErrorCode)
	})

	t.Run("when update user with valid data, Then return status code 204", func(t *testing.T) {
		mockUserRepo.ExpectedCalls = nil
		mockUserGroupRepo.ExpectedCalls = nil
		userID := "test-user-id"
		existingUser := &userModel.User{
			Id:        1,
			UserId:    userID,
			UserName:  "testuser",
			FirstName: "OldFirst",
			LastName:  "OldLast",
			Email:     "old@example.com",
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockUserRepo.On("GetUserByUserId", userID).Return(existingUser, nil)
		existingUser.UpdateUser("NewFirst", "NewLast", "new@example.com")
		mockUserRepo.On("UpdateUser", mock.MatchedBy(func(u userModel.User) bool {
			return u.Id == existingUser.Id &&
				u.UserId == existingUser.UserId &&
				u.UserName == existingUser.UserName &&
				u.FirstName == "NewFirst" &&
				u.LastName == "NewLast" &&
				u.Email == "new@example.com" &&
				u.IsActive == existingUser.IsActive &&
				u.CreatedAt.Equal(existingUser.CreatedAt)
		})).Return(nil)

		requestBody := `{"first_name": "NewFirst", "last_name": "NewLast", "email": "new@example.com"}`
		req := httptest.NewRequest(http.MethodPut, "/api/users/"+userID, strings.NewReader(requestBody))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		router := chi.NewRouter()
		router.Route("/api/users/{userID}", func(r chi.Router) {
			r.Put("/", handler.handleUpdateUser)
		})

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNoContent, rr.Code)
		var response map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &response)
		assert.Equal(t, "User is updated.", response["message"])
	})

	t.Run("when update user but user not found, Then return status code 404", func(t *testing.T) {
		mockUserRepo.ExpectedCalls = nil
		mockUserGroupRepo.ExpectedCalls = nil
		userID := "non-existent-user"
		mockUserRepo.On("GetUserByUserId", userID).Return((*userModel.User)(nil), nil)

		requestBody := `{"first_name": "NewFirst", "last_name": "NewLast", "email": "new@example.com"}`
		req := httptest.NewRequest(http.MethodPut, "/api/users/"+userID, strings.NewReader(requestBody))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		router := chi.NewRouter()
		router.Route("/api/users/{userID}", func(r chi.Router) {
			r.Put("/", handler.handleUpdateUser)
		})

		router.ServeHTTP(rr, req)

		var details response.ProblemDetails
		json.Unmarshal(rr.Body.Bytes(), &details)
		assert.Equal(t, http.StatusNotFound, rr.Code)
		assert.Equal(t, "Not found", details.Title)
		assert.Equal(t, "cannot find a user", details.Detail)
		assert.Equal(t, UserGetNotFound, details.ErrorCode)
	})

}
