package models

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTgUser_Permissions(t *testing.T) {
	id := int64(12345)
	name := "TestUser"
	info := "{}"
	user := NewTgUser(id, name, info)

	// Test AddPermission
	user.AddPermission(CanChat)
	assert.True(t, user.HasPermission(CanChat), "user should have CanChat permission")

	// Test RemovePermission
	user.RemovePermission(CanChat)
	assert.False(t, user.HasPermission(CanChat), "user should not have CanChat permission")

	// Test AddPermissionsFromString
	user.AddPermissionsFromString("CanDraw, CanUseSound")
	assert.True(t, user.HasPermission(CanDraw), "user should have CanDraw permission")
	assert.True(t, user.HasPermission(CanUseSound), "user should have CanUseSound permission")

	// Test PermissionsToString
	permissionsString := user.PermissionsToString()
	assert.Contains(t, permissionsString, CanDraw, "permissions string should contain CanDraw")
	assert.Contains(t, permissionsString, CanUseSound, "permissions string should contain CanUseSound")

	// Test ClearPermissions
	user.ClearPermissions()
	assert.False(t, user.HasPermission(CanDraw), "user should not have CanDraw permission after clearing")
	assert.False(t, user.HasPermission(CanUseSound), "user should not have CanUseSound permission after clearing")

	// Test HasPermission with CanEverything
	user.AddPermission(CanEverything)
	assert.True(t, user.HasPermission(CanChat), "user should have all permissions with CanEverything")
	assert.True(t, user.HasPermission(CanDraw), "user should have all permissions with CanEverything")
}

func TestNewTgUser(t *testing.T) {
	id := int64(12345)
	name := "TestUser"
	info := "TestInfo"

	user := NewTgUser(id, name, info)

	assert.Equal(t, id, user.Id, "user ID should match the input")
	assert.Equal(t, name, user.Name, "user name should match the input")
	assert.Equal(t, info, user.Info, "user info should match the input")
	assert.NotNil(t, user.Permissions, "user permissions map should be initialized")
	assert.WithinDuration(t, time.Now(), user.CreatedAt, time.Second, "user CreatedAt should be close to now")
	assert.WithinDuration(t, time.Now(), user.SeenAt, time.Second, "user SeenAt should be close to now")
}

func TestTgUser_String(t *testing.T) {
	id := int64(12345)
	name := "TestUser"
	info := "TestInfo"
	user := NewTgUser(id, name, info)
	user.AddPermission(CanChat)
	user.AddPermission(CanDraw)

	expected := fmt.Sprintf("TgUser{Id: %d, Name: %q, CreatedAt: %q, SeenAt: %q, Permissions: %q, Info: %q}",
		id, name, user.CreatedAt.Format(time.RFC3339), user.SeenAt.Format(time.RFC3339), "[CanChat,CanDraw]", info)

	assert.Equal(t, expected, user.String(), "String method output mismatch")
}
