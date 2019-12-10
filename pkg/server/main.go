package server

import (
	"github.com/cmelgarejo/go-gql-server/internal/logger"
	"github.com/cmelgarejo/go-gql-server/internal/orm"
	"github.com/cmelgarejo/go-gql-server/pkg/utils"
	"github.com/gin-gonic/gin"
)

// Run spins up the server
func Run(serverconf *utils.ServerConfig, orm *orm.ORM) {
	r := gin.Default()

	// Initialize the Auth providers
	InitalizeAuthProviders(serverconf)

	// Routes and Handlers
	RegisterRoutes(serverconf, r, orm)

	// Inform the user where the server is listening
	logger.Info("Running @ " + serverconf.SchemaVersionedEndpoint(""))

	// Run the server
	// Print out and exit(1) to the OS if the server cannot run
	logger.Fatal(r.Run(serverconf.ListenEndpoint()))
}
