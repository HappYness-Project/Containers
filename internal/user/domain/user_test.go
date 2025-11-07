package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewUser(t *testing.T) {
	t.Run("when creating new user with valid data, Then return user with correct fields", func(t *testing.T) {
		// Given
		userId := "test-user-123"
		userName := "testuser"
		firstName := "John"
		lastName := "Doe"
		email := "john.doe@example.com"

		// When
		user := NewUser(userId, userName, firstName, lastName, email)

		// Then
		require.NotNil(t, user)
		require.NotNil(t, user.Id)
		assert.Equal(t, userId, user.UserId)
		assert.Equal(t, userName, user.UserName)
		assert.Equal(t, firstName, user.FirstName)
		assert.Equal(t, lastName, user.LastName)
		assert.Equal(t, email, user.Email)
		assert.True(t, user.IsActive)
		assert.Equal(t, 0, user.DefaultGroupId)
		assert.NotZero(t, user.CreatedAt)
		assert.NotZero(t, user.UpdatedAt)
	})

	t.Run("when creating new user with empty strings, Then return user with empty fields", func(t *testing.T) {
		// Given
		userId := ""
		userName := ""
		firstName := ""
		lastName := ""
		email := ""

		// When
		user := NewUser(userId, userName, firstName, lastName, email)

		// Then
		require.NotNil(t, user)
		assert.Equal(t, userId, user.UserId)
		assert.Equal(t, userName, user.UserName)
		assert.Equal(t, firstName, user.FirstName)
		assert.Equal(t, lastName, user.LastName)
		assert.Equal(t, email, user.Email)
		assert.True(t, user.IsActive)
		assert.Equal(t, 0, user.DefaultGroupId)
	})

	t.Run("when creating multiple users, Then each user has unique timestamps", func(t *testing.T) {
		// Given
		user1 := NewUser("user1", "username1", "First1", "Last1", "user1@example.com")
		time.Sleep(1 * time.Second) // Ensure different timestamps
		user2 := NewUser("user2", "username2", "First2", "Last2", "user2@example.com")

		// Then - compare only up to seconds precision
		user1Created := user1.CreatedAt.Truncate(time.Millisecond)
		user2Created := user2.CreatedAt.Truncate(time.Millisecond)
		user1Updated := user1.UpdatedAt.Truncate(time.Millisecond)
		user2Updated := user2.UpdatedAt.Truncate(time.Millisecond)

		assert.NotEqual(t, user1Created, user2Created)
		assert.NotEqual(t, user1Updated, user2Updated)
		assert.True(t, user2Created.After(user1Created))
		assert.True(t, user2Updated.After(user1Updated))
	})

	t.Run("when creating user, Then CreatedAt and UpdatedAt are initially equal", func(t *testing.T) {
		user := NewUser("test-user", "testuser", "Test", "User", "test@example.com")

		createdAt := user.CreatedAt.Truncate(time.Millisecond)
		updatedAt := user.UpdatedAt.Truncate(time.Millisecond)
		assert.Equal(t, createdAt, updatedAt)
	})
}

func TestUserStruct(t *testing.T) {
	t.Run("when creating user struct directly, Then all fields can be set", func(t *testing.T) {
		// Given
		now := time.Now()
		user := User{
			Id:             1,
			UserId:         "direct-user",
			UserName:       "directuser",
			FirstName:      "Direct",
			LastName:       "User",
			Email:          "direct@example.com",
			IsActive:       false,
			CreatedAt:      now,
			UpdatedAt:      now,
			DefaultGroupId: 5,
		}

		// Then
		assert.Equal(t, 1, user.Id)
		assert.Equal(t, "direct-user", user.UserId)
		assert.Equal(t, "directuser", user.UserName)
		assert.Equal(t, "Direct", user.FirstName)
		assert.Equal(t, "User", user.LastName)
		assert.Equal(t, "direct@example.com", user.Email)
		assert.False(t, user.IsActive)
		assert.Equal(t, now, user.CreatedAt)
		assert.Equal(t, now, user.UpdatedAt)
		assert.Equal(t, 5, user.DefaultGroupId)
	})

	t.Run("when modifying user fields, Then changes are reflected", func(t *testing.T) {
		// Given
		user := NewUser("original", "originaluser", "Original", "User", "original@example.com")
		originalCreatedAt := user.CreatedAt

		// When
		user.FirstName = "Modified"
		user.LastName = "Name"
		user.Email = "modified@example.com"
		user.IsActive = false
		user.DefaultGroupId = 10
		user.UpdatedAt = time.Now()

		// Then
		assert.Equal(t, "Modified", user.FirstName)
		assert.Equal(t, "Name", user.LastName)
		assert.Equal(t, "modified@example.com", user.Email)
		assert.False(t, user.IsActive)
		assert.Equal(t, 10, user.DefaultGroupId)
		assert.Equal(t, originalCreatedAt, user.CreatedAt)
		assert.True(t, user.UpdatedAt.After(originalCreatedAt))
	})
}

