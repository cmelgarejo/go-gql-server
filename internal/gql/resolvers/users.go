package resolvers

import (
	"context"
	"errors"

	"github.com/cmelgarejo/go-gql-server/pkg/utils/consts"

	"github.com/cmelgarejo/go-gql-server/internal/logger"

	"github.com/cmelgarejo/go-gql-server/internal/gql/models"
	tf "github.com/cmelgarejo/go-gql-server/internal/gql/resolvers/transformations"
	dbm "github.com/cmelgarejo/go-gql-server/internal/orm/models"
)

// CreateUser creates a record
func (r *mutationResolver) CreateUser(ctx context.Context, input models.UserInput) (*models.User, error) {
	return userCreateUpdate(r, input, false)
}

// UpdateUser updates a record
func (r *mutationResolver) UpdateUser(ctx context.Context, id string, input models.UserInput) (*models.User, error) {
	return userCreateUpdate(r, input, true, id)
}

// DeleteUser deletes a record
func (r *mutationResolver) DeleteUser(ctx context.Context, id string) (bool, error) {
	return userDelete(r, id)
}

// Users lists records
func (r *queryResolver) Users(ctx context.Context, id *string) (*models.Users, error) {
	cu := getCurrentUser(ctx)
	if ok, err := cu.HasPermission(consts.Permissions.List, consts.GetTableName(consts.EntityNames.Users)); !ok || err != nil {
		return nil, logger.Errorfn(consts.EntityNames.Users, err)
	}
	return userList(r, id)
}

// ## Helper functions

func userCreateUpdate(r *mutationResolver, input models.UserInput, update bool, ids ...string) (*models.User, error) {
	dbo, err := tf.GQLInputUserToDBUser(&input, update, ids...)
	if err != nil {
		return nil, err
	}
	// Create scoped clean db interface
	tx := r.ORM.DB.Begin()
	defer tx.RollbackUnlessCommitted()
	if !update {
		tx = tx.Create(dbo).First(dbo) // Create the user
		if tx.Error != nil {
			return nil, tx.Error
		}
	} else {
		tx = tx.Model(&dbo).Update(dbo).First(dbo) // Or update it
	}
	gql, err := tf.DBUserToGQLUser(dbo)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx = tx.Commit()
	return gql, tx.Error
}

func userDelete(r *mutationResolver, id string) (bool, error) {
	return false, errors.New("not implemented")
}

func userList(r *queryResolver, id *string) (*models.Users, error) {
	entity := consts.GetTableName(consts.EntityNames.Users)
	whereID := "id = ?"
	record := &models.Users{}
	dbRecords := []*dbm.User{}
	tx := r.ORM.DB.Begin()
	defer tx.RollbackUnlessCommitted()
	if id != nil {
		tx = tx.Where(whereID, *id)
	}
	tx = tx.Find(&dbRecords).Count(&record.Count)
	for _, dbRec := range dbRecords {
		if rec, err := tf.DBUserToGQLUser(dbRec); err != nil {
			logger.Errorfn(entity, err)
		} else {
			record.List = append(record.List, rec)
		}
	}
	return record, tx.Error
}
