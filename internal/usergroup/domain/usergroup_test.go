package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewUserGroup(t *testing.T) {
	t.Run("when creating new user group with valid data, Then return user group with correct fields", func(t *testing.T) {
		// Given
		name := "Admin Group"
		desc := "Administrators group"
		groupType := "admin"

		// When
		userGroup, err := NewUserGroup(name, desc, groupType)

		// Then
		require.NoError(t, err)
		require.NotNil(t, userGroup)
		assert.Equal(t, name, userGroup.GroupName)
		assert.Equal(t, desc, userGroup.GroupDesc)
		assert.Equal(t, groupType, userGroup.Type)
		assert.True(t, userGroup.IsActive)
		assert.Equal(t, "", userGroup.Thumbnail)
		assert.Equal(t, 0, userGroup.GroupId) // Should be zero for new group
	})

	t.Run("when creating new user group with empty name, Then return error", func(t *testing.T) {
		// Given
		name := ""
		desc := "Some description"
		groupType := "regular"

		// When
		userGroup, err := NewUserGroup(name, desc, groupType)

		// Then
		assert.Error(t, err)
		assert.Nil(t, userGroup)
		assert.Contains(t, err.Error(), "GroupName field cannot be empty")
	})

	t.Run("when creating new user group with empty group type, Then return error", func(t *testing.T) {
		// Given
		name := "Test Group"
		desc := "Some description"
		groupType := ""

		// When
		userGroup, err := NewUserGroup(name, desc, groupType)

		// Then
		assert.Error(t, err)
		assert.Nil(t, userGroup)
		assert.Contains(t, err.Error(), "GroupType field cannot be empty")
	})

	t.Run("when creating new user group with empty description, Then return user group with empty description", func(t *testing.T) {
		// Given
		name := "Test Group"
		desc := ""
		groupType := "regular"

		// When
		userGroup, err := NewUserGroup(name, desc, groupType)

		// Then
		require.NoError(t, err)
		require.NotNil(t, userGroup)
		assert.Equal(t, name, userGroup.GroupName)
		assert.Equal(t, desc, userGroup.GroupDesc)
		assert.Equal(t, groupType, userGroup.Type)
		assert.True(t, userGroup.IsActive)
		assert.Equal(t, "", userGroup.Thumbnail)
	})

	t.Run("when creating new user group with whitespace name, Then return user group", func(t *testing.T) {
		// Given
		name := "   Test Group   "
		desc := "Some description"
		groupType := "regular"

		// When
		userGroup, err := NewUserGroup(name, desc, groupType)

		// Then
		require.NoError(t, err)
		require.NotNil(t, userGroup)
		assert.Equal(t, name, userGroup.GroupName)
		assert.Equal(t, desc, userGroup.GroupDesc)
		assert.Equal(t, groupType, userGroup.Type)
	})

	t.Run("when creating new user group with whitespace group type, Then return user group", func(t *testing.T) {
		// Given
		name := "Test Group"
		desc := "Some description"
		groupType := "   regular   "

		// When
		userGroup, err := NewUserGroup(name, desc, groupType)

		// Then
		require.NoError(t, err)
		require.NotNil(t, userGroup)
		assert.Equal(t, name, userGroup.GroupName)
		assert.Equal(t, desc, userGroup.GroupDesc)
		assert.Equal(t, groupType, userGroup.Type)
	})
}
