package repository

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/happYness-Project/taskManagementGolang/internal/taskcontainer/domain"
	"github.com/stretchr/testify/require"
)

func TestContainerRepo_AllTaskContainers(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()
	containerRepo := NewContainerRepository(db)

	t.Run("Container Exists", func(t *testing.T) {
		mockContainer := mockContainerObj()
		rows := mockContainerRows(mockContainer)
		mock.ExpectQuery(sqlGetAllContainers).WillReturnRows(rows)

		containers, err := containerRepo.AllTaskContainers()

		require.Nil(t, err)
		require.Equal(t, mockContainer, *containers[0])
	})

	t.Run("Container does not exist, then should return nothing with no result error msg", func(t *testing.T) {
		mock.ExpectQuery(sqlGetAllContainers).WillReturnError(sql.ErrNoRows)

		containers, err := containerRepo.AllTaskContainers()

		require.NotNil(t, err)
		require.Equal(t, "sql: no rows in result set", err.Error())
		require.Len(t, containers, 0)
	})
}
func TestContainerRepo_ContainerById(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()
	containerRepo := NewContainerRepository(db)
	t.Run("Container Exists", func(t *testing.T) {
		mockContainer := mockContainerObj()
		rows := mockContainerRows(mockContainer)
		mock.ExpectQuery(sqlGetById).
			WithArgs(mockContainer.Id).
			WillReturnRows(rows)

		container, err := containerRepo.GetById(mockContainer.Id)

		require.Nil(t, err)
		require.Equal(t, &mockContainer, container)
	})

	t.Run("When Cannot find Container, Then return null", func(t *testing.T) {
		mock.ExpectQuery(sqlGetById).
			WithArgs("invalid")

		container, err := containerRepo.GetById("invalid")

		require.NotNil(t, err)
		require.Nil(t, container)
	})

	t.Run("Internal server error.", func(t *testing.T) {
		mock.ExpectQuery(sqlGetById).
			WithArgs("SomeError").WillReturnError(errors.New("random internal error occurred"))

		container, err := containerRepo.GetById("invalid")
		require.NotNil(t, err)
		require.Nil(t, container)
		// TODO Need to update GetById to handle different scenario.
	})
}
func TestContainerRepo_GetContainersByGroupId(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()
	containerRepo := NewContainerRepository(db)

	t.Run("Containers exist for the given group id", func(t *testing.T) {
		mockContainer := mockContainerObj()
		rows := mockContainerRows(mockContainer)
		mock.ExpectQuery(sqlGetContainersByGroupId).
			WithArgs(mockContainer.UsergroupId).
			WillReturnRows(rows)

		_, err := containerRepo.GetContainersByGroupId(mockContainer.UsergroupId)

		require.Nil(t, err)
	})
}

func mockContainerObj() domain.TaskContainer {
	return domain.TaskContainer{
		Id:          uuid.New().String(),
		Name:        "testuser",
		Description: "testdesc",
		IsActive:    true,
		UsergroupId: 1,
	}
}
func mockContainerRows(c domain.TaskContainer) *sqlmock.Rows {
	return sqlmock.NewRows([]string{"id", "name", "description", "is_active", "usergroup_id"}).
		AddRow(c.Id, c.Name, c.Description, c.IsActive, c.UsergroupId)
}
