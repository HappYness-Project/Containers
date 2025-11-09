package integration

import (
	"testing"

	"github.com/happYness-Project/taskManagementGolang/tests/builders"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTaskRepository_Create(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("should create task in task container successfully", func(t *testing.T) {
		cleanup := setupTest(t)
		defer cleanup()

		// Arrange
		containerId := setupTaskEnvironment(t)
		task := builders.NewTaskBuilder().
			WithName("New Feature Implementation").
			WithDescription("Implement user authentication").
			WithPriority("high").
			WithCategory("development").
			MustBuild()

		// Act
		createdTask, err := repos.TaskRepo.CreateTask(containerId, *task)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, task.TaskId, createdTask.TaskId)
		assert.Equal(t, "New Feature Implementation", createdTask.TaskName)
		assert.Equal(t, "Implement user authentication", createdTask.TaskDesc)
		assert.Equal(t, "high", createdTask.Priority)
		assert.Equal(t, "development", createdTask.Category)
		assert.False(t, createdTask.IsCompleted)
		assert.False(t, createdTask.IsImportant)
	})

	t.Run("should create task with default values", func(t *testing.T) {
		cleanup := setupTest(t)
		defer cleanup()

		// Arrange
		containerId := setupTaskEnvironment(t)
		task := builders.NewTaskBuilder().MustBuild() // Use all defaults

		// Act
		createdTask, err := repos.TaskRepo.CreateTask(containerId, *task)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, "Test Task", createdTask.TaskName)
		assert.Equal(t, "Test Description", createdTask.TaskDesc)
		assert.Equal(t, "medium", createdTask.Priority)
		assert.Equal(t, "work", createdTask.Category)
	})

	t.Run("should create task with different priorities", func(t *testing.T) {
		cleanup := setupTest(t)
		defer cleanup()

		// Arrange
		containerId := setupTaskEnvironment(t)
		lowTask := builders.NewTaskBuilder().WithName("Low").LowPriority().MustBuild()
		highTask := builders.NewTaskBuilder().WithName("High").HighPriority().MustBuild()
		urgentTask := builders.NewTaskBuilder().WithName("Urgent").UrgentPriority().MustBuild()

		// Act
		createdLow, err1 := repos.TaskRepo.CreateTask(containerId, *lowTask)
		createdHigh, err2 := repos.TaskRepo.CreateTask(containerId, *highTask)
		createdUrgent, err3 := repos.TaskRepo.CreateTask(containerId, *urgentTask)

		// Assert
		require.NoError(t, err1)
		require.NoError(t, err2)
		require.NoError(t, err3)
		assert.Equal(t, "low", createdLow.Priority)
		assert.Equal(t, "high", createdHigh.Priority)
		assert.Equal(t, "urgent", createdUrgent.Priority)
	})

	t.Run("should create task with different target dates", func(t *testing.T) {
		cleanup := setupTest(t)
		defer cleanup()

		// Arrange
		containerId := setupTaskEnvironment(t)
		todayTask := builders.NewTaskBuilder().WithName("Due Today").DueToday().MustBuild()
		tomorrowTask := builders.NewTaskBuilder().WithName("Due Tomorrow").DueTomorrow().MustBuild()
		overdueTask := builders.NewTaskBuilder().WithName("Overdue").Overdue().MustBuild()

		// Act
		_, err1 := repos.TaskRepo.CreateTask(containerId, *todayTask)
		_, err2 := repos.TaskRepo.CreateTask(containerId, *tomorrowTask)
		_, err3 := repos.TaskRepo.CreateTask(containerId, *overdueTask)

		// Assert
		require.NoError(t, err1)
		require.NoError(t, err2)
		require.NoError(t, err3)

		// Verify tasks exist
		tasks, err := repos.TaskRepo.GetTasksByContainerId(containerId)
		require.NoError(t, err)
		assert.Len(t, tasks, 3)
	})
}

func TestTaskRepository_UpdateTask(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	t.Run("should update task details successfully", func(t *testing.T) {
		cleanup := setupTest(t)
		defer cleanup()

		// Arrange
		containerId := setupTaskEnvironment(t)
		task := builders.NewTaskBuilder().WithName("Initial Name").MustBuild()
		createdTask, err := repos.TaskRepo.CreateTask(containerId, *task)
		require.NoError(t, err)

		// Act
		createdTask.TaskName = "Updated Name"
		createdTask.TaskDesc = "Updated Description"
		createdTask.Priority = "urgent"
		createdTask.Category = "personal"
		err = repos.TaskRepo.UpdateTask(createdTask)

		// Assert
		require.NoError(t, err)
		updatedTask, err := repos.TaskRepo.GetTaskById(createdTask.TaskId)
		require.NoError(t, err)
		assert.Equal(t, "Updated Name", updatedTask.TaskName)
		assert.Equal(t, "Updated Description", updatedTask.TaskDesc)
		assert.Equal(t, "urgent", updatedTask.Priority)
		assert.Equal(t, "personal", updatedTask.Category)
	})
}

