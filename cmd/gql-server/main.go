package main

import (
	log "github.com/cmelgarejo/go-gql-server/internal/logger"

	"github.com/cmelgarejo/go-gql-server/internal/orm"
	"github.com/cmelgarejo/go-gql-server/pkg/server"
)

func main() {
	orm, err := orm.Factory()
	if err != nil {
		log.Panic(err)
	}
	server.Run(orm)
}
