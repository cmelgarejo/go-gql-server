package resolvers

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"

	"github.com/cmelgarejo/go-gql-server/pkg/utils/consts"

	"github.com/cmelgarejo/go-gql-server/internal/logger"

	"github.com/cmelgarejo/go-gql-server/internal/gql/models"
	// tf "github.com/cmelgarejo/go-gql-server/internal/gql/resolvers/transformations"
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
func (r *queryResolver) Users(ctx context.Context, id *string, filters []*models.QueryFilter, limit *int, offset *int, orderBy *string, sortDirection *string) ([]*models.User, error) {
	return userList(r, id, filters, limit, offset, orderBy, sortDirection)
}

// ## Helper functions

func userCreateUpdate(r *mutationResolver, input models.UserInput, update bool, cu *dbm.User, ids ...string) (*models.User, error) {
	return nil, errors.New("not implemented")
}

func userDelete(r *mutationResolver, id string) (bool, error) {
	return false, errors.New("not implemented")
}

func userList(r *queryResolver, id *string, filters []*models.QueryFilter, limit *int, offset *int, orderBy *string, sortDirection *string) ([]*models.User, error) {
	usersURL := "https://jsonplaceholder.typicode.com/users"
	if id != nil {
		usersURL = usersURL + "/" + *id
	}
	resp, err := http.Get(usersURL)
	if err != nil {
		logger.Error(err)
	}
	users := []*models.User{}
	jsonResponse, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err)
	}
	if id != nil {
		user := models.User{}
		err = json.Unmarshal(jsonResponse, &user)
		users = append(users, &user)
	} else {
		err = json.Unmarshal(jsonResponse, &users)
	}
	if err != nil {
		logger.Error(err)
	}
	if orderBy != nil && sortDirection != nil {
		sort.Slice(users, func(i, j int) bool {
			if *sortDirection == "ASC" {
				return (strings.Compare(*users[i].Name, *users[j].Name) == -1)
			}
			return (strings.Compare(*users[i].Name, *users[j].Name) == 1)
		})
	}
	for _, filter := range filters {
		if filter.Field == "email" {
			filteredUsers := []*models.User{}
			for _, u := range users {
				if strings.Contains(u.Email, filter.Value.(string)) {
					filteredUsers = append(filteredUsers, u)
				}
			}
			return filteredUsers, nil
		}
	}
	posts := getAllPosts()
	for _, u := range users {
		for _, p := range posts {
			if p.UserID == u.ID {
				u.Posts = append(u.Posts, p)
			}
		}
	}
	return users, err
}

func getAllPosts() []*models.Post {
	postsURL := "https://jsonplaceholder.typicode.com/posts"
	resp, err := http.Get(postsURL)
	if err != nil {
		logger.Error(err)
	}
	records := []*models.Post{}
	jsonResponse, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err)
	}
	err = json.Unmarshal(jsonResponse, &records)
	if err != nil {
		logger.Error(err)
	}
	return records
}
