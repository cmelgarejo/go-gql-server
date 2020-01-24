package jobs

import (
	"reflect"

	"github.com/cmelgarejo/go-gql-server/internal/logger"
	"github.com/cmelgarejo/go-gql-server/internal/orm/models"
	"github.com/cmelgarejo/go-gql-server/pkg/utils/consts"

	"github.com/jinzhu/gorm"
	"gopkg.in/gormigrate.v1"
)

// SeedRBAC inserts the first users
var SeedRBAC *gormigrate.Migration = &gormigrate.Migration{
	ID: "SEED_RBAC",
	Migrate: func(db *gorm.DB) error {
		tx := db.Begin()
		defer tx.RollbackUnlessCommitted()
		v := reflect.ValueOf(consts.EntityNames)
		tablenames := make([]interface{}, v.NumField())
		for i := 0; i < v.NumField(); i++ {
			tablenames[i] = consts.GetTableName(v.Field(i).Interface().(string))
		}
		v = reflect.ValueOf(consts.Permissions)
		permissions := make([]interface{}, v.NumField())
		for i := 0; i < v.NumField(); i++ {
			permissions[i] = v.Field(i).Interface()
		}
		padmin := []models.Permission{}
		for _, t := range tablenames {
			for _, p := range permissions {
				permission := models.Permission{
					Tag:         consts.FormatPermissionTag(p.(string), t.(string)),
					Description: consts.FormatPermissionDesc(p.(string), t.(string)),
				}
				if err := tx.Create(&permission).First(&permission).Error; err != nil {
					logger.Error("[Migration.Jobs.SeedRBAC.permissions] error: ", err)
					return err
				}
				padmin = append(padmin, permission)
			}
		}
		for _, r := range consts.Roles {
			role := &models.Role{
				Name:        r.Name,
				Description: r.Description,
			}
			if err := tx.Create(role).First(&role).Error; err != nil {
				logger.Error("[Migration.Jobs.SeedRBAC.roles] error: ", err)
				return err
			}
			switch r.Name {
			case "admin":
				for _, p := range padmin {
					tx.Model(role).Association(consts.EntityNames.Permissions).Append(p)
				}
			case "user":
				// Permissions for user role
				// for _, p := range puser {
				// 	tx.Model(role).Association(consts.EntityNames.Permissions).Append(p)
				// }
			}
		}
		tx.Commit()
		return nil
	},
	Rollback: func(db *gorm.DB) error {
		tx := db.Begin()
		defer tx.RollbackUnlessCommitted()
		for _, u := range users {
			if err := tx.Delete(u).Error; err != nil {
				return err
			}
		}
		tx.Commit()
		return nil
	},
}
