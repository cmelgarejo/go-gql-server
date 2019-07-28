package jobs

import (
	"github.com/cmelgarejo/go-gql-server/internal/orm/models"
	"github.com/jinzhu/gorm"
	"gopkg.in/gormigrate.v1"
)

var (
	uname                    = "Test User"
	fname                    = "Test"
	lname                    = "User"
	nname                    = "Foo Bar"
	description              = "This is the first user ever!"
	location                 = "His house, maybe?"
	firstUser   *models.User = &models.User{
		Email:       "test@test.com",
		Name:        &uname,
		FirstName:   &fname,
		LastName:    &lname,
		NickName:    &nname,
		Description: &description,
		Location:    &location,
	}
)

// SeedUsers inserts the first users
var SeedUsers *gormigrate.Migration = &gormigrate.Migration{
	ID: "SEED_USERS",
	Migrate: func(db *gorm.DB) error {
		return db.Create(&firstUser).Error
	},
	Rollback: func(db *gorm.DB) error {
		return db.Delete(&firstUser).Error
	},
}
