package models

// Role defines a role for the user
type Role struct {
	BaseModelSeq
	Name        string       `gorm:"not null"`
	Description string       `gorm:"size:1024"`
	ParentRoles []Role       `gorm:"many2many:role_parents;association_jointable_foreignkey:parent_role_id"`
	ChildRoles  []Role       `gorm:"many2many:role_parents;association_jointable_foreignkey:role_id"`
	Permissions []Permission `gorm:"many2many:role_permissions;association_autoupdate:false;association_autocreate:false"`
}

// Permission defines a permission scope for the user
type Permission struct {
	BaseModelSeq
	Tag         string `gorm:"not null;unique_index"`
	Description string `gorm:"size:1024"`
}
