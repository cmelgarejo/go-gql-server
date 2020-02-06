// Package orm provides `GORM` helpers for the creation, migration and access
// on the project's database
package orm

import (
	"errors"
	"fmt"

	"github.com/cmelgarejo/go-gql-server/pkg/utils/consts"

	"github.com/cmelgarejo/go-gql-server/internal/gql/resolvers/transformations"

	"github.com/markbates/goth"

	"github.com/cmelgarejo/go-gql-server/internal/logger"
	"github.com/cmelgarejo/go-gql-server/internal/orm/models"

	"github.com/cmelgarejo/go-gql-server/internal/orm/migration"

	"github.com/cmelgarejo/go-gql-server/pkg/utils"

	//Imports the database dialect of choice
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/jinzhu/gorm"
)

var (
	sUserTbl  = "User"
	nestedFmt = "%s.%s"
)

// ORM struct to holds the gorm pointer to db
type ORM struct {
	DB *gorm.DB
}

// Factory creates a db connection with the selected dialect and connection
// string
func Factory(cfg *utils.ServerConfig) (*ORM, error) {
	db, err := gorm.Open(cfg.Database.Dialect, cfg.Database.DSN)
	if err != nil {
		logger.Panic("[ORM] err: ", err)
	}
	orm := &ORM{DB: db}
	// Log every SQL command on dev, @prod: this should be disabled? Maybe.
	db.LogMode(cfg.Database.LogMode)
	// Automigrate tables
	if cfg.Database.AutoMigrate {
		err = migration.ServiceAutoMigration(orm.DB)
		if err != nil {
			logger.Error("[ORM.autoMigrate] err: ", err)
		}
	}
	logger.Info("[ORM] Database connection initialized.")
	return orm, nil
}

//FindUserByAPIKey finds the user that is related to the API key
func (o *ORM) FindUserByAPIKey(apiKey string) (*models.User, error) {
	if apiKey == "" {
		return nil, errors.New("API key is empty")
	}
	uak := &models.UserAPIKey{}
	up := fmt.Sprintf(nestedFmt, sUserTbl, consts.EntityNames.Permissions)
	ur := fmt.Sprintf(nestedFmt, sUserTbl, consts.EntityNames.Roles)
	if err := o.DB.Preload(sUserTbl).Preload(up).Preload(ur).
		Where("api_key = ?", apiKey).Find(uak).Error; err != nil {
		return nil, err
	}
	return &uak.User, nil
}

// FindUserByJWT finds the user that is related to the APIKey token
func (o *ORM) FindUserByJWT(email string, provider string, userID string) (*models.User, error) {
	if provider == "" || userID == "" {
		return nil, errors.New("provider or userId empty")
	}
	tx := o.DB.Begin()
	p := &models.UserProfile{}
	up := fmt.Sprintf(nestedFmt, sUserTbl, consts.EntityNames.Permissions)
	ur := fmt.Sprintf(nestedFmt, sUserTbl, consts.EntityNames.Roles)
	if err := tx.Preload(sUserTbl).Preload(up).Preload(ur).
		Where("email  = ? AND provider = ? AND external_user_id = ?", email, provider, userID).
		First(p).Error; err != nil {
		return nil, err
	}
	return &p.User, nil
}

// UpsertUserProfile saves the user if doesn't exists and adds the OAuth profile
func (o *ORM) UpsertUserProfile(input *goth.User) (*models.User, error) {
	db := o.DB.New()
	up := &models.UserProfile{}
	u, err := transformations.GothUserToDBUser(input, false)
	if err != nil {
		return nil, err
	}
	if tx := db.Where("email = ?", input.Email).First(u); !tx.RecordNotFound() && tx.Error != nil {
		return nil, tx.Error
	}
	if tx := db.Model(u).Save(u); tx.Error != nil {
		return nil, err
	}
	if tx := db.Where("email = ? AND provider = ? AND external_user_id = ?",
		input.Email, input.Provider, input.UserID).First(up); !tx.RecordNotFound() && tx.Error != nil {
		return nil, err
	}
	up, err = transformations.GothUserToDBUserProfile(input, false)
	if err != nil {
		return nil, err
	}
	up.User = *u
	if tx := db.Model(up).Save(up); tx.Error != nil {
		return nil, tx.Error
	}
	return u, nil
}
