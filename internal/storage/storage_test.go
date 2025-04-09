package storage

import (
	"testing"
	"time"

	"gourbot/internal/models"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func TestStorage_OpenClose(t *testing.T) {
	// Use in-memory SQLite database
	storage := NewStorage(":memory:")
	err := storage.Open()
	assert.NoError(t, err, "failed to open storage")

	err = storage.Close()
	assert.NoError(t, err, "failed to close storage")
}

func TestStorage_CreateTables(t *testing.T) {
	storage := NewStorage(":memory:")
	err := storage.Open()
	assert.NoError(t, err, "failed to open storage")
	defer storage.Close()

	// Ensure tables are created without errors
	err = storage.createTables()
	assert.NoError(t, err, "failed to create tables")
}

func TestStorage_AddTgRecord(t *testing.T) {
	storage := NewStorage(":memory:")
	err := storage.Open()
	assert.NoError(t, err, "failed to open storage")
	defer storage.Close()

	// Add a record to tgdump
	err = storage.AddTgRecord(true, "{\"key\":\"value\"}")
	assert.NoError(t, err, "failed to add record to tgdump")

	// Verify the record exists
	row := storage.db.QueryRow("SELECT out, data FROM tgdump WHERE uid = 1")
	var out bool
	var data string
	err = row.Scan(&out, &data)
	assert.NoError(t, err, "failed to query record from tgdump")
	assert.Equal(t, true, out, "unexpected value for 'out'")
	assert.Equal(t, "{\"key\":\"value\"}", data, "unexpected value for 'data'")
}

func TestStorage_AddTgUser(t *testing.T) {
	storage := NewStorage(":memory:")
	err := storage.Open()
	assert.NoError(t, err, "failed to open storage")
	defer storage.Close()

	id := int64(12345)
	name := "TestUser"
	info := "{}"
	user := models.NewTgUser(id, name, info)

	err = storage.AddTgUser(user)
	assert.NoError(t, err, "failed to add user")
}

func TestStorage_TgUserExists(t *testing.T) {
	storage := NewStorage(":memory:")
	err := storage.Open()
	assert.NoError(t, err, "failed to open storage")
	defer storage.Close()

	id := int64(12345)
	name := "TestUser"
	info := "{}"
	user := models.NewTgUser(id, name, info)

	err = storage.AddTgUser(user)
	assert.NoError(t, err, "failed to add user")

	exists, err := storage.TgUserExists(user.Id)
	assert.NoError(t, err, "failed to check if user exists")
	assert.True(t, exists, "user should exist")
}

func TestStorage_GetTgUser(t *testing.T) {
	storage := NewStorage(":memory:")
	err := storage.Open()
	assert.NoError(t, err, "failed to open storage")
	defer storage.Close()

	id := int64(12345)
	name := "TestUser"
	info := "{}"
	user := models.NewTgUser(id, name, info)

	err = storage.AddTgUser(user)
	assert.NoError(t, err, "failed to add user")

	retrievedUser, err := storage.GetTgUser(user.Id)
	assert.NoError(t, err, "failed to get user")
	assert.Equal(t, user.Id, retrievedUser.Id, "user ID mismatch")
	assert.Equal(t, user.Name, retrievedUser.Name, "user name mismatch")
	assert.Equal(t, user.Permissions, retrievedUser.Permissions, "user permissions mismatch")
	assert.Equal(t, user.Info, retrievedUser.Info, "user info mismatch")
}

func TestStorage_GetAllTgUsers(t *testing.T) {
	storage := NewStorage(":memory:")
	err := storage.Open()
	assert.NoError(t, err, "failed to open storage")
	defer storage.Close()

	id1 := int64(12345)
	name1 := "TestUser1"
	info1 := "{}"
	user1 := models.NewTgUser(id1, name1, info1)

	id2 := int64(67890)
	name2 := "TestUser2"
	info2 := "{}"
	user2 := models.NewTgUser(id2, name2, info2)

	err = storage.AddTgUser(user1)
	assert.NoError(t, err, "failed to add user1")
	err = storage.AddTgUser(user2)
	assert.NoError(t, err, "failed to add user2")

	users, err := storage.GetAllTgUsers()
	assert.NoError(t, err, "failed to get all users")
	assert.Len(t, users, 2, "unexpected number of users")
}

func TestStorage_UpdateTgUser(t *testing.T) {
	storage := NewStorage(":memory:")
	err := storage.Open()
	assert.NoError(t, err, "failed to open storage")
	defer storage.Close()

	id := int64(12345)
	name := "TestUser"
	info := "{}"
	user := models.NewTgUser(id, name, info)

	err = storage.AddTgUser(user)
	assert.NoError(t, err, "failed to add user")

	user.Name = "UpdatedUser"
	user.SeenAt = time.Now()
	user.Permissions = map[string]bool{"read": true, "write": true}
	user.Info = "{\"updated\":true}"

	err = storage.UpdateTgUser(user)
	assert.NoError(t, err, "failed to update user")

	updatedUser, err := storage.GetTgUser(user.Id)
	assert.NoError(t, err, "failed to get updated user")
	assert.Equal(t, user.Name, updatedUser.Name, "user name mismatch after update")
	assert.Equal(t, user.Permissions, updatedUser.Permissions, "user permissions mismatch after update")
	assert.Equal(t, user.Info, updatedUser.Info, "user info mismatch after update")
}