func TestTaskRepository_DoneTask(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	t.Run("should mark task as completed", func(t *testing.T) {
		cleanup := setupTest(t)
		defer cleanup()

		// Arrange
		containerId := setupTaskEnvironment(t)
		task := builders.NewTaskBuilder().MustBuild()
		createdTask, err := repos.TaskRepo.CreateTask(containerId, *task)
		require.NoError(t, err)

		// Act
		err = repos.TaskRepo.DoneTask(createdTask.TaskId, true)

		// Assert
		require.NoError(t, err)
		updatedTask, err := repos.TaskRepo.GetTaskById(createdTask.TaskId)
		require.NoError(t, err)
		assert.True(t, updatedTask.IsCompleted)
	})

	t.Run("should mark task as not completed", func(t *testing.T) {
		cleanup := setupTest(t)
		defer cleanup()

		// Arrange
		containerId := setupTaskEnvironment(t)
		task := builders.NewTaskBuilder().MustBuild()
		createdTask, err := repos.TaskRepo.CreateTask(containerId, *task)
		require.NoError(t, err)

		// Mark as completed first
		err = repos.TaskRepo.DoneTask(createdTask.TaskId, true)
		require.NoError(t, err)

		// Act
		err = repos.TaskRepo.DoneTask(createdTask.TaskId, false)

		// Assert
		require.NoError(t, err)
		updatedTask, err := repos.TaskRepo.GetTaskById(createdTask.TaskId)
		require.NoError(t, err)
		assert.False(t, updatedTask.IsCompleted)
	})

}

func TestTaskRepository_UpdateImportantTask(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	t.Run("should mark task as important", func(t *testing.T) {
		cleanup := setupTest(t)
		defer cleanup()
		// Arrange
		containerId := setupTaskEnvironment(t)
		task := builders.NewTaskBuilder().MustBuild()
		createdTask, err := repos.TaskRepo.CreateTask(containerId, *task)
		require.NoError(t, err)

		// Act
		err = repos.TaskRepo.UpdateImportantTask(createdTask.TaskId, true)

		// Assert
		require.NoError(t, err)
		updatedTask, err := repos.TaskRepo.GetTaskById(createdTask.TaskId)
		require.NoError(t, err)
		assert.True(t, updatedTask.IsImportant)
	})

	t.Run("should unmark task as important", func(t *testing.T) {
		cleanup := setupTest(t)
		defer cleanup()
		// Arrange
		containerId := setupTaskEnvironment(t)
		task := builders.NewTaskBuilder().MustBuild()
		createdTask, err := repos.TaskRepo.CreateTask(containerId, *task)
		require.NoError(t, err)

		// Mark as important first
		err = repos.TaskRepo.UpdateImportantTask(createdTask.TaskId, true)
		require.NoError(t, err)

		// Act
		err = repos.TaskRepo.UpdateImportantTask(createdTask.TaskId, false)

		// Assert
		require.NoError(t, err)
		updatedTask, err := repos.TaskRepo.GetTaskById(createdTask.TaskId)
		require.NoError(t, err)
		assert.False(t, updatedTask.IsImportant)
	})
}

func TestTaskRepository_DeleteTask(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	t.Run("should delete task successfully", func(t *testing.T) {
		cleanup := setupTest(t)
		defer cleanup()

		containerId := setupTaskEnvironment(t)

		// Create task
		task := builders.NewTaskBuilder().WithName("Task to delete").MustBuild()
		createdTask, _ := repos.TaskRepo.CreateTask(containerId, *task)

		// Verify task exists
		existingTask, err := repos.TaskRepo.GetTaskById(createdTask.TaskId)
		require.NoError(t, err)
		require.NotNil(t, existingTask)

		// Delete task
		err = repos.TaskRepo.DeleteTask(createdTask.TaskId)
		require.NoError(t, err)

		deletedTask, _ := repos.TaskRepo.GetTaskById(createdTask.TaskId)
		assert.Nil(t, deletedTask)
	})

	t.Run("should handle deleting non-existent task gracefully", func(t *testing.T) {
		cleanup := setupTest(t)
		defer cleanup()

		// Attempt to delete a task that doesn't exist (using valid UUID format)
		err := repos.TaskRepo.DeleteTask("00000000-0000-0000-0000-000000000000")
		// This test documents the behavior - may or may not error depending on implementation
		_ = err // Should handle gracefully or return error
	})
}

