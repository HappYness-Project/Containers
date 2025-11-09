package integration

import (
	"testing"

	"github.com/happYness-Project/taskManagementGolang/internal/user/repository"
	usergroupRepo "github.com/happYness-Project/taskManagementGolang/internal/usergroup/repository"
	"github.com/happYness-Project/taskManagementGolang/tests/builders"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserGroupFlow_CreateGroupWithAdmin(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	cleanup := setupTest(t)
	defer cleanup()

	// Create repositories
	userRepo := repository.NewUserRepository(testDB)
	groupRepo := usergroupRepo.NewUserGroupRepository(testDB)

	// Create a user
	user := builders.NewUserBuilder().
		WithUserName("groupadmin").
		WithEmail("admin@example.com").
		Build()

	err := userRepo.CreateUser(*user)
	require.NoError(t, err)

	// Verify user was created
	createdUser, err := userRepo.GetUserByUserId(user.UserId)
	require.NoError(t, err)
	require.NotNil(t, createdUser)

	// Create a group with the user as admin
	group := builders.NewUserGroupBuilder().
		WithName("Development Team").
		WithDescription("Our dev team").
		TeamType().
		MustBuild()

	groupId, err := groupRepo.CreateGroupWithUsers(*group, createdUser.Id)
	require.NoError(t, err)
	require.Greater(t, groupId, 0)

	// Verify group was created
	createdGroup, err := groupRepo.GetById(groupId)
	require.NoError(t, err)
	assert.Equal(t, "Development Team", createdGroup.GroupName)
	assert.Equal(t, "Our dev team", createdGroup.GroupDesc)
	assert.Equal(t, "team", createdGroup.Type)
	assert.True(t, createdGroup.IsActive)

	// Verify user is assigned to group
	groups, err := groupRepo.GetUserGroupsByUserId(createdUser.Id)
	require.NoError(t, err)
	assert.Len(t, groups, 1)
	assert.Equal(t, groupId, groups[0].GroupId)

	// Verify user has admin role
	role, err := userRepo.GetUserRoleInGroup(createdUser.UserId, groupId)
	require.NoError(t, err)
	assert.Equal(t, "admin", role)
}

func TestUserGroupFlow_AddMemberToGroup(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	cleanup := setupTest(t)
	defer cleanup()

	userRepo := repository.NewUserRepository(testDB)
	groupRepo := usergroupRepo.NewUserGroupRepository(testDB)

	// Create admin user
	admin := builders.NewUserBuilder().
		WithUserName("admin").
		WithEmail("admin@example.com").
		Build()
	err := userRepo.CreateUser(*admin)
	require.NoError(t, err)
	adminFromDB, _ := userRepo.GetUserByUserId(admin.UserId)

	// Create member user
	member := builders.NewUserBuilder().
		WithUserName("member").
		WithEmail("member@example.com").
		Build()
	err = userRepo.CreateUser(*member)
	require.NoError(t, err)
	memberFromDB, _ := userRepo.GetUserByUserId(member.UserId)

	// Create group with admin
	group := builders.NewUserGroupBuilder().
		WithName("Project Alpha").
		MustBuild()
	groupId, err := groupRepo.CreateGroupWithUsers(*group, adminFromDB.Id)
	require.NoError(t, err)

	// Add member to group
	err = groupRepo.InsertUserGroupUserTable(groupId, memberFromDB.Id)
	require.NoError(t, err)

	// Verify both users are in the group
	usersWithRoles, err := userRepo.GetUsersByGroupIdWithRoles(groupId)
	require.NoError(t, err)
	assert.Len(t, usersWithRoles, 2)

	// Verify roles
	roles := make(map[string]string)
	for _, u := range usersWithRoles {
		roles[u.UserName] = u.Role
	}
	assert.Equal(t, "admin", roles["admin"])
	assert.Equal(t, "member", roles["member"])
}

func TestUserGroupFlow_ChangeUserRole(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	cleanup := setupTest(t)
	defer cleanup()

	userRepo := repository.NewUserRepository(testDB)
	groupRepo := usergroupRepo.NewUserGroupRepository(testDB)

	// Create users
	admin := builders.NewUserBuilder().WithUserName("admin").Build()
	member := builders.NewUserBuilder().WithUserName("member").Build()

	userRepo.CreateUser(*admin)
	userRepo.CreateUser(*member)

	adminFromDB, _ := userRepo.GetUserByUserId(admin.UserId)
	memberFromDB, _ := userRepo.GetUserByUserId(member.UserId)

	// Create group
	group := builders.NewUserGroupBuilder().WithName("Test Group").MustBuild()
	groupId, _ := groupRepo.CreateGroupWithUsers(*group, adminFromDB.Id)

	// Add member
	groupRepo.InsertUserGroupUserTable(groupId, memberFromDB.Id)

	// Promote member to admin
	err := groupRepo.UpdateUserRoleInGroup(groupId, memberFromDB.Id, "admin")
	require.NoError(t, err)

	// Verify role change
	role, err := userRepo.GetUserRoleInGroup(memberFromDB.UserId, groupId)
	require.NoError(t, err)
	assert.Equal(t, "admin", role)
}

func TestUserGroupFlow_RemoveMemberFromGroup(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	cleanup := setupTest(t)
	defer cleanup()

	userRepo := repository.NewUserRepository(testDB)
	groupRepo := usergroupRepo.NewUserGroupRepository(testDB)

	// Create users
	admin := builders.NewUserBuilder().WithUserName("admin").Build()
	member := builders.NewUserBuilder().WithUserName("member").Build()

	userRepo.CreateUser(*admin)
	userRepo.CreateUser(*member)

	adminFromDB, _ := userRepo.GetUserByUserId(admin.UserId)
	memberFromDB, _ := userRepo.GetUserByUserId(member.UserId)

	// Create group and add both users
	group := builders.NewUserGroupBuilder().WithName("Test Group").MustBuild()
	groupId, _ := groupRepo.CreateGroupWithUsers(*group, adminFromDB.Id)
	groupRepo.InsertUserGroupUserTable(groupId, memberFromDB.Id)

	// Verify 2 users in group
	usersBefore, _ := userRepo.GetUsersByGroupIdWithRoles(groupId)
	assert.Len(t, usersBefore, 2)

	// Remove member from group
	err := groupRepo.RemoveUserFromUserGroup(groupId, memberFromDB.Id)
	require.NoError(t, err)

	// Verify only 1 user remains
	usersAfter, _ := userRepo.GetUsersByGroupIdWithRoles(groupId)
	assert.Len(t, usersAfter, 1)
	assert.Equal(t, "admin", usersAfter[0].UserName)
}

func TestUserGroupFlow_GetUserGroups(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	cleanup := setupTest(t)
	defer cleanup()

	userRepo := repository.NewUserRepository(testDB)
	groupRepo := usergroupRepo.NewUserGroupRepository(testDB)

	// Create user
	user := builders.NewUserBuilder().Build()
	userRepo.CreateUser(*user)
	userFromDB, _ := userRepo.GetUserByUserId(user.UserId)

	// Create multiple groups
	group1 := builders.NewUserGroupBuilder().
		WithName("Team A").
		TeamType().
		MustBuild()
	group2 := builders.NewUserGroupBuilder().
		WithName("Project X").
		ProjectType().
		MustBuild()
	group3 := builders.NewUserGroupBuilder().
		WithName("Personal Tasks").
		PersonalType().
		MustBuild()

	groupId1, _ := groupRepo.CreateGroupWithUsers(*group1, userFromDB.Id)
	groupId2, _ := groupRepo.CreateGroupWithUsers(*group2, userFromDB.Id)
	groupId3, _ := groupRepo.CreateGroupWithUsers(*group3, userFromDB.Id)

	// Get all groups for user
	userGroups, err := groupRepo.GetUserGroupsByUserId(userFromDB.Id)
	require.NoError(t, err)
	assert.Len(t, userGroups, 3)

	// Verify group details
	groupNames := make(map[int]string)
	for _, g := range userGroups {
		groupNames[g.GroupId] = g.GroupName
	}
	assert.Equal(t, "Team A", groupNames[groupId1])
	assert.Equal(t, "Project X", groupNames[groupId2])
	assert.Equal(t, "Personal Tasks", groupNames[groupId3])
}

func TestUserGroupFlow_DeleteGroup(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	cleanup := setupTest(t)
	defer cleanup()

	userRepo := repository.NewUserRepository(testDB)
	groupRepo := usergroupRepo.NewUserGroupRepository(testDB)

	// Create user and group
	user := builders.NewUserBuilder().Build()
	userRepo.CreateUser(*user)
	userFromDB, _ := userRepo.GetUserByUserId(user.UserId)

	group := builders.NewUserGroupBuilder().WithName("Temporary Group").MustBuild()
	groupId, _ := groupRepo.CreateGroupWithUsers(*group, userFromDB.Id)

	// Verify group exists
	createdGroup, err := groupRepo.GetById(groupId)
	require.NoError(t, err)
	require.NotNil(t, createdGroup)

	// Delete group
	err = groupRepo.DeleteUserGroup(groupId)
	require.NoError(t, err)

	// Verify group is deleted (GetById should return empty or error)
	deletedGroup, err := groupRepo.GetById(groupId)
	// Group should either not be found or be inactive
	// This depends on your implementation (soft delete vs hard delete)
	if err == nil {
		assert.Equal(t, 0, deletedGroup.GroupId)
	}
}
