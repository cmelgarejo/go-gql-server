// Package orm provides `GORM` helpers for the creation, migration and access
// on the project's database
package orm

import (
	log "github.com/cmelgarejo/go-gql-server/internal/logger"

	"github.com/cmelgarejo/go-gql-server/internal/orm/migration"

	"github.com/cmelgarejo/go-gql-server/pkg/utils"
	//Imports the database dialect of choice
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/jinzhu/gorm"
)

var autoMigrate, logMode, seedDB bool
var dsn, dialect string

// ORM struct to holds the gorm pointer to db
type ORM struct {
	DB *gorm.DB
}

func init() {
	dialect = utils.MustGet("GORM_DIALECT")
	dsn = utils.MustGet("GORM_CONNECTION_DSN")
	seedDB = utils.MustGetBool("GORM_SEED_DB")
	logMode = utils.MustGetBool("GORM_LOGMODE")
	autoMigrate = utils.MustGetBool("GORM_AUTOMIGRATE")
}

// Factory creates a db connection with the selected dialect and connection
// string
func Factory() (*ORM, error) {
	db, err := gorm.Open(dialect, dsn)
	if err != nil {
		log.Panic("[ORM] err: ", err)
	}
	orm := &ORM{
		DB: db,
	}
	// Log every SQL command on dev, @prod: this should be disabled?
	db.LogMode(logMode)
	// Automigrate tables
	if autoMigrate {
		err = migration.ServiceAutoMigration(orm.DB)
	}
	log.Info("[ORM] Database connection initialized.")
	return orm, err
}
