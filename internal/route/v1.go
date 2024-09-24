package route

import (
	"context"

	"gorm.io/gorm"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-gonic/gin"

	"github.com/genefriendway/onchain-handler/blockchain"
	"github.com/genefriendway/onchain-handler/conf"
	"github.com/genefriendway/onchain-handler/internal/module/membership"
	"github.com/genefriendway/onchain-handler/internal/module/reward"
)

func RegisterRoutes(r *gin.Engine, config *conf.Configuration, db *gorm.DB, ethClient *ethclient.Client) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	v1 := r.Group("/api/v1")
	appRouter := v1.Group("")

	// SECTION: reward tokens
	rewardRepository := reward.NewRewardRepository(db)
	rewardUCase := reward.NewRewardUCase(rewardRepository, ethClient, config)
	rewardHandler := reward.NewRewardHandler(rewardUCase)
	appRouter.POST("/rewards", rewardHandler.Reward)

	// SECTION: membership purchase
	membershipRepository := membership.NewMembershipRepository(db)

	// SECTION: events listener
	membershipEventListener := blockchain.NewMembershipEventListener(
		ethClient,
		config.Blockchain.MembershipContractAddress,
		membershipRepository,
	)
	go membershipEventListener.RunListener(ctx)
}