func TestUpdateDefaultGroupId(t *testing.T) {
	t.Run("when updating default group ID with valid positive ID, Then update succeeds", func(t *testing.T) {
		user := NewUser("test-user", "testuser", "Test", "User", "test@example.com")
		originalUpdatedAt := user.UpdatedAt
		time.Sleep(1 * time.Millisecond) // Ensure time difference

		err := user.UpdateDefaultGroupId(5)

		assert.NoError(t, err)
		assert.Equal(t, 5, user.DefaultGroupId)
		assert.True(t, user.UpdatedAt.After(originalUpdatedAt))
	})

	t.Run("when updating default group ID with zero, Then update succeeds", func(t *testing.T) {
		// Given
		user := NewUser("test-user", "testuser", "Test", "User", "test@example.com")
		user.DefaultGroupId = 10 // Set initial group ID
		originalUpdatedAt := user.UpdatedAt
		time.Sleep(1 * time.Millisecond) // Ensure time difference

		err := user.UpdateDefaultGroupId(0)

		assert.NoError(t, err)
		assert.Equal(t, 0, user.DefaultGroupId)
		assert.True(t, user.UpdatedAt.After(originalUpdatedAt))
	})

	t.Run("when updating default group ID with negative ID, Then return error", func(t *testing.T) {
		user := NewUser("test-user", "testuser", "Test", "User", "test@example.com")
		originalGroupId := user.DefaultGroupId
		originalUpdatedAt := user.UpdatedAt

		err := user.UpdateDefaultGroupId(-1)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "group ID cannot be negative")
		assert.Equal(t, originalGroupId, user.DefaultGroupId)
		assert.Equal(t, originalUpdatedAt, user.UpdatedAt) // Should not be updated
	})

	t.Run("when updating default group ID with same ID, Then return error", func(t *testing.T) {
		// Given
		user := NewUser("test-user", "testuser", "Test", "User", "test@example.com")
		user.DefaultGroupId = 5 // Set initial group ID
		originalUpdatedAt := user.UpdatedAt

		// When
		err := user.UpdateDefaultGroupId(5)

		// Then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "group ID is already set to the specified value")
		assert.Equal(t, 5, user.DefaultGroupId)
		assert.Equal(t, originalUpdatedAt, user.UpdatedAt) // Should not be updated
	})

	t.Run("when updating default group ID multiple times, Then each update succeeds", func(t *testing.T) {
		// Given
		user := NewUser("test-user", "testuser", "Test", "User", "test@example.com")

		// When & Then - First update
		err1 := user.UpdateDefaultGroupId(5)
		assert.NoError(t, err1)
		assert.Equal(t, 5, user.DefaultGroupId)

		// Second update
		err2 := user.UpdateDefaultGroupId(10)
		assert.NoError(t, err2)
		assert.Equal(t, 10, user.DefaultGroupId)

		// Third update
		err3 := user.UpdateDefaultGroupId(3)
		assert.NoError(t, err3)
		assert.Equal(t, 3, user.DefaultGroupId)
	})

	t.Run("when updating default group ID, Then UpdatedAt timestamp is modified", func(t *testing.T) {
		// Given
		user := NewUser("test-user", "testuser", "Test", "User", "test@example.com")
		originalUpdatedAt := user.UpdatedAt
		time.Sleep(1 * time.Millisecond) // Ensure time difference

		// When
		err := user.UpdateDefaultGroupId(7)

		// Then
		assert.NoError(t, err)
		assert.True(t, user.UpdatedAt.After(originalUpdatedAt))
		assert.Equal(t, 7, user.DefaultGroupId)
	})

	t.Run("when updating default group ID on user with existing group, Then update succeeds", func(t *testing.T) {
		// Given
		user := &User{
			Id:             1,
			UserId:         "existing-user",
			UserName:       "existinguser",
			FirstName:      "Existing",
			LastName:       "User",
			Email:          "existing@example.com",
			IsActive:       true,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
			DefaultGroupId: 3, // Existing group ID
		}
		originalUpdatedAt := user.UpdatedAt
		time.Sleep(1 * time.Millisecond) // Ensure time difference

		// When
		err := user.UpdateDefaultGroupId(8)

		// Then
		assert.NoError(t, err)
		assert.Equal(t, 8, user.DefaultGroupId)
		assert.True(t, user.UpdatedAt.After(originalUpdatedAt))
	})
}

func TestClearDefaultGroup(t *testing.T) {
	t.Run("when clearing default group, Then group ID is set to zero", func(t *testing.T) {
		// Given
		user := NewUser("test-user", "testuser", "Test", "User", "test@example.com")
		user.DefaultGroupId = 5 // Set initial group ID
		originalUpdatedAt := user.UpdatedAt
		time.Sleep(1 * time.Millisecond) // Ensure time difference

		// When
		user.ClearDefaultGroup()

		// Then
		assert.Equal(t, 0, user.DefaultGroupId)
		assert.True(t, user.UpdatedAt.After(originalUpdatedAt))
	})

	t.Run("when clearing default group that is already zero, Then no change occurs", func(t *testing.T) {
		// Given
		user := NewUser("test-user", "testuser", "Test", "User", "test@example.com")
		user.DefaultGroupId = 0 // Already zero
		originalUpdatedAt := user.UpdatedAt
		time.Sleep(1 * time.Millisecond) // Ensure time difference

		// When
		user.ClearDefaultGroup()

		// Then
		assert.Equal(t, 0, user.DefaultGroupId)
		assert.True(t, user.UpdatedAt.After(originalUpdatedAt))
	})
}

func TestUpdateUser(t *testing.T) {
	t.Run("when updating user fields, Then fields and UpdatedAt are updated", func(t *testing.T) {
		user := NewUser("test-user", "testuser", "OldFirst", "OldLast", "old@example.com")
		originalUpdatedAt := user.UpdatedAt
		time.Sleep(1 * time.Millisecond) // Ensure UpdatedAt changes

		user.UpdateUser("NewFirst", "NewLast", "new@example.com")

		assert.Equal(t, "NewFirst", user.FirstName)
		assert.Equal(t, "NewLast", user.LastName)
		assert.Equal(t, "new@example.com", user.Email)
		assert.True(t, user.UpdatedAt.After(originalUpdatedAt))
	})
}
