package models

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/cmelgarejo/go-gql-server/pkg/utils/consts"

	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

// ## Entity definitions

// User defines a user for the app
type User struct {
	BaseModelSoftDelete        // We don't to actually delete the users, audit
	Email               string `gorm:"not null;index"`
	Password            string
	Name                *string `gorm:"null"`
	NickName            *string
	FirstName           *string
	LastName            *string
	Location            *string
	AvatarURL           *string       `gorm:"size:1024"`
	Description         *string       `gorm:"size:1024"`
	UserProfiles        []UserProfile `gorm:"association_autocreate:false;association_autoupdate:false"`
	Roles               []Role        `gorm:"many2many:user_roles;association_autocreate:false;association_autoupdate:false"`
	Permissions         []Permission  `gorm:"many2many:user_permissions;association_autocreate:false;association_autoupdate:false"`
	CreatedBy           *User         `gorm:"association_autoupdate:false;association_autocreate:false"`
	UpdatedBy           *User         `gorm:"association_autoupdate:false;association_autocreate:false"`
}

// UserProfile saves all the related OAuth Profiles
type UserProfile struct {
	BaseModelSeq
	Email          string    `gorm:"unique_index:idx_email_provider_external_user_id"`
	UserID         uuid.UUID `gorm:"not null;index"`
	User           User      `gorm:"association_autocreate:false;association_autoupdate:false"`
	Provider       string    `gorm:"not null;index;unique_index:idx_email_provider_external_user_id;default:'DB'"` // DB means database or no ExternalUserID
	ExternalUserID string    `gorm:"not null;index;unique_index:idx_email_provider_external_user_id"`              // User ID
	Name           string
	NickName       string
	FirstName      string
	LastName       string
	Location       string `gorm:"size:512"`
	AvatarURL      string `gorm:"size:1024"`
	Description    string `gorm:"size:1024"`
	CreatedBy      *User  `gorm:"association_autoupdate:false;association_autocreate:false"`
	UpdatedBy      *User  `gorm:"association_autoupdate:false;association_autocreate:false"`
}

// UserAPIKey generated api keys for the users
type UserAPIKey struct {
	BaseModelSeq
	Name        string
	User        User         `gorm:"association_autocreate:false;association_autoupdate:false"`
	UserID      uuid.UUID    `gorm:"not null;index"`
	APIKey      string       `gorm:"size:128;unique_index"`
	Permissions []Permission `gorm:"many2many:user_api_key_permissions;association_autocreate:false;association_autoupdate:false"`
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

// AfterSave hook for User
func (u *User) AfterSave(scope *gorm.Scope) error {
	db := scope.DB().
		Preload(consts.EntityNames.Roles).Preload(consts.EntityNames.Permissions).
		First(u)
	// Deal with role changes
	db.Model(u).Association(consts.EntityNames.Permissions).Clear()
	for _, r := range u.Roles {
		if err := scope.DB().Model(r).Preload(consts.EntityNames.Permissions).First(&r).Error; err != nil {
			return err
		}
		if len(r.Permissions) > 0 {
			if err := db.Model(u).Association(consts.EntityNames.Permissions).Append(r.Permissions).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

// AfterSave hook (assigning roles, fill all permissions for example)
func (ur *UserRole) AfterSave(scope *gorm.Scope) error {
	db := scope.DB()
	role := Role{}
	user := User{}
	user.ID = ur.UserID
	role.ID = ur.RoleID
	db.Model(role).Preload(consts.EntityNames.Permissions).First(&role)
	if err := db.Model(user).First(&user).Association(consts.EntityNames.Permissions).
		Replace(role.Permissions).Error; err != nil {
		return err
	}
	return nil
}

// BeforeSave hook for UserAPIKey
func (k *UserAPIKey) BeforeSave(scope *gorm.Scope) error {
	db := scope.DB()
	if k.Name == "" {
		u := &User{}
		if err := db.Where("id = ?", k.UserID).First(u).Error; err != nil {
			return err
		}
	}
	if hash, err := bcrypt.GenerateFromPassword([]byte(k.UserID.String()), 0); err == nil {
		hasher := sha1.New()
		hasher.Write(hash)
		scope.SetColumn("APIKey", hex.EncodeToString(hasher.Sum(nil)))
	}
	return nil
}

// ## Helper functions

// HasRole verifies if user possesses a role
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
	tag := fmt.Sprintf(permission, consts.GetTableName(entity))
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

// GetDisplayName returns the displayName if not nil, or the first + last name
func (u *User) GetDisplayName() string {
	displayName := ""
	if u.FirstName != nil {
		displayName += *u.FirstName
	}
	if u.LastName != nil {
		displayName += " " + *u.LastName
	}
	if u.LastName != nil {
		displayName = *u.LastName
	}
	return strings.TrimSpace(displayName)
}
