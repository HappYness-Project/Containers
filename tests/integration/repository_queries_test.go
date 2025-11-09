package integration

import (
	"testing"

	"github.com/happYness-Project/taskManagementGolang/tests/builders"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRepository_GetUsersByGroupIdWithRoles(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	cleanup := setupTest(t)
	defer cleanup()

	// Create multiple users
	admin := builders.NewUserBuilder().
		WithUserName("admin").
		WithEmail("admin@test.com").
		Build()
	member1 := builders.NewUserBuilder().
		WithUserName("member1").
		WithEmail("member1@test.com").
		Build()
	member2 := builders.NewUserBuilder().
		WithUserName("member2").
		WithEmail("member2@test.com").
		Build()

	repos.UserRepo.CreateUser(*admin)
	repos.UserRepo.CreateUser(*member1)
	repos.UserRepo.CreateUser(*member2)

	adminFromDB, _ := repos.UserRepo.GetUserByUserId(admin.UserId)
	member1FromDB, _ := repos.UserRepo.GetUserByUserId(member1.UserId)
	member2FromDB, _ := repos.UserRepo.GetUserByUserId(member2.UserId)

	// Create group with admin
	group := builders.NewUserGroupBuilder().WithName("Dev Team").MustBuild()
	groupId, _ := repos.UserGroupRepo.CreateGroupWithUsers(*group, adminFromDB.Id)

	// Add members
	repos.UserGroupRepo.InsertUserGroupUserTable(groupId, member1FromDB.Id)
	repos.UserGroupRepo.InsertUserGroupUserTable(groupId, member2FromDB.Id)

	// Promote member1 to admin
	repos.UserGroupRepo.UpdateUserRoleInGroup(groupId, member1FromDB.Id, "admin")

	// Test: Get users with roles
	usersWithRoles, err := repos.UserRepo.GetUsersByGroupIdWithRoles(groupId)
	require.NoError(t, err)
	require.Len(t, usersWithRoles, 3)

	// Verify roles
	roleMap := make(map[string]string)
	for _, u := range usersWithRoles {
		roleMap[u.UserName] = u.Role
		// Verify UserWithRole has all user fields
		assert.NotEmpty(t, u.UserId)
		assert.NotEmpty(t, u.Email)
		assert.NotZero(t, u.JoinedAt)
	}

	assert.Equal(t, "admin", roleMap["admin"])
	assert.Equal(t, "admin", roleMap["member1"])
	assert.Equal(t, "member", roleMap["member2"])
}

func TestRepository_GetAllTasksByGroupId(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	cleanup := setupTest(t)
	defer cleanup()

	// Setup: Create user and two groups
	user := builders.NewUserBuilder().Build()
	repos.UserRepo.CreateUser(*user)
	userFromDB, _ := repos.UserRepo.GetUserByUserId(user.UserId)

	group1 := builders.NewUserGroupBuilder().WithName("Group 1").MustBuild()
	group2 := builders.NewUserGroupBuilder().WithName("Group 2").MustBuild()

	groupId1, _ := repos.UserGroupRepo.CreateGroupWithUsers(*group1, userFromDB.Id)
	groupId2, _ := repos.UserGroupRepo.CreateGroupWithUsers(*group2, userFromDB.Id)

	// Create containers in each group
	container1A := builders.NewTaskContainerBuilder().
		WithName("Container 1A").
		WithUsergroupId(groupId1).
		Build()
	container1B := builders.NewTaskContainerBuilder().
		WithName("Container 1B").
		WithUsergroupId(groupId1).
		Build()
	container2A := builders.NewTaskContainerBuilder().
		WithName("Container 2A").
		WithUsergroupId(groupId2).
		Build()

	repos.TaskContainerRepo.CreateContainer(*container1A)
	repos.TaskContainerRepo.CreateContainer(*container1B)
	repos.TaskContainerRepo.CreateContainer(*container2A)

	// Create tasks in different containers
	task1 := builders.NewTaskBuilder().WithName("Task in 1A").MustBuild()
	task2 := builders.NewTaskBuilder().WithName("Task in 1B").MustBuild()
	task3 := builders.NewTaskBuilder().WithName("Another in 1A").MustBuild()
	task4 := builders.NewTaskBuilder().WithName("Task in 2A").MustBuild()

	repos.TaskRepo.CreateTask(container1A.Id, *task1)
	repos.TaskRepo.CreateTask(container1B.Id, *task2)
	repos.TaskRepo.CreateTask(container1A.Id, *task3)
	repos.TaskRepo.CreateTask(container2A.Id, *task4)

	// Test: Get all tasks in group 1
	group1Tasks, err := repos.TaskRepo.GetAllTasksByGroupId(groupId1)
	require.NoError(t, err)
	assert.Len(t, group1Tasks, 3)

	// Test: Get all tasks in group 2
	group2Tasks, err := repos.TaskRepo.GetAllTasksByGroupId(groupId2)
	require.NoError(t, err)
	assert.Len(t, group2Tasks, 1)

	// Verify task names in group 1
	taskNames := make(map[string]bool)
	for _, task := range group1Tasks {
		taskNames[task.TaskName] = true
	}
	assert.True(t, taskNames["Task in 1A"])
	assert.True(t, taskNames["Task in 1B"])
	assert.True(t, taskNames["Another in 1A"])
	assert.False(t, taskNames["Task in 2A"])
}

// TestRepository_GetUserRoleInGroup tests the role lookup query
func TestRepository_GetUserRoleInGroup(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	cleanup := setupTest(t)
	defer cleanup()

	// Create user and group
	user := builders.NewUserBuilder().Build()
	repos.UserRepo.CreateUser(*user)
	userFromDB, _ := repos.UserRepo.GetUserByUserId(user.UserId)

	group := builders.NewUserGroupBuilder().WithName("Test Group").MustBuild()
	groupId, _ := repos.UserGroupRepo.CreateGroupWithUsers(*group, userFromDB.Id)

	// Test: Get role (should be admin)
	role, err := repos.UserRepo.GetUserRoleInGroup(userFromDB.UserId, groupId)
	require.NoError(t, err)
	assert.Equal(t, "admin", role)

	// Test: Get role for non-existent user-group relationship
	_, err = repos.UserRepo.GetUserRoleInGroup("non-existent-user", groupId)
	assert.Error(t, err)
}

// TestRepository_GetUsersByGroupId tests basic user retrieval by group
func TestRepository_GetUsersByGroupId(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	cleanup := setupTest(t)
	defer cleanup()

	// Create users
	users := []*builders.UserBuilder{
		builders.NewUserBuilder().WithUserName("user1"),
		builders.NewUserBuilder().WithUserName("user2"),
		builders.NewUserBuilder().WithUserName("user3"),
	}

	var userIds []int
	for _, userBuilder := range users {
		user := userBuilder.Build()
		repos.UserRepo.CreateUser(*user)
		userFromDB, _ := repos.UserRepo.GetUserByUserId(user.UserId)
		userIds = append(userIds, userFromDB.Id)
	}

	// Create group with first user as admin
	group := builders.NewUserGroupBuilder().WithName("Team").MustBuild()
	groupId, _ := repos.UserGroupRepo.CreateGroupWithUsers(*group, userIds[0])

	// Add other users
	repos.UserGroupRepo.InsertUserGroupUserTable(groupId, userIds[1])
	repos.UserGroupRepo.InsertUserGroupUserTable(groupId, userIds[2])

	// Test: Get all users in group
	groupUsers, err := repos.UserRepo.GetUsersByGroupId(groupId)
	require.NoError(t, err)
	assert.Len(t, groupUsers, 3)
}

// TestRepository_GetContainersByGroupId tests container retrieval by group
func TestRepository_GetContainersByGroupId(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	cleanup := setupTest(t)
	defer cleanup()

	// Setup
	user := builders.NewUserBuilder().Build()
	repos.UserRepo.CreateUser(*user)
	userFromDB, _ := repos.UserRepo.GetUserByUserId(user.UserId)

	group1 := builders.NewUserGroupBuilder().WithName("Group 1").MustBuild()
	group2 := builders.NewUserGroupBuilder().WithName("Group 2").MustBuild()

	groupId1, _ := repos.UserGroupRepo.CreateGroupWithUsers(*group1, userFromDB.Id)
	groupId2, _ := repos.UserGroupRepo.CreateGroupWithUsers(*group2, userFromDB.Id)

	// Create containers
	repos.TaskContainerRepo.CreateContainer(*builders.NewTaskContainerBuilder().
		WithName("Container A").
		WithUsergroupId(groupId1).
		Build())
	repos.TaskContainerRepo.CreateContainer(*builders.NewTaskContainerBuilder().
		WithName("Container B").
		WithUsergroupId(groupId1).
		Build())
	repos.TaskContainerRepo.CreateContainer(*builders.NewTaskContainerBuilder().
		WithName("Container C").
		WithUsergroupId(groupId2).
		Build())

	// Test: Get containers for group 1
	containers1, err := repos.TaskContainerRepo.GetContainersByGroupId(groupId1)
	require.NoError(t, err)
	assert.Len(t, containers1, 2)

	// Test: Get containers for group 2
	containers2, err := repos.TaskContainerRepo.GetContainersByGroupId(groupId2)
	require.NoError(t, err)
	assert.Len(t, containers2, 1)
}

func TestRepository_UserLookupMethods(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	cleanup := setupTest(t)
	defer cleanup()

	// Create user
	user := builders.NewUserBuilder().
		WithUserName("testuser").
		WithEmail("test@example.com").
		Build()
	repos.UserRepo.CreateUser(*user)

	// Test: Get by UserId
	byUserId, err := repos.UserRepo.GetUserByUserId(user.UserId)
	require.NoError(t, err)
	assert.Equal(t, user.UserId, byUserId.UserId)

	// Test: Get by Username
	byUsername, err := repos.UserRepo.GetUserByUsername("testuser")
	require.NoError(t, err)
	assert.Equal(t, "testuser", byUsername.UserName)

	// Test: Get by Email
	byEmail, err := repos.UserRepo.GetUserByEmail("test@example.com")
	require.NoError(t, err)
	assert.Equal(t, "test@example.com", byEmail.Email)

	// Test: Get all users
	allUsers, err := repos.UserRepo.GetAllUsers()
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(allUsers), 1)
}

