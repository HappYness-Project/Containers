package integration

import (
	"testing"
	"time"

	"github.com/happYness-Project/taskManagementGolang/internal/task/repository"
	taskcontainerRepo "github.com/happYness-Project/taskManagementGolang/internal/taskcontainer/repository"
	userRepo "github.com/happYness-Project/taskManagementGolang/internal/user/repository"
	usergroupRepo "github.com/happYness-Project/taskManagementGolang/internal/usergroup/repository"
	"github.com/happYness-Project/taskManagementGolang/tests/builders"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTaskFlow_CreateTaskInContainer(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	cleanup := setupTest(t)
	defer cleanup()

	// Setup: Create user, group, and container
	userRepository := userRepo.NewUserRepository(testDB)
	groupRepository := usergroupRepo.NewUserGroupRepository(testDB)
	containerRepository := taskcontainerRepo.NewContainerRepository(testDB)
	taskRepository := repository.NewTaskRepository(testDB)

	user := builders.NewUserBuilder().Build()
	userRepository.CreateUser(*user)
	userFromDB, _ := userRepository.GetUserByUserId(user.UserId)

	group := builders.NewUserGroupBuilder().WithName("My Tasks").MustBuild()
	groupId, _ := groupRepository.CreateGroupWithUsers(*group, userFromDB.Id)

	container := builders.NewTaskContainerBuilder().
		WithName("TODO List").
		WithUsergroupId(groupId).
		Build()
	err := containerRepository.CreateContainer(*container)
	require.NoError(t, err)

	// Create a task
	task := builders.NewTaskBuilder().
		WithName("Write integration tests").
		WithDescription("Create comprehensive test suite").
		WithPriority("high").
		WithCategory("development").
		MustBuild()

	createdTask, err := taskRepository.CreateTask(container.Id, *task)
	require.NoError(t, err)
	assert.Equal(t, "Write integration tests", createdTask.TaskName)

	// Verify task was created
	retrievedTask, err := taskRepository.GetTaskById(task.TaskId)
	require.NoError(t, err)
	assert.Equal(t, task.TaskId, retrievedTask.TaskId)
	assert.Equal(t, "Write integration tests", retrievedTask.TaskName)
	assert.Equal(t, "high", retrievedTask.Priority)
	assert.False(t, retrievedTask.IsCompleted)
	assert.False(t, retrievedTask.IsImportant)

	// Verify task is in container
	tasksInContainer, err := taskRepository.GetTasksByContainerId(container.Id)
	require.NoError(t, err)
	assert.Len(t, tasksInContainer, 1)
	assert.Equal(t, task.TaskId, tasksInContainer[0].TaskId)
}

func TestTaskFlow_UpdateTask(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	cleanup := setupTest(t)
	defer cleanup()

	// Setup
	taskRepository, containerId := setupTaskEnvironment(t)

	// Create task
	task := builders.NewTaskBuilder().
		WithName("Original Name").
		WithPriority("low").
		MustBuild()
	createdTask, _ := taskRepository.CreateTask(containerId, *task)

	// Update task
	err := createdTask.UpdateTask(
		"Updated Name",
		"Updated Description",
		time.Now().Add(48*time.Hour),
		"urgent",
		"personal",
	)
	require.NoError(t, err)

	err = taskRepository.UpdateTask(createdTask)
	require.NoError(t, err)

	// Verify update
	updatedTask, err := taskRepository.GetTaskById(createdTask.TaskId)
	require.NoError(t, err)
	assert.Equal(t, "Updated Name", updatedTask.TaskName)
	assert.Equal(t, "Updated Description", updatedTask.TaskDesc)
	assert.Equal(t, "urgent", updatedTask.Priority)
	assert.Equal(t, "personal", updatedTask.Category)
}

func TestTaskFlow_ToggleCompletion(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	cleanup := setupTest(t)
	defer cleanup()

	taskRepository, containerId := setupTaskEnvironment(t)

	// Create task
	task := builders.NewTaskBuilder().WithName("Task to complete").MustBuild()
	createdTask, _ := taskRepository.CreateTask(containerId, *task)

	// Mark as completed
	err := taskRepository.DoneTask(createdTask.TaskId, true)
	require.NoError(t, err)

	// Verify completion
	completedTask, _ := taskRepository.GetTaskById(createdTask.TaskId)
	assert.True(t, completedTask.IsCompleted)

	// Mark as incomplete
	err = taskRepository.DoneTask(createdTask.TaskId, false)
	require.NoError(t, err)

	// Verify incomplete
	incompleteTask, _ := taskRepository.GetTaskById(createdTask.TaskId)
	assert.False(t, incompleteTask.IsCompleted)
}

func TestTaskFlow_ToggleImportant(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	cleanup := setupTest(t)
	defer cleanup()

	taskRepository, containerId := setupTaskEnvironment(t)

	// Create task
	task := builders.NewTaskBuilder().WithName("Important task").MustBuild()
	createdTask, _ := taskRepository.CreateTask(containerId, *task)

	// Mark as important
	err := taskRepository.UpdateImportantTask(createdTask.TaskId, true)
	require.NoError(t, err)

	// Verify importance
	importantTask, _ := taskRepository.GetTaskById(createdTask.TaskId)
	assert.True(t, importantTask.IsImportant)

	// Mark as not important
	err = taskRepository.UpdateImportantTask(createdTask.TaskId, false)
	require.NoError(t, err)

	// Verify not important
	normalTask, _ := taskRepository.GetTaskById(createdTask.TaskId)
	assert.False(t, normalTask.IsImportant)
}

func TestTaskFlow_DeleteTask(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	cleanup := setupTest(t)
	defer cleanup()

	taskRepository, containerId := setupTaskEnvironment(t)

	// Create task
	task := builders.NewTaskBuilder().WithName("Task to delete").MustBuild()
	createdTask, _ := taskRepository.CreateTask(containerId, *task)

	// Verify task exists
	existingTask, err := taskRepository.GetTaskById(createdTask.TaskId)
	require.NoError(t, err)
	require.NotNil(t, existingTask)

	// Delete task
	err = taskRepository.DeleteTask(createdTask.TaskId)
	require.NoError(t, err)

	// Verify task is deleted
	deletedTask, _ := taskRepository.GetTaskById(createdTask.TaskId)
	assert.Nil(t, deletedTask.TaskId) // Should return empty task
}

func TestTaskFlow_GetTasksByGroupId(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	cleanup := setupTest(t)
	defer cleanup()

	// Setup multiple containers in same group
	userRepository := userRepo.NewUserRepository(testDB)
	groupRepository := usergroupRepo.NewUserGroupRepository(testDB)
	containerRepository := taskcontainerRepo.NewContainerRepository(testDB)
	taskRepository := repository.NewTaskRepository(testDB)

	user := builders.NewUserBuilder().Build()
	userRepository.CreateUser(*user)
	userFromDB, _ := userRepository.GetUserByUserId(user.UserId)

	group := builders.NewUserGroupBuilder().WithName("Project").MustBuild()
	groupId, _ := groupRepository.CreateGroupWithUsers(*group, userFromDB.Id)

	// Create two containers
	container1 := builders.NewTaskContainerBuilder().
		WithName("Backlog").
		WithUsergroupId(groupId).
		Build()
	container2 := builders.NewTaskContainerBuilder().
		WithName("In Progress").
		WithUsergroupId(groupId).
		Build()

	containerRepository.CreateContainer(*container1)
	containerRepository.CreateContainer(*container2)

	// Create tasks in both containers
	task1 := builders.NewTaskBuilder().WithName("Task 1").MustBuild()
	task2 := builders.NewTaskBuilder().WithName("Task 2").MustBuild()
	task3 := builders.NewTaskBuilder().WithName("Task 3").Important().MustBuild()

	taskRepository.CreateTask(container1.Id, *task1)
	taskRepository.CreateTask(container1.Id, *task2)
	taskRepository.CreateTask(container2.Id, *task3)

	// Mark task3 as important in DB
	taskRepository.UpdateImportantTask(task3.TaskId, true)

	// Get all tasks by group
	allTasks, err := taskRepository.GetAllTasksByGroupId(groupId)
	require.NoError(t, err)
	assert.Len(t, allTasks, 3)

	// Get only important tasks
	importantTasks, err := taskRepository.GetAllTasksByGroupIdOnlyImportant(groupId)
	require.NoError(t, err)
	assert.Len(t, importantTasks, 1)
	assert.Equal(t, "Task 3", importantTasks[0].TaskName)
}

func TestTaskFlow_GetTasksByContainerId(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	cleanup := setupTest(t)
	defer cleanup()

	taskRepository, containerId := setupTaskEnvironment(t)

	// Create multiple tasks
	task1 := builders.NewTaskBuilder().WithName("Task 1").HighPriority().MustBuild()
	task2 := builders.NewTaskBuilder().WithName("Task 2").LowPriority().MustBuild()
	task3 := builders.NewTaskBuilder().WithName("Task 3").UrgentPriority().MustBuild()

	taskRepository.CreateTask(containerId, *task1)
	taskRepository.CreateTask(containerId, *task2)
	taskRepository.CreateTask(containerId, *task3)

	// Get all tasks in container
	tasks, err := taskRepository.GetTasksByContainerId(containerId)
	require.NoError(t, err)
	assert.Len(t, tasks, 3)

	// Verify task names
	taskNames := make(map[string]bool)
	for _, task := range tasks {
		taskNames[task.TaskName] = true
	}
	assert.True(t, taskNames["Task 1"])
	assert.True(t, taskNames["Task 2"])
	assert.True(t, taskNames["Task 3"])
}

func TestTaskFlow_MultipleContainersInGroup(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	cleanup := setupTest(t)
	defer cleanup()

	userRepository := userRepo.NewUserRepository(testDB)
	groupRepository := usergroupRepo.NewUserGroupRepository(testDB)
	containerRepository := taskcontainerRepo.NewContainerRepository(testDB)

	// Setup group
	user := builders.NewUserBuilder().Build()
	userRepository.CreateUser(*user)
	userFromDB, _ := userRepository.GetUserByUserId(user.UserId)

	group := builders.NewUserGroupBuilder().WithName("Team").MustBuild()
	groupId, _ := groupRepository.CreateGroupWithUsers(*group, userFromDB.Id)

	// Create multiple containers
	containers := []string{"TODO", "In Progress", "Done", "Backlog"}
	for _, name := range containers {
		container := builders.NewTaskContainerBuilder().
			WithName(name).
			WithUsergroupId(groupId).
			Build()
		err := containerRepository.CreateContainer(*container)
		require.NoError(t, err)
	}

	// Verify all containers are in group
	groupContainers, err := containerRepository.GetContainersByGroupId(groupId)
	require.NoError(t, err)
	assert.Len(t, groupContainers, 4)

	containerNames := make(map[string]bool)
	for _, container := range groupContainers {
		containerNames[container.Name] = true
	}
	assert.True(t, containerNames["TODO"])
	assert.True(t, containerNames["In Progress"])
	assert.True(t, containerNames["Done"])
	assert.True(t, containerNames["Backlog"])
}

// Helper function to setup a basic task environment (user, group, container)
func setupTaskEnvironment(t *testing.T) (*repository.TaskRepo, string) {
	t.Helper()

	userRepository := userRepo.NewUserRepository(testDB)
	groupRepository := usergroupRepo.NewUserGroupRepository(testDB)
	containerRepository := taskcontainerRepo.NewContainerRepository(testDB)
	taskRepository := repository.NewTaskRepository(testDB)

	user := builders.NewUserBuilder().Build()
	userRepository.CreateUser(*user)
	userFromDB, _ := userRepository.GetUserByUserId(user.UserId)

	group := builders.NewUserGroupBuilder().MustBuild()
	groupId, _ := groupRepository.CreateGroupWithUsers(*group, userFromDB.Id)

	container := builders.NewTaskContainerBuilder().
		WithUsergroupId(groupId).
		Build()
	containerRepository.CreateContainer(*container)

	return taskRepository, container.Id
}
