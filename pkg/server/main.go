package server

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/cmelgarejo/go-gql-server/internal/handlers"
	"github.com/cmelgarejo/go-gql-server/pkg/utils"
)

var host, port string

func init() {
	host = utils.MustGet("GQL_SERVER_HOST")
	port = utils.MustGet("GQL_SERVER_PORT")
}

// Run spins up the server
func Run() {
	r := gin.Default()
	// Simple keep-alive/ping handler
	r.GET("/ping", handlers.Ping())
	// Inform the user where the server is listening
	log.Println("Running @ http://" + host + ":" + port)
	// Print out and exit(1) to the OS if the server cannot run
	log.Fatalln(r.Run(host + ":" + port))

}
