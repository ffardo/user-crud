package routes

import (
	"net/http"

	"github.com/ffardo/user-crud/interfaces"
	"github.com/gin-gonic/gin"
)

func InitRouter(uc interfaces.UserController, apiKey string) *gin.Engine {

	r := gin.Default()
	g := r.Group("/api/users", func(ctx *gin.Context) {
		if ctx.Request.Header.Get("X-API-KEY") != apiKey {
			ctx.AbortWithStatus(http.StatusUnauthorized)
		}
	})
	g.POST("/", uc.Post)
	g.GET("/:uuid", uc.Get)
	g.PATCH("/:uuid", uc.Patch)
	g.DELETE("/:uuid", uc.Delete)

	return r
}
