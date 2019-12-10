package consts

import (
	"fmt"
	"strings"
)

type permissionTypes struct {
	Create string
	Read   string
	Update string
	Delete string
	List   string
	Assign string
	Upload string
}

type tablenames struct {
	Users       string
	Roles       string
	Permissions string
}

type roles struct {
	Admin string
	User  string
	// You can add more as you need
}

var (
	// Permissions has the types of permissions that can be assigned
	Permissions = permissionTypes{
		Create: "create:%s",
		Read:   "read:%s",
		Update: "update:%s",
		Delete: "delete:%s",
		List:   "list:%s",
		Assign: "assign:%s",
		Upload: "upload:%s",
	}
	// Tablenames the names of the tables in the server
	Tablenames = tablenames{
		Users:       "users",
		Roles:       "roles",
		Permissions: "permissions",
	}
)

// FormatPermissionTag returns a string formatted action:entity permission
func FormatPermissionTag(action string, entity string) string {
	return fmt.Sprintf(action, entity)
}

// FormatPermissionDesc returns a string with the description of the
// action:entity permission
func FormatPermissionDesc(action string, entity string) string {
	return "Allows the user to " +
		strings.ReplaceAll(FormatPermissionTag(action, entity), ":", " ")
}
