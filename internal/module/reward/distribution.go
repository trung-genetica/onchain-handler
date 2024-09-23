package reward

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/genefriendway/onchain-handler/conf"
	"github.com/genefriendway/onchain-handler/contracts/lptoken"
	util "github.com/genefriendway/onchain-handler/internal/utils/ethereum"
	"github.com/genefriendway/onchain-handler/internal/utils/log"
)

// DistributeReward distributes reward LifePoint tokens from the reward address to user wallets
func DistributeReward(config *conf.Configuration, recipients map[string]*big.Int) error {
	// Load Blockchain configuration
	rpcUrl := config.Blockchain.RpcUrl
	chainID := config.Blockchain.ChainID
	privateKey := config.Blockchain.PrivateKeyReward
	rewardAddress := config.Blockchain.RewardAddress

	// Connect to the blockchain network
	client, err := util.ConnectToNetwork(rpcUrl)
	if err != nil {
		return fmt.Errorf("failed to connect to network: %w", err)
	}

	// Get authentication for signing transactions
	privateKeyECDSA, err := util.PrivateKeyFromHex(privateKey)
	if err != nil {
		return fmt.Errorf("failed to get private key: %w", err)
	}

	auth, err := util.GetAuth(client, privateKeyECDSA, new(big.Int).SetUint64(uint64(chainID)))
	if err != nil {
		return fmt.Errorf("failed to get auth: %w", err)
	}

	// Set up the reward token contract instance
	tokenAddress := common.HexToAddress(rewardAddress)
	LPToken, err := lptoken.NewLptoken(tokenAddress, client)
	if err != nil {
		return fmt.Errorf("failed to instantiate ERC20 contract: %w", err)
	}

	// Distribute tokens to each recipient
	for recipientAddress, amount := range recipients {
		log.LG.Infof("Transferring %s tokens to %s...\n", amount.String(), recipientAddress)

		recipient := common.HexToAddress(recipientAddress)
		tx, err := LPToken.Transfer(auth, recipient, amount)
		if err != nil {
			log.LG.Infof("Failed to transfer tokens to %s: %v", recipientAddress, err)
			continue
		}

		// Log the transaction hash for tracking
		log.LG.Infof("Tokens transferred to %s. Tx hash: %s\n", recipientAddress, tx.Hash().Hex())
	}
	return nil
}
