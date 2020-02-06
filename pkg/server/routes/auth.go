package routes

import (
	"github.com/cmelgarejo/go-gql-server/internal/handlers/auth"
	"github.com/cmelgarejo/go-gql-server/internal/orm"
	"github.com/cmelgarejo/go-gql-server/pkg/utils"
	"github.com/gin-gonic/gin"
)

// Auth routes
func Auth(cfg *utils.ServerConfig, r *gin.Engine, orm *orm.ORM) error {
	provider := string(utils.ProjectContextKeys.ProviderCtxKey)
	// OAuth handlers
	g := r.Group(cfg.VersionedEndpoint("/auth"))
	g.GET("/:"+provider, auth.Begin())
	g.GET("/:"+provider+"/callback", auth.Callback(cfg, orm))
	// g.GET(:"+provider+"/refresh", auth.Refresh(cfg, orm))
	return nil
}
