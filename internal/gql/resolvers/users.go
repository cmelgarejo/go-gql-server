package resolvers

import (
	"context"
	"errors"

	"github.com/cmelgarejo/go-gql-server/internal/orm"
	"github.com/cmelgarejo/go-gql-server/pkg/utils"

	"github.com/cmelgarejo/go-gql-server/pkg/utils/consts"

	"github.com/cmelgarejo/go-gql-server/internal/logger"

	"github.com/cmelgarejo/go-gql-server/internal/gql/models"
	tf "github.com/cmelgarejo/go-gql-server/internal/gql/resolvers/transformations"
	dbm "github.com/cmelgarejo/go-gql-server/internal/orm/models"
)

// CreateUser creates a record
func (r *mutationResolver) CreateUser(ctx context.Context, input models.UserInput) (*models.User, error) {
	cu := getCurrentUser(ctx)
	if ok, err := cu.HasPermission(consts.Permissions.Create, consts.EntityNames.Users); !ok || err != nil {
		return nil, logger.Errorfn(consts.EntityNames.Users, err)
	}
	return userCreateUpdate(r, input, false, cu)
}

// UpdateUser updates a record
func (r *mutationResolver) UpdateUser(ctx context.Context, id string, input models.UserInput) (*models.User, error) {
	cu := getCurrentUser(ctx)
	if ok, err := cu.HasPermission(consts.Permissions.Create, consts.EntityNames.Users); !ok || err != nil {
		return nil, logger.Errorfn(consts.EntityNames.Users, err)
	}
	return userCreateUpdate(r, input, true, cu, id)
}

// DeleteUser deletes a record
func (r *mutationResolver) DeleteUser(ctx context.Context, id string) (bool, error) {
	cu := getCurrentUser(ctx)
	if ok, err := cu.HasPermission(consts.Permissions.Delete, consts.EntityNames.Users); !ok || err != nil {
		return false, logger.Errorfn(consts.EntityNames.Users, err)
	}
	return userDelete(r, id)
}

// Users lists records
func (r *queryResolver) Users(ctx context.Context, id *string, filters []*models.QueryFilter, limit *int, offset *int, orderBy *string, sortDirection *string) (*models.Users, error) {
	cu := getCurrentUser(ctx)
	if ok, err := cu.HasPermission(consts.Permissions.List, consts.EntityNames.Users); !ok || err != nil {
		return nil, logger.Errorfn(consts.EntityNames.Users, err)
	}
	return userList(r, id, filters, limit, offset, orderBy, sortDirection)
}

// ## Helper functions

func userCreateUpdate(r *mutationResolver, input models.UserInput, update bool, cu *dbm.User, ids ...string) (*models.User, error) {
	dbo, err := tf.GQLInputUserToDBUser(&input, update, cu, ids...)
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
	tx = tx.Commit()
	return tf.DBUserToGQLUser(dbo), tx.Error
}

func userDelete(r *mutationResolver, id string) (bool, error) {
	return false, errors.New("not implemented")
}

func userList(r *queryResolver, id *string, filters []*models.QueryFilter, limit *int, offset *int, orderBy *string, sortDirection *string) (*models.Users, error) {
	whereID := "id = ?"
	record := &models.Users{}
	dbRecords := []*dbm.User{}
	tx := r.ORM.DB.Begin().
		Offset(*offset).Limit(*limit).Order(utils.ToSnakeCase(*orderBy) + " " + *sortDirection).
		Preload(consts.EntityNames.UserProfiles)
	if id != nil {
		tx = tx.Where(whereID, *id)
	}
	if filters != nil {
		if filtered, err := orm.ParseFilters(tx, filters); err == nil {
			tx = filtered
		} else {
			return nil, err
		}
	}
	tx = tx.Find(&dbRecords).Count(&record.Count)
	for _, dbRec := range dbRecords {
		record.List = append(record.List, tf.DBUserToGQLUser(dbRec))
	}
	return record, tx.Error
}
