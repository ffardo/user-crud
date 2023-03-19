package interfaces

import "github.com/gin-gonic/gin"

type UserController interface {
	Post(ctx *gin.Context)
	Get(ctx *gin.Context)
	Patch(ctx *gin.Context)
	Delete(ctx *gin.Context)
}
