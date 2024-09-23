package route

import (
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"

	"github.com/genefriendway/onchain-handler/conf"
	"github.com/genefriendway/onchain-handler/internal/module/reward"
)

func RegisterRoutes(r *gin.Engine, config *conf.Configuration, db *gorm.DB) {
	v1 := r.Group("/api/v1")
	appRouter := v1.Group("")

	// SECTION: reward tokens
	rewardRepository := reward.NewRewardRepository(db)
	rewardUCase := reward.NewRewardUCase(rewardRepository)
	rewardHandler := reward.NewRewardHandler(rewardUCase, config)
	appRouter.POST("/rewards", rewardHandler.Reward)
}
