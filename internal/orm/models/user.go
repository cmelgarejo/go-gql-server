package models

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

// ## Entity definitions

// User defines a user for the app
type User struct {
	BaseModelSoftDelete        // We don't to actually delete the users, audit
	Email               string `gorm:"not null"`
	Password            string
	Name                *string `gorm:"null"`
	NickName            *string
	FirstName           *string
	LastName            *string
	Location            *string
	AvatarURL           *string `gorm:"size:1024"`
	Description         *string `gorm:"size:1024"`
	Profiles            []UserProfile
	Roles               []Role       `gorm:"many2many:user_roles;association_autoupdate:false;association_autocreate:false"`
	Permissions         []Permission `gorm:"many2many:user_permissions;association_autoupdate:false;association_autocreate:false"`
}

// UserProfile saves all the related OAuth Profiles
type UserProfile struct {
	BaseModelSeq
	Email          string    `gorm:"unique_index:idx_email_provider_external_user_id"`
	User           User      `gorm:"association_autoupdate:false;association_autocreate:false"`
	UserID         uuid.UUID `gorm:"not null;index"`
	Provider       string    `gorm:"not null;index;unique_index:idx_email_provider_external_user_id;default:'DB'"` // DB means database or no ExternalUserID
	ExternalUserID string    `gorm:"not null;index;unique_index:idx_email_provider_external_user_id"`              // User ID
	Name           string
	NickName       string
	FirstName      string
	LastName       string
	Location       string  `gorm:"size:1024"`
	AvatarURL      string  `gorm:"size:1024"`
	Description    *string `gorm:"size:1024"`
}

// UserAPIKey generated api keys for the users
type UserAPIKey struct {
	BaseModelSeq
	User        User      `gorm:"association_autoupdate:false;association_autocreate:false"`
	UserID      uuid.UUID `gorm:"not null;index"`
	APIKey      string    `gorm:"size:128;unique_index"`
	Name        string
	Permissions []Permission `gorm:"many2many:user_permissions;association_autoupdate:false;association_autocreate:false"`
}

// UserRole relation between an user and its roles
type UserRole struct {
	UserID uuid.UUID `gorm:"index"`
	RoleID int       `gorm:"index"`
}

// UserPermission relation between an user and its permissions
type UserPermission struct {
	UserID       uuid.UUID `gorm:"index"`
	PermissionID int       `gorm:"index"`
}

// ## Hooks

// BeforeSave hook for User
func (u *User) BeforeSave(scope *gorm.Scope) error {
	if u.Password != "" {
		if pw, err := bcrypt.GenerateFromPassword([]byte(u.Password), 11); err == nil {
			scope.SetColumn("Password", pw)
		}
	}
	return nil
}

// AfterSave hook (assigning roles, fill all permissions for example)

// BeforeSave hook for UserAPIKey
func (k *UserAPIKey) BeforeSave(scope *gorm.Scope) error {
	db := scope.NewDB()
	if k.Name == "" {
		u := &User{}
		if err := db.Where("id = ?", k.UserID).First(u).Error; err != nil {
			return err
		}
		// k.Name =
	}
	if hash, err := bcrypt.GenerateFromPassword([]byte(k.UserID.String()), 0); err == nil {
		hasher := sha1.New()
		hasher.Write(hash)
		scope.SetColumn("APIKey", hex.EncodeToString(hasher.Sum(nil)))
	}
	return nil
}

// ## Helper functions

// HasRole verifies if user posseses a role
func (u *User) HasRole(roleID int) (bool, error) {
	for _, r := range u.Roles {
		if r.ID == roleID {
			return true, nil
		}
	}
	return false, fmt.Errorf("The user has no [%d] roleID", roleID)
}

// HasPermission verifies if user has a specific permission
func (u *User) HasPermission(permission string, entity string) (bool, error) {
	tag := fmt.Sprintf(permission, entity)
	for _, r := range u.Permissions {
		if r.Tag == tag {
			return true, nil
		}
	}
	return false, fmt.Errorf("user has no permission: [%s]", tag)
}

// HasPermissionBool verifies if user has a specific permission - returns t/f
func (u *User) HasPermissionBool(permission string, entity string) bool {
	p, _ := u.HasPermission(permission, entity)
	return p
}

// HasPermissionTag verifies if user has a specific permission tag
func (u *User) HasPermissionTag(tag string) (bool, error) {
	for _, r := range u.Permissions {
		if r.Tag == tag {
			return true, nil
		}
	}
	return false, fmt.Errorf("The user has no [%s] permission", tag)
}

// GetName returns the displayName if not nil, or the first + last name
func (u *User) GetName() string {
	userName := ""
	if u.FirstName != nil {
		userName += *u.FirstName
	}
	if u.LastName != nil {
		userName += " " + *u.LastName
	}
	if u.LastName != nil {
		userName = *u.LastName
	}
	return userName
}
