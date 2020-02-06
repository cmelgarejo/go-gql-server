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
		{
			Email:       "admin@test.com",
			Name:        &uname,
			FirstName:   &fname,
			LastName:    &lname,
			NickName:    &nname,
			Description: &description,
			Location:    &location,
			Roles:       []models.Role{{BaseModelSeq: models.BaseModelSeq{ID: 1}}},
		},
		{
			Email:       "user@test.com",
			Name:        &uname,
			FirstName:   &fname,
			LastName:    &lname,
			NickName:    &nname,
			Description: &description,
			Location:    &location,
			Roles:       []models.Role{{BaseModelSeq: models.BaseModelSeq{ID: 2}}},
		},
	}
)

// SeedUsers inserts the first users
var SeedUsers *gormigrate.Migration = &gormigrate.Migration{
	ID: "SEED_USERS",
	Migrate: func(db *gorm.DB) error {
		tx := db.Begin()
		defer tx.RollbackUnlessCommitted()
		for _, u := range users {
			if err := tx.Create(u).Error; err != nil {
				return err
			}
			if err := tx.Create(&models.UserAPIKey{UserID: u.ID}).Error; err != nil {
				return err
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
