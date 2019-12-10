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
		db = db.Begin()
		v := reflect.ValueOf(consts.Tablenames)
		tablenames := make([]interface{}, v.NumField())
		for i := 0; i < v.NumField(); i++ {
			tablenames[i] = v.Field(i).Interface()
		}
		v = reflect.ValueOf(consts.Permissions)
		permissions := make([]interface{}, v.NumField())
		for i := 0; i < v.NumField(); i++ {
			permissions[i] = v.Field(i).Interface()
		}
		for _, t := range tablenames {
			for _, p := range permissions {
				if err := db.Create(&models.Permission{
					Tag:         consts.FormatPermissionTag(p.(string), t.(string)),
					Description: consts.FormatPermissionDesc(p.(string), t.(string)),
				}).Error; err != nil {
					// db.RollbackUnlessCommitted()
					logger.Error("[Migration.Jobs.SeedRBAC] error: ", err)
					return err
				}
			}
		}
		db.Commit()
		return nil
	},
	Rollback: func(db *gorm.DB) error {
		for _, u := range users {
			if err := db.Delete(u).Error; err != nil {
				return err
			}
		}
		return nil
	},
}
