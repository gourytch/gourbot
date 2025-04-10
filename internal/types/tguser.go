package types

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// Permission constants for TgUser.
const (
	CanEverything       = "CanEverything"
	CanChat             = "CanChat"
	CanDraw             = "CanDraw"
	CanUseSound         = "CanUseSound"
	CanUseRoles         = "CanUseRoles"
	CanManageRoles      = "CanManageRoles"
	CanGetStatistics    = "CanGetStatistics"
	CanGetAllStatistics = "CanGetAllStatistics"
)

// TgUser represents a Telegram user.
type TgUser struct {
	Id          int64           // Unique identifier from Telegram API, stored as INTEGER in the database
	Name        string          // Name or nickname of the user, stored as TEXT in the database
	CreatedAt   time.Time       // When the user was first seen, stored as INTEGER (Unix time) in the database
	SeenAt      time.Time       // When the user was last seen, stored as INTEGER (Unix time) in the database
	Permissions map[string]bool // Set of permissions for the user, stored as TEXT (comma-separated) in the database
	Info        []byte          // Additional information about the user record, stored as TEXT in the database
}

// Constructor for TgUser that initializes Permissions as an empty map.
func NewTgUser(id int64, name string, info []byte) *TgUser {
	return &TgUser{
		Id:          id,
		Name:        name,
		CreatedAt:   time.Now(),
		SeenAt:      time.Now(),
		Permissions: make(map[string]bool),
		Info:        info,
	}
}

// ClearPermissions removes all permissions from the user.
func (u *TgUser) ClearPermissions() {
	u.Permissions = make(map[string]bool)
}

// AddPermissionsFromString parses a comma-separated string and adds permissions to the user.
func (u *TgUser) AddPermissionsFromString(permissions string) {
	if u.Permissions == nil {
		u.Permissions = make(map[string]bool)
	}
	for _, perm := range strings.Split(permissions, ",") {
		perm = strings.TrimSpace(perm)
		if perm == "" {
			continue // Skip empty permissions
		}
		u.Permissions[perm] = true
	}
}

// PermissionsToString converts the user's permissions to a comma-separated string.
func (u *TgUser) PermissionsToString() string {
	var perms []string
	for perm := range u.Permissions {
		perms = append(perms, perm)
	}
	sort.Strings(perms) // Ensure permissions are sorted alphabetically
	return strings.Join(perms, ",")
}

// AddPermission adds a single permission to the user.
func (u *TgUser) AddPermission(permission string) {
	if u.Permissions == nil {
		u.Permissions = make(map[string]bool)
	}
	u.Permissions[permission] = true
}

// RemovePermission removes a single permission from the user.
func (u *TgUser) RemovePermission(permission string) {
	if u.Permissions != nil {
		delete(u.Permissions, permission)
	}
}

// HasPermission checks if the user has a specific permission or the CanEverything permission.
func (u *TgUser) HasPermission(permission string) bool {
	if u.Permissions == nil {
		return false
	}
	return u.Permissions[CanEverything] || u.Permissions[permission]
}

// String formats the TgUser fields into a human-readable string.
// - Id: Displayed as is.
// - Name: Quoted string.
// - CreatedAt and SeenAt: Formatted as RFC3339 and quoted.
// - Permissions: Converted to a comma-separated string and enclosed in square brackets.
// - Info: Quoted string.
func (u *TgUser) String() string {
	permissions := u.PermissionsToString()
	if permissions == "" {
		permissions = "[]"
	} else {
		permissions = fmt.Sprintf("[%s]", permissions)
	}
	return fmt.Sprintf("TgUser{Id: %d, Name: %q, CreatedAt: %q, SeenAt: %q, Permissions: %q, Info: %q}",
		u.Id, u.Name, u.CreatedAt.Format(time.RFC3339), u.SeenAt.Format(time.RFC3339), permissions, u.Info)
}
