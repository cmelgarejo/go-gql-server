package transformations

import (
	"errors"

	gql "github.com/cmelgarejo/go-gql-server/internal/gql/models"
	dbm "github.com/cmelgarejo/go-gql-server/internal/orm/models"
	"github.com/gofrs/uuid"
)

// DBUserToGQLUser transforms [user] db input to gql type
func DBUserToGQLUser(i *dbm.User) (o *gql.User, err error) {
	o = &gql.User{
		ID:          i.ID.String(),
		Email:       i.Email,
		UserID:      i.UserID,
		Name:        i.Name,
		FirstName:   i.FirstName,
		LastName:    i.LastName,
		NickName:    i.NickName,
		Description: i.Description,
		Location:    i.Location,
		CreatedAt:   i.CreatedAt,
		UpdatedAt:   i.UpdatedAt,
	}
	return o, err
}

// GQLInputUserToDBUser transforms [user] gql input to db model
func GQLInputUserToDBUser(i *gql.UserInput, update bool, ids ...string) (o *dbm.User, err error) {
	o = &dbm.User{
		UserID:      i.UserID,
		Name:        i.Name,
		FirstName:   i.FirstName,
		LastName:    i.LastName,
		NickName:    i.NickName,
		Description: i.Description,
		Location:    i.Location,
	}
	if i.Email == nil && !update {
		return nil, errors.New("Field [email] is required")
	}
	if i.Email != nil {
		o.Email = *i.Email
	}
	if len(ids) > 0 {
		updID, err := uuid.FromString(ids[0])
		if err != nil {
			return nil, err
		}
		o.ID = updID
	}
	return o, err
}
