package route

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os/user"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/happYness-Project/taskManagementGolang/internal/mocks"
	"github.com/happYness-Project/taskManagementGolang/internal/taskcontainer/domain"
	"github.com/happYness-Project/taskManagementGolang/pkg/configs"
	"github.com/happYness-Project/taskManagementGolang/pkg/loggers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTaskContainerHandler(t *testing.T) {
	env := configs.Env{}
	logger := loggers.Setup(env)
	mockContainerRepo := new(mocks.MockContainerRepo)
	mockUserRepo := new(mocks.MockUserRepo)
	handler := NewHandler(logger, mockContainerRepo, mockUserRepo)

	t.Run("when get all task containers, Then return status code 200 and containers array", func(t *testing.T) {
		// Arrange
		expectedContainers := []*domain.TaskContainer{
			{Id: "1", Name: "Container1", Description: "Desc1", Type: "typeA", IsActive: true, Activity_level: 0, UsergroupId: 2},
		}
		mockContainerRepo.On("AllTaskContainers").Return(expectedContainers, nil)
		req := httptest.NewRequest(http.MethodGet, "/api/task-containers", nil)
		rr := httptest.NewRecorder()
		router := chi.NewRouter()
		router.Get("/api/task-containers", handler.handleGetTaskContainers)

		// Act
		router.ServeHTTP(rr, req)

		// Assert
		assert.Equal(t, http.StatusOK, rr.Code)
		var containers []*domain.TaskContainer
		err := json.Unmarshal(rr.Body.Bytes(), &containers)
		require.NoError(t, err)
		assert.Len(t, containers, 1)
		assert.Equal(t, expectedContainers[0].Id, containers[0].Id)
		assert.Equal(t, expectedContainers[0].Name, containers[0].Name)
		assert.Equal(t, expectedContainers[0].Description, containers[0].Description)
		assert.Equal(t, expectedContainers[0].Type, containers[0].Type)
		assert.Equal(t, expectedContainers[0].IsActive, containers[0].IsActive)
		assert.Equal(t, expectedContainers[0].Activity_level, containers[0].Activity_level)
		assert.Equal(t, expectedContainers[0].UsergroupId, containers[0].UsergroupId)
		mockContainerRepo.AssertExpectations(t)
	})

	t.Run("when get task container by containerId, Then return status code 200 and container", func(t *testing.T) {
		containerId := "abcd"
		expectedContainer := &domain.TaskContainer{Id: containerId, Name: "Container2", Description: "Desc2", Type: "typeB", IsActive: false, Activity_level: 1, UsergroupId: 3}
		mockContainerRepo.On("GetById", containerId).Return(expectedContainer, nil)
		req := httptest.NewRequest(http.MethodGet, "/api/task-containers/"+containerId, nil)
		rr := httptest.NewRecorder()
		router := chi.NewRouter()
		router.Get("/api/task-containers/{containerID}", handler.handleGetTaskContainerById)

		// Act
		router.ServeHTTP(rr, req)

		// Assert
		assert.Equal(t, http.StatusOK, rr.Code)
		var container domain.TaskContainer
		err := json.Unmarshal(rr.Body.Bytes(), &container)
		require.NoError(t, err)
		assert.Equal(t, expectedContainer.Id, container.Id)
		assert.Equal(t, expectedContainer.Name, container.Name)
		assert.Equal(t, expectedContainer.Description, container.Description)
		assert.Equal(t, expectedContainer.Type, container.Type)
		assert.Equal(t, expectedContainer.IsActive, container.IsActive)
		assert.Equal(t, expectedContainer.Activity_level, container.Activity_level)
		assert.Equal(t, expectedContainer.UsergroupId, container.UsergroupId)
		mockContainerRepo.AssertExpectations(t)
	})
}

type mockContainerRepo struct{}

func (m *mockContainerRepo) AllTaskContainers() ([]*domain.TaskContainer, error) {
	return []*domain.TaskContainer{}, nil
}
func (m *mockContainerRepo) GetById(id string) (*domain.TaskContainer, error) {
	return &domain.TaskContainer{}, nil
}
func (m *mockContainerRepo) GetContainersByGroupId(groupId int) ([]domain.TaskContainer, error) {
	return []domain.TaskContainer{}, nil
}
func (m *mockContainerRepo) CreateContainer(c domain.TaskContainer) error {
	return nil
}
func (m *mockContainerRepo) DeleteContainer(id string) error {
	return nil
}
func (m *mockContainerRepo) RemoveContainerByUsergroupId(groupId int) error {
	return nil
}

type mockUserRepo struct{}

func (m *mockUserRepo) GetAllUsers() ([]*user.User, error) {
	return []*user.User{}, nil
}