func TestRepository_GetAllTasksByGroupIdOnlyImportant(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	cleanup := setupTest(t)
	defer cleanup()

	// Setup
	user := builders.NewUserBuilder().Build()
	repos.UserRepo.CreateUser(*user)
	userFromDB, _ := repos.UserRepo.GetUserByUserId(user.UserId)

	group := builders.NewUserGroupBuilder().WithName("Project").MustBuild()
	groupId, _ := repos.UserGroupRepo.CreateGroupWithUsers(*group, userFromDB.Id)

	container := builders.NewTaskContainerBuilder().
		WithUsergroupId(groupId).
		Build()
	repos.TaskContainerRepo.CreateContainer(*container)

	// Create tasks with different importance
	task1 := builders.NewTaskBuilder().WithName("Important Task 1").MustBuild()
	task2 := builders.NewTaskBuilder().WithName("Normal Task").MustBuild()
	task3 := builders.NewTaskBuilder().WithName("Important Task 2").MustBuild()

	repos.TaskRepo.CreateTask(container.Id, *task1)
	repos.TaskRepo.CreateTask(container.Id, *task2)
	repos.TaskRepo.CreateTask(container.Id, *task3)

	// Mark tasks as important
	repos.TaskRepo.UpdateImportantTask(task1.TaskId, true)
	repos.TaskRepo.UpdateImportantTask(task3.TaskId, true)

	// Test: Get only important tasks
	importantTasks, err := repos.TaskRepo.GetAllTasksByGroupIdOnlyImportant(groupId)
	require.NoError(t, err)
	assert.Len(t, importantTasks, 2)

	// Verify all returned tasks are important
	for _, task := range importantTasks {
		assert.True(t, task.IsImportant)
	}
}