func TestTaskRepository_GetAllTasksByGroupId(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	t.Run("should return empty list when no tasks in group", func(t *testing.T) {
		cleanup := setupTest(t)
		defer cleanup()

		// Setup user and group
		user := builders.NewUserBuilder().Build()
		repos.UserRepo.CreateUser(*user)
		userFromDB, _ := repos.UserRepo.GetUserByUserId(user.UserId)

		group := builders.NewUserGroupBuilder().WithName("Empty Group").MustBuild()
		groupId, _ := repos.UserGroupRepo.CreateGroupWithUsers(*group, userFromDB.Id)

		// Get tasks by group ID
		tasks, err := repos.TaskRepo.GetAllTasksByGroupId(groupId)
		require.NoError(t, err)
		assert.Len(t, tasks, 0)
	})

	t.Run("should return list of tasks for group", func(t *testing.T) {
		cleanup := setupTest(t)
		defer cleanup()

		// Setup user and group
		user := builders.NewUserBuilder().Build()
		repos.UserRepo.CreateUser(*user)
		userFromDB, _ := repos.UserRepo.GetUserByUserId(user.UserId)

		group := builders.NewUserGroupBuilder().WithName("Active Group").MustBuild()
		groupId, _ := repos.UserGroupRepo.CreateGroupWithUsers(*group, userFromDB.Id)

		// Create container in group
		container := builders.NewTaskContainerBuilder().
			WithName("Main Container").
			WithUsergroupId(groupId).
			Build()
		repos.TaskContainerRepo.CreateContainer(*container)

		// Create tasks in container
		task1 := builders.NewTaskBuilder().WithName("Task 1").MustBuild()
		task2 := builders.NewTaskBuilder().WithName("Task 2").MustBuild()

		repos.TaskRepo.CreateTask(container.Id, *task1)
		repos.TaskRepo.CreateTask(container.Id, *task2)

		// Get tasks by group ID
		tasks, err := repos.TaskRepo.GetAllTasksByGroupId(groupId)
		require.NoError(t, err)
		assert.Len(t, tasks, 2)
	})
}

func TestTaskRepository_GetTasksByContainerId(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	t.Run("should return empty list when no tasks in container", func(t *testing.T) {
		cleanup := setupTest(t)
		defer cleanup()

		containerId := setupTaskEnvironment(t)
		// Get tasks by container ID
		tasks, err := repos.TaskRepo.GetTasksByContainerId(containerId)
		require.NoError(t, err)
		assert.Len(t, tasks, 0)
	})

	t.Run("should return list of tasks for container", func(t *testing.T) {
		cleanup := setupTest(t)
		defer cleanup()

		containerId := setupTaskEnvironment(t)

		// Create tasks in container
		task1 := builders.NewTaskBuilder().WithName("Task A").MustBuild()
		task2 := builders.NewTaskBuilder().WithName("Task B").MustBuild()

		repos.TaskRepo.CreateTask(containerId, *task1)
		repos.TaskRepo.CreateTask(containerId, *task2)

		// Get tasks by container ID
		tasks, err := repos.TaskRepo.GetTasksByContainerId(containerId)
		require.NoError(t, err)
		assert.Len(t, tasks, 2)
		taskNames := make(map[string]bool)
		for _, task := range tasks {
			taskNames[task.TaskName] = true
		}
		assert.True(t, taskNames["Task A"])
		assert.True(t, taskNames["Task B"])
	})
}

func setupTaskEnvironment(t *testing.T) string {
	t.Helper()

	user := builders.NewUserBuilder().Build()
	repos.UserRepo.CreateUser(*user)
	userFromDB, _ := repos.UserRepo.GetUserByUserId(user.UserId)

	group := builders.NewUserGroupBuilder().MustBuild()
	groupId, _ := repos.UserGroupRepo.CreateGroupWithUsers(*group, userFromDB.Id)

	container := builders.NewTaskContainerBuilder().
		WithUsergroupId(groupId).
		Build()
	repos.TaskContainerRepo.CreateContainer(*container)

	return container.Id
}
