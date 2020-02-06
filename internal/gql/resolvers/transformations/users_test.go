package transformations

import (
	"reflect"
	"testing"
	"time"

	"github.com/gofrs/uuid"

	gql "github.com/cmelgarejo/go-gql-server/internal/gql/models"
	dbm "github.com/cmelgarejo/go-gql-server/internal/orm/models"
	"github.com/markbates/goth"
)

var (
	provider  = "test"
	emptyStr  = ""
	gUUID, _  = uuid.NewV4()
	gUUID2, _ = uuid.NewV4()
	userID    = gUUID2.String()
	now       = time.Now()
	email     = "test@test.com"
	password  = func() string {
		u, _ := uuid.NewV4()
		return u.String()
	}()
	dbmUser       = &dbm.User{Email: email, Password: password}
	dbmUserNoPass = &dbm.User{
		Email:       email,
		Description: &emptyStr,
		FirstName:   &emptyStr,
		LastName:    &emptyStr,
		Location:    &emptyStr,
		Name:        &emptyStr,
		NickName:    &emptyStr,
		AvatarURL:   &emptyStr,
	}
	dbmUserIDNoTS = &dbm.User{
		BaseModelSoftDelete: dbm.BaseModelSoftDelete{
			BaseModel: dbm.BaseModel{ID: gUUID},
		}, Email: email, Password: password,
	}
	dbmUserIDNoTSOrPassFull = &dbm.User{
		BaseModelSoftDelete: dbm.BaseModelSoftDelete{
			BaseModel: dbm.BaseModel{ID: gUUID},
		},
		Email:       email,
		Description: &emptyStr,
		FirstName:   &emptyStr,
		LastName:    &emptyStr,
		Location:    &emptyStr,
		Name:        &emptyStr,
		NickName:    &emptyStr,
		AvatarURL:   &emptyStr,
	}
	dbmUserID = &dbm.User{
		BaseModelSoftDelete: dbm.BaseModelSoftDelete{
			BaseModel: dbm.BaseModel{
				ID: gUUID, CreatedAt: &now, UpdatedAt: &now,
			},
		}, Email: email, Password: password,
		UserProfiles: []dbm.UserProfile{*userProfile},
	}
	dbmUserNoID = &dbm.User{
		BaseModelSoftDelete: dbm.BaseModelSoftDelete{
			BaseModel: dbm.BaseModel{CreatedAt: &now, UpdatedAt: &now},
		},
		Email: email, Password: password,
	}
	gothUser = &goth.User{
		Email:       email,
		Provider:    provider,
		Description: emptyStr,
		FirstName:   emptyStr,
		LastName:    emptyStr,
		Location:    emptyStr,
		Name:        emptyStr,
		NickName:    emptyStr,
		AvatarURL:   emptyStr,
		UserID:      userID,
	}
	userProfile = &dbm.UserProfile{
		BaseModelSeq: dbm.BaseModelSeq{
			ID:        0,
			CreatedAt: &now,
			UpdatedAt: &now,
		},
		Email:          email,
		ExternalUserID: gUUID2.String(),
		Provider:       provider,
		Description:    emptyStr,
		FirstName:      emptyStr,
		LastName:       emptyStr,
		Location:       emptyStr,
		Name:           emptyStr,
		NickName:       emptyStr,
	}
	userProfileNoID = func() *dbm.UserProfile {
		up := userProfile
		up.ID = 0
		up.CreatedAt = nil
		up.UpdatedAt = nil
		return up
	}()
)

