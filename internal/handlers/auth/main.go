package auth

import (
	"context"
	"net/http"

	"github.com/cmelgarejo/go-gql-server/pkg/utils"
	"github.com/gin-gonic/gin"
)

func addProviderToContext(c *gin.Context, value interface{}) *http.Request {
	return c.Request.WithContext(context.WithValue(c.Request.Context(),
		string(utils.ProjectContextKeys.GothicProviderCtxKey), value))
}
