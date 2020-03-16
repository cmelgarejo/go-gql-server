package main

import (
	"github.com/cmelgarejo/go-gql-server/cmd/gql-server/config"
	"github.com/cmelgarejo/go-gql-server/internal/logger"

	"github.com/cmelgarejo/go-gql-server/internal/orm"
	"github.com/cmelgarejo/go-gql-server/pkg/server"
)

// main
func main() {
	sc := config.Server()
	orm, err := orm.Factory(sc)
	defer orm.DB.Close()
	if err != nil {
		logger.Panic(err)
	}
	server.Run(sc, orm)
}
