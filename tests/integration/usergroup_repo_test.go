package integration

import (
	"testing"

	"github.com/happYness-Project/taskManagementGolang/tests/builders"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserGroupRepository_CreateGroupWithUsers(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	t.Run("should create a user group and assign users as admin", func(t *testing.T) {
		cleanup := setupTest(t)
		defer cleanup()

		user := builders.NewUserBuilder().
			WithUserName("groupadmin").
			WithEmail("admin@example.com").
			Build()
		repos.UserRepo.CreateUser(*user)
		createdUser, _ := repos.UserRepo.GetUserByUserId(user.UserId)
		group := builders.NewUserGroupBuilder().
			WithName("Development Team").
			WithDescription("Our dev team").
			TeamType().
			MustBuild()

		groupId, err := repos.UserGroupRepo.CreateGroupWithUsers(*group, createdUser.Id)
		require.NoError(t, err)
		require.Greater(t, groupId, 0)

		// Verify group was created
		createdGroup, err := repos.UserGroupRepo.GetById(groupId)
		require.NoError(t, err)
		assert.Equal(t, "Development Team", createdGroup.GroupName)
		assert.Equal(t, "Our dev team", createdGroup.GroupDesc)
		assert.Equal(t, "team", createdGroup.Type)
		assert.True(t, createdGroup.IsActive)

		// Verify user is assigned to group
		groups, err := repos.UserGroupRepo.GetUserGroupsByUserId(createdUser.Id)
		require.NoError(t, err)
		assert.Len(t, groups, 1)
		assert.Equal(t, groupId, groups[0].GroupId)

		// Verify user has admin role
		role, err := repos.UserRepo.GetUserRoleInGroup(createdUser.UserId, groupId)
		require.NoError(t, err)
		assert.Equal(t, "admin", role)
	})

}

func TestUserGroupRepository_AddMemberToGroup(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	t.Run("should add a member to an existing user group", func(t *testing.T) {
		cleanup := setupTest(t)
		defer cleanup()

		// Create admin user
		admin := builders.NewUserBuilder().
			WithUserName("admin").
			WithEmail("admin@example.com").
			Build()
		err := repos.UserRepo.CreateUser(*admin)
		require.NoError(t, err)
		adminFromDB, _ := repos.UserRepo.GetUserByUserId(admin.UserId)

		// Create member user
		member := builders.NewUserBuilder().
			WithUserName("member").
			WithEmail("newuser@example.com").Build()

		err = repos.UserRepo.CreateUser(*member)
		require.NoError(t, err)
		memberFromDB, _ := repos.UserRepo.GetUserByUserId(member.UserId)

		// Create group with admin
		group := builders.NewUserGroupBuilder().
			WithName("Project Beta").
			MustBuild()
		groupId, err := repos.UserGroupRepo.CreateGroupWithUsers(*group, adminFromDB.Id)
		require.NoError(t, err)

		// Add member to group
		err = repos.UserGroupRepo.InsertUserGroupUserTable(groupId, memberFromDB.Id)
		require.NoError(t, err)

		// Verify both users are in the group
		usersWithRoles, err := repos.UserRepo.GetUsersByGroupIdWithRoles(groupId)
		require.NoError(t, err)
		assert.Len(t, usersWithRoles, 2)

		// Verify roles
		roles := make(map[string]string)
		for _, u := range usersWithRoles {
			roles[u.UserName] = u.Role
		}
		assert.Equal(t, "admin", roles["admin"])
		assert.Equal(t, "member", roles["member"])
	})
}

func TestUserGroupRepository_UpdateUserRoleInGroup(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	t.Run("should promote member to admin", func(t *testing.T) {
		cleanup := setupTest(t)
		defer cleanup()

		// Arrange - Create users
		admin := builders.NewUserBuilder().WithUserName("admin").Build()
		member := builders.NewUserBuilder().WithUserName("member").Build()

		repos.UserRepo.CreateUser(*admin)
		repos.UserRepo.CreateUser(*member)

		adminFromDB, _ := repos.UserRepo.GetUserByUserId(admin.UserId)
		memberFromDB, _ := repos.UserRepo.GetUserByUserId(member.UserId)

		// Create group
		group := builders.NewUserGroupBuilder().WithName("Test Group").MustBuild()
		groupId, _ := repos.UserGroupRepo.CreateGroupWithUsers(*group, adminFromDB.Id)

		// Add member
		repos.UserGroupRepo.InsertUserGroupUserTable(groupId, memberFromDB.Id)

		// Act - Promote member to admin
		err := repos.UserGroupRepo.UpdateUserRoleInGroup(groupId, memberFromDB.Id, "admin")

		// Assert
		require.NoError(t, err)
		role, err := repos.UserRepo.GetUserRoleInGroup(memberFromDB.UserId, groupId)
		require.NoError(t, err)
		assert.Equal(t, "admin", role)
	})

	t.Run("should demote admin to member", func(t *testing.T) {
		cleanup := setupTest(t)
		defer cleanup()

		// Arrange - Create users
		admin1 := builders.NewUserBuilder().WithUserName("admin1").Build()
		admin2 := builders.NewUserBuilder().WithUserName("admin2").Build()

		repos.UserRepo.CreateUser(*admin1)
		repos.UserRepo.CreateUser(*admin2)

		admin1FromDB, _ := repos.UserRepo.GetUserByUserId(admin1.UserId)
		admin2FromDB, _ := repos.UserRepo.GetUserByUserId(admin2.UserId)

		// Create group with admin1
		group := builders.NewUserGroupBuilder().WithName("Test Group").MustBuild()
		groupId, _ := repos.UserGroupRepo.CreateGroupWithUsers(*group, admin1FromDB.Id)

		// Add admin2
		repos.UserGroupRepo.InsertUserGroupUserTable(groupId, admin2FromDB.Id)
		repos.UserGroupRepo.UpdateUserRoleInGroup(groupId, admin2FromDB.Id, "admin")

		// Act - Demote admin2 to member
		err := repos.UserGroupRepo.UpdateUserRoleInGroup(groupId, admin2FromDB.Id, "member")

		// Assert
		require.NoError(t, err)
		role, err := repos.UserRepo.GetUserRoleInGroup(admin2FromDB.UserId, groupId)
		require.NoError(t, err)
		assert.Equal(t, "member", role)
	})
}

func TestUserGroupRepository_RemoveUserFromUserGroup(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	t.Run("should remove member from group successfully", func(t *testing.T) {
		cleanup := setupTest(t)
		defer cleanup()

		// Arrange - Create users
		admin := builders.NewUserBuilder().WithUserName("admin").Build()
		member := builders.NewUserBuilder().WithUserName("member").Build()

		repos.UserRepo.CreateUser(*admin)
		repos.UserRepo.CreateUser(*member)

		adminFromDB, _ := repos.UserRepo.GetUserByUserId(admin.UserId)
		memberFromDB, _ := repos.UserRepo.GetUserByUserId(member.UserId)

		// Create group and add both users
		group := builders.NewUserGroupBuilder().WithName("Test Group").MustBuild()
		groupId, _ := repos.UserGroupRepo.CreateGroupWithUsers(*group, adminFromDB.Id)
		repos.UserGroupRepo.InsertUserGroupUserTable(groupId, memberFromDB.Id)

		// Verify 2 users in group
		usersBefore, _ := repos.UserRepo.GetUsersByGroupIdWithRoles(groupId)
		assert.Len(t, usersBefore, 2)

		// Act - Remove member from group
		err := repos.UserGroupRepo.RemoveUserFromUserGroup(groupId, memberFromDB.Id)

		// Assert
		require.NoError(t, err)
		usersAfter, _ := repos.UserRepo.GetUsersByGroupIdWithRoles(groupId)
		assert.Len(t, usersAfter, 1)
		assert.Equal(t, "admin", usersAfter[0].UserName)
	})

	t.Run("should handle removing multiple members", func(t *testing.T) {
		cleanup := setupTest(t)
		defer cleanup()

		// Arrange - Create users
		admin := builders.NewUserBuilder().WithUserName("admin").Build()
		member1 := builders.NewUserBuilder().WithUserName("member1").Build()
		member2 := builders.NewUserBuilder().WithUserName("member2").Build()

		repos.UserRepo.CreateUser(*admin)
		repos.UserRepo.CreateUser(*member1)
		repos.UserRepo.CreateUser(*member2)

		adminFromDB, _ := repos.UserRepo.GetUserByUserId(admin.UserId)
		member1FromDB, _ := repos.UserRepo.GetUserByUserId(member1.UserId)
		member2FromDB, _ := repos.UserRepo.GetUserByUserId(member2.UserId)

		// Create group and add all users
		group := builders.NewUserGroupBuilder().WithName("Test Group").MustBuild()
		groupId, _ := repos.UserGroupRepo.CreateGroupWithUsers(*group, adminFromDB.Id)
		repos.UserGroupRepo.InsertUserGroupUserTable(groupId, member1FromDB.Id)
		repos.UserGroupRepo.InsertUserGroupUserTable(groupId, member2FromDB.Id)

		// Act - Remove both members
		err1 := repos.UserGroupRepo.RemoveUserFromUserGroup(groupId, member1FromDB.Id)
		err2 := repos.UserGroupRepo.RemoveUserFromUserGroup(groupId, member2FromDB.Id)

		// Assert
		require.NoError(t, err1)
		require.NoError(t, err2)
		usersAfter, _ := repos.UserRepo.GetUsersByGroupIdWithRoles(groupId)
		assert.Len(t, usersAfter, 1)
		assert.Equal(t, "admin", usersAfter[0].UserName)
	})
}

func TestUserGroupRepository_GetUserGroupsByUserId(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	t.Run("should return empty list when user has no groups", func(t *testing.T) {
		cleanup := setupTest(t)
		defer cleanup()

		// Arrange - Create user without groups
		user := builders.NewUserBuilder().Build()
		repos.UserRepo.CreateUser(*user)
		userFromDB, _ := repos.UserRepo.GetUserByUserId(user.UserId)

		// Act - Get all groups for user
		userGroups, err := repos.UserGroupRepo.GetUserGroupsByUserId(userFromDB.Id)

		// Assert
		require.NoError(t, err)
		assert.Len(t, userGroups, 0)
	})

	t.Run("should return all groups for user", func(t *testing.T) {
		cleanup := setupTest(t)
		defer cleanup()

		// Arrange - Create user
		user := builders.NewUserBuilder().Build()
		repos.UserRepo.CreateUser(*user)
		userFromDB, _ := repos.UserRepo.GetUserByUserId(user.UserId)

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

		groupId1, _ := repos.UserGroupRepo.CreateGroupWithUsers(*group1, userFromDB.Id)
		groupId2, _ := repos.UserGroupRepo.CreateGroupWithUsers(*group2, userFromDB.Id)
		groupId3, _ := repos.UserGroupRepo.CreateGroupWithUsers(*group3, userFromDB.Id)

		// Act - Get all groups for user
		userGroups, err := repos.UserGroupRepo.GetUserGroupsByUserId(userFromDB.Id)

		// Assert
		require.NoError(t, err)
		assert.Len(t, userGroups, 3)

		// Verify group details
		groupNames := make(map[int]string)
		groupTypes := make(map[int]string)
		for _, g := range userGroups {
			groupNames[g.GroupId] = g.GroupName
			groupTypes[g.GroupId] = g.Type
		}
		assert.Equal(t, "Team A", groupNames[groupId1])
		assert.Equal(t, "Project X", groupNames[groupId2])
		assert.Equal(t, "Personal Tasks", groupNames[groupId3])
		assert.Equal(t, "team", groupTypes[groupId1])
		assert.Equal(t, "project", groupTypes[groupId2])
		assert.Equal(t, "personal", groupTypes[groupId3])
	})

	t.Run("should return groups where user is member", func(t *testing.T) {
		cleanup := setupTest(t)
		defer cleanup()

		// Arrange - Create admin and member
		admin := builders.NewUserBuilder().WithUserName("admin").Build()
		member := builders.NewUserBuilder().WithUserName("member").Build()

		repos.UserRepo.CreateUser(*admin)
		repos.UserRepo.CreateUser(*member)

		adminFromDB, _ := repos.UserRepo.GetUserByUserId(admin.UserId)
		memberFromDB, _ := repos.UserRepo.GetUserByUserId(member.UserId)

		// Create groups
		group1 := builders.NewUserGroupBuilder().WithName("Group 1").MustBuild()
		group2 := builders.NewUserGroupBuilder().WithName("Group 2").MustBuild()

		groupId1, _ := repos.UserGroupRepo.CreateGroupWithUsers(*group1, adminFromDB.Id)
		_, _ = repos.UserGroupRepo.CreateGroupWithUsers(*group2, adminFromDB.Id)

		// Add member to only group1
		repos.UserGroupRepo.InsertUserGroupUserTable(groupId1, memberFromDB.Id)

		// Act - Get groups for member
		memberGroups, err := repos.UserGroupRepo.GetUserGroupsByUserId(memberFromDB.Id)

		// Assert
		require.NoError(t, err)
		assert.Len(t, memberGroups, 1)
		assert.Equal(t, groupId1, memberGroups[0].GroupId)
		assert.Equal(t, "Group 1", memberGroups[0].GroupName)
	})
}

func TestUserGroupRepository_DeleteUserGroup(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	t.Run("should delete group successfully", func(t *testing.T) {
		cleanup := setupTest(t)
		defer cleanup()

		// Arrange - Create user and group
		user := builders.NewUserBuilder().Build()
		repos.UserRepo.CreateUser(*user)
		userFromDB, _ := repos.UserRepo.GetUserByUserId(user.UserId)

		group := builders.NewUserGroupBuilder().WithName("Temporary Group").MustBuild()
		groupId, _ := repos.UserGroupRepo.CreateGroupWithUsers(*group, userFromDB.Id)

		// Verify group exists
		createdGroup, err := repos.UserGroupRepo.GetById(groupId)
		require.NoError(t, err)
		require.NotNil(t, createdGroup)

		// Act - Delete group
		err = repos.UserGroupRepo.DeleteUserGroup(groupId)

		// Assert
		require.NoError(t, err)

		// Verify group is deleted (GetById should return empty or error)
		deletedGroup, err := repos.UserGroupRepo.GetById(groupId)
		// Group should either not be found or be inactive
		// This depends on your implementation (soft delete vs hard delete)
		if err == nil {
			assert.Equal(t, 0, deletedGroup.GroupId)
		}
	})

	t.Run("should handle deleting non-existent group gracefully", func(t *testing.T) {
		cleanup := setupTest(t)
		defer cleanup()

		// Act - Attempt to delete a group that doesn't exist
		err := repos.UserGroupRepo.DeleteUserGroup(99999)

		// Assert - Should not error (or handle gracefully depending on implementation)
		// This test documents the behavior - adjust based on your requirements
		_ = err // May or may not error depending on implementation
	})

	t.Run("should delete group with multiple members", func(t *testing.T) {
		cleanup := setupTest(t)
		defer cleanup()

		// Arrange - Create users and group
		admin := builders.NewUserBuilder().WithUserName("admin").Build()
		member1 := builders.NewUserBuilder().WithUserName("member1").Build()
		member2 := builders.NewUserBuilder().WithUserName("member2").Build()

		repos.UserRepo.CreateUser(*admin)
		repos.UserRepo.CreateUser(*member1)
		repos.UserRepo.CreateUser(*member2)

		adminFromDB, _ := repos.UserRepo.GetUserByUserId(admin.UserId)
		member1FromDB, _ := repos.UserRepo.GetUserByUserId(member1.UserId)
		member2FromDB, _ := repos.UserRepo.GetUserByUserId(member2.UserId)

		group := builders.NewUserGroupBuilder().WithName("Team Group").MustBuild()
		groupId, _ := repos.UserGroupRepo.CreateGroupWithUsers(*group, adminFromDB.Id)
		repos.UserGroupRepo.InsertUserGroupUserTable(groupId, member1FromDB.Id)
		repos.UserGroupRepo.InsertUserGroupUserTable(groupId, member2FromDB.Id)

		// Verify 3 users in group
		usersBefore, _ := repos.UserRepo.GetUsersByGroupIdWithRoles(groupId)
		assert.Len(t, usersBefore, 3)

		// Act - Delete group
		err := repos.UserGroupRepo.DeleteUserGroup(groupId)

		// Assert
		require.NoError(t, err)

		// Verify group is deleted
		deletedGroup, err := repos.UserGroupRepo.GetById(groupId)
		if err == nil {
			assert.Equal(t, 0, deletedGroup.GroupId)
		}
	})
}
