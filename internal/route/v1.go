package route

import (
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"

	"github.com/genefriendway/onchain-handler/conf"
)

func RegisterRoutes(r *gin.Engine, config *conf.Configuration, db *gorm.DB) {
	v1 := r.Group("/api/v1")
	appRouter := v1.Group("")
	appRouter.GET("/healthcheck", nil)
}
