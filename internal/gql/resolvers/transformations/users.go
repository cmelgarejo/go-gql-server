package transformations

import (
	"errors"

	"github.com/markbates/goth"

	gql "github.com/cmelgarejo/go-gql-server/internal/gql/models"
	dbm "github.com/cmelgarejo/go-gql-server/internal/orm/models"
	"github.com/gofrs/uuid"
)

// DBUserToGQLUser transforms [user] db input to gql type
func DBUserToGQLUser(i *dbm.User) *gql.User {
	if i == nil {
		return nil
	}
	profiles := []*gql.UserProfile{}
	for _, p := range i.UserProfiles {
		profiles = append(profiles, DBUserProfileToGQLUserProfile(&p))
	}
	return &gql.User{
		AvatarURL:   i.AvatarURL,
		ID:          i.ID.String(),
		Email:       i.Email,
		Name:        i.Name,
		FirstName:   i.FirstName,
		LastName:    i.LastName,
		NickName:    i.NickName,
		Description: i.Description,
		Location:    i.Location,
		Profiles:    profiles,
		CreatedAt:   i.CreatedAt,
		UpdatedAt:   i.UpdatedAt,
	}
}

// DBUserProfileToGQLUserProfile transforms [user] db input to gql type
func DBUserProfileToGQLUserProfile(i *dbm.UserProfile) *gql.UserProfile {
	if i == nil {
		return nil
	}
	return &gql.UserProfile{
		AvatarURL:      &i.AvatarURL,
		ID:             i.ID,
		ExternalUserID: &i.ExternalUserID,
		Email:          i.Email,
		Name:           &i.Name,
		FirstName:      &i.FirstName,
		LastName:       &i.LastName,
		NickName:       &i.NickName,
		Description:    &i.Description,
		Location:       &i.Location,
		CreatedAt:      *i.CreatedAt,
		UpdatedAt:      i.UpdatedAt,
		CreatedBy:      DBUserToGQLUser(i.CreatedBy),
		UpdatedBy:      DBUserToGQLUser(i.UpdatedBy),
	}
}

// GQLInputUserToDBUser transforms [user] gql input to db model
func GQLInputUserToDBUser(i *gql.UserInput, update bool, u *dbm.User, ids ...string) (o *dbm.User, err error) {
	if i.Email == nil && !update {
		return nil, errors.New("field [email] is required")
	}
	if i.Password == nil && !update {
		return nil, errors.New("field [password] is required")
	}
	o = &dbm.User{
		Name:        i.Name,
		FirstName:   i.FirstName,
		LastName:    i.LastName,
		NickName:    i.NickName,
		Description: i.Description,
		Location:    i.Location,
	}
	if i.Email != nil {
		o.Email = *i.Email
	}
	if i.Password != nil {
		o.Password = *i.Password
	}
	if !update {
		o.CreatedBy = u
	}
	o.UpdatedBy = u
	if len(ids) > 0 {
		updID, err := uuid.FromString(ids[0])
		if err != nil {
			return nil, err
		}
		o.ID = updID
	}
	return o, err
}

// GothUserToDBUser transforms [user] goth to db model
func GothUserToDBUser(i *goth.User, update bool, ids ...string) (o *dbm.User, err error) {
	if i.Email == "" && !update {
		return nil, errors.New("field [Email] is required")
	}
	o = &dbm.User{
		Email:       i.Email,
		Name:        &i.Name,
		FirstName:   &i.FirstName,
		LastName:    &i.LastName,
		NickName:    &i.NickName,
		Location:    &i.Location,
		AvatarURL:   &i.AvatarURL,
		Description: &i.Description,
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

// GothUserToDBUserProfile transforms [user] goth to db model
func GothUserToDBUserProfile(i *goth.User, update bool, ids ...int) (o *dbm.UserProfile, err error) {
	if i.UserID == "" && !update {
		return nil, errors.New("field [UserID] is required")
	}
	if i.Email == "" && !update {
		return nil, errors.New("field [Email] is required")
	}
	o = &dbm.UserProfile{
		ExternalUserID: i.UserID,
		Provider:       i.Provider,
		Email:          i.Email,
		Name:           i.Name,
		FirstName:      i.FirstName,
		LastName:       i.LastName,
		NickName:       i.NickName,
		Location:       i.Location,
		AvatarURL:      i.AvatarURL,
		Description:    i.Description,
	}
	if len(ids) > 0 {
		updID := ids[0]
		o.ID = updID
	}
	return o, err
}