func TestDBUserToGQLUser(t *testing.T) {
	type args struct {
		i *dbm.User
	}
	profiles := []*gql.UserProfile{
		{
			ID:             0,
			CreatedAt:      now,
			UpdatedAt:      &now,
			Email:          email,
			AvatarURL:      &emptyStr,
			ExternalUserID: &userID,
			Name:           &emptyStr,
			FirstName:      &emptyStr,
			LastName:       &emptyStr,
			NickName:       &emptyStr,
			Description:    &emptyStr,
			Location:       &emptyStr,
		},
	}
	tests := []struct {
		name    string
		args    args
		wantO   *gql.User
		wantErr bool
	}{
		{
			name: "DBUser OK",
			args: args{
				i: dbmUserID,
			},
			wantO: &gql.User{
				ID:        gUUID.String(),
				Email:     email,
				CreatedAt: &now,
				UpdatedAt: &now,
				Profiles:  profiles,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotO, err := DBUserToGQLUser(tt.args.i)
			if (err != nil) != tt.wantErr {
				t.Errorf("DBUserToGQLUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotO, tt.wantO) {
				t.Errorf("DBUserToGQLUser() = \n%#v\n, want \n%#v\n", gotO.Profiles[0], tt.wantO.Profiles[0])
			}
		})
	}
}

func TestGQLInputUserToDBUser(t *testing.T) {
	type args struct {
		i      *gql.UserInput
		update bool
		ids    []string
	}
	tests := []struct {
		name    string
		args    args
		wantO   *dbm.User
		wantErr bool
	}{
		{
			name: "GQLInput Create OK",
			args: args{
				i: &gql.UserInput{
					Email:    &email,
					Password: &password,
				},
				update: false,
			},
			wantO: dbmUser,
		},
		{
			name: "GQLInput Update OK",
			args: args{
				i: &gql.UserInput{
					Email:    &email,
					Password: &password,
				},
				update: true,
				ids:    []string{gUUID.String()},
			},
			wantO: dbmUserIDNoTS,
		},
		{
			name: "GQLInput Update FAIL",
			args: args{
				i: &gql.UserInput{
					Email:    &email,
					Password: &password,
				},
				update: true,
				ids:    []string{"badID"},
			},
			wantErr: true,
		},
		{
			name: "GQLInput Create w/no Email FAIL",
			args: args{
				i: &gql.UserInput{},
			},
			wantErr: true,
		},
		{
			name: "GQLInput Create w/no Password FAIL",
			args: args{
				i: &gql.UserInput{
					Email: &email,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotO, err := GQLInputUserToDBUser(tt.args.i, tt.args.update, tt.args.ids...)
			if (err != nil) != tt.wantErr {
				t.Errorf("GQLInputUserToDBUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotO, tt.wantO) {
				t.Errorf("GQLInputUserToDBUser() = \n%v, want: \n%v", gotO, tt.wantO)
			}
		})
	}
}

func TestGothUserToDBUser(t *testing.T) {
	type args struct {
		i      *goth.User
		update bool
		ids    []string
	}
	tests := []struct {
		name    string
		args    args
		wantO   *dbm.User
		wantErr bool
	}{
		{
			name: "GothUser OK",
			args: args{
				i: gothUser,
			},
			wantO: dbmUserNoPass,
		},
		{
			name: "GothUser w/no Email FAIL",
			args: args{
				i: &goth.User{},
			},
			wantErr: true,
		},
		{
			name: "GothUser w/UUID OK",
			args: args{
				i:   gothUser,
				ids: []string{gUUID.String()},
			},
			wantO: dbmUserIDNoTSOrPassFull,
		},
		{
			name: "GothUser w/bad UUID string FAIL",
			args: args{
				i:   gothUser,
				ids: []string{"badID"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotO, err := GothUserToDBUser(tt.args.i, tt.args.update, tt.args.ids...)
			if (err != nil) != tt.wantErr {
				t.Errorf("GothUserToDBUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotO, tt.wantO) {
				t.Errorf("GothUserToDBUser() = \n%v, want: \n%v", gotO, tt.wantO)
			}
		})
	}
}

func TestGothUserToDBUserProfile(t *testing.T) {
	type args struct {
		i      *goth.User
		update bool
		ids    []int
	}
	tests := []struct {
		name    string
		args    args
		wantO   *dbm.UserProfile
		wantErr bool
	}{
		{
			name: "Create GothUser to DBUser OK",
			args: args{
				update: false,
				i:      gothUser,
			},
			wantO: userProfileNoID,
		},
		{
			name: "Update GothUser to DBUser OK",
			args: args{
				update: true,
				i:      gothUser,
				ids:    []int{0},
			},
			wantO: userProfile,
		},
		{
			name: "GothUser to DBUser UserID FAIL",
			args: args{
				update: false,
				i:      &goth.User{},
			},
			wantErr: true,
		},
		{
			name: "GothUser to DBUser Email FAIL",
			args: args{
				update: false,
				i:      &goth.User{UserID: "bah"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotO, err := GothUserToDBUserProfile(tt.args.i, tt.args.update, tt.args.ids...)
			if (err != nil) != tt.wantErr {
				t.Errorf("GothUserToDBUserProfile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotO, tt.wantO) {
				t.Errorf("GothUserToDBUserProfile() = \n%v, want \n%v", gotO, tt.wantO)
			}
		})
	}
}
