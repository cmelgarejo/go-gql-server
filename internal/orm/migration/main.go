package migration

import (
	"fmt"

	"github.com/cmelgarejo/go-gql-server/internal/logger"
	"github.com/cmelgarejo/go-gql-server/internal/orm/migration/jobs"
	"github.com/cmelgarejo/go-gql-server/internal/orm/models"
	"github.com/cmelgarejo/go-gql-server/pkg/utils/consts"
	"github.com/jinzhu/gorm"
	"gopkg.in/gormigrate.v1"
)

func updateMigration(db *gorm.DB) (err error) {
	db.AutoMigrate(
		&models.Role{},
		&models.Permission{},
		&models.UserProfile{},
		&models.UserAPIKey{},
		&models.User{},
	)
	return addIndexes(db)
}

func addIndexes(db *gorm.DB) (err error) {
	// Entity names
	//db.NewScope(&models.User{}).GetModelStruct().TableName(db)

	// FKs
	if err := db.Model(&models.UserProfile{}).
		AddForeignKey("user_id", consts.GetTableName(consts.EntityNames.Users)+"(id)", "RESTRICT", "RESTRICT").Error; err != nil {
		return err
	}
	if err := db.Model(&models.UserAPIKey{}).
		AddForeignKey("user_id", consts.GetTableName(consts.EntityNames.Users)+"(id)", "RESTRICT", "RESTRICT").Error; err != nil {
		return err
	}
	if err := db.Model(&models.UserRole{}).
		AddForeignKey("user_id", consts.GetTableName(consts.EntityNames.Users)+"(id)", "CASCADE", "CASCADE").Error; err != nil {
		return err
	}
	if err := db.Model(&models.UserRole{}).
		AddForeignKey("role_id", consts.GetTableName(consts.EntityNames.Roles)+"(id)", "CASCADE", "CASCADE").Error; err != nil {
		return err
	}
	if err := db.Model(&models.UserPermission{}).
		AddForeignKey("user_id", consts.GetTableName(consts.EntityNames.Users)+"(id)", "CASCADE", "CASCADE").Error; err != nil {
		return err
	}
	if err := db.Model(&models.UserPermission{}).
		AddForeignKey("permission_id", consts.GetTableName(consts.EntityNames.Permissions)+"(id)", "CASCADE", "CASCADE").Error; err != nil {
		return err
	}
	// Indexes
	// None needed so far
	return nil
}

// ServiceAutoMigration migrates all the tables and modifications to the connected source
func ServiceAutoMigration(db *gorm.DB) error {
	// Initialize the migration empty so InitSchema runs always first on creation
	m := gormigrate.New(db, gormigrate.DefaultOptions, nil)
	m.InitSchema(func(db *gorm.DB) error {
		logger.Info("[Migration.InitSchema] Initializing database schema")
		switch db.Dialect().GetName() {
		case consts.Dialects.PostgresSQL:
			db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)
		}
		if err := updateMigration(db); err != nil {
			return fmt.Errorf("[Migration.InitSchema]: %v", err)
		}
		return nil
	})
	m.Migrate()

	if err := updateMigration(db); err != nil {
		return err
	}
	// Keep a list of migrations here
	m = gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		jobs.SeedRBAC,
		jobs.SeedUsers,
	})
	return m.Migrate()
}
