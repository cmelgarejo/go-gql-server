package jobs

import (
	"github.com/cmelgarejo/go-gql-server/internal/orm/models"
	"github.com/jinzhu/gorm"
	"gopkg.in/gormigrate.v1"
)

var (
	uname       = "Test User"
	fname       = "Test"
	lname       = "User"
	nname       = "Foo Bar"
	description = "This is the first user ever!"
	location    = "His house, maybe?"
	users       = []*models.User{
		&models.User{
			Email:       "test@test.com",
			Name:        &uname,
			FirstName:   &fname,
			LastName:    &lname,
			NickName:    &nname,
			Description: &description,
			Location:    &location,
		},
		&models.User{
			Email:       "test2@test.com",
			Name:        &uname,
			FirstName:   &fname,
			LastName:    &lname,
			NickName:    &nname,
			Description: &description,
			Location:    &location,
		},
	}
)

// SeedUsers inserts the first users
var SeedUsers *gormigrate.Migration = &gormigrate.Migration{
	ID: "SEED_USERS",
	Migrate: func(db *gorm.DB) error {
		for _, u := range users {
			if err := db.Create(u).Error; err != nil {
				return err
			}
			if err := db.Create(&models.UserAPIKey{UserID: u.ID}).Error; err != nil {
				return err
			}
		}
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
