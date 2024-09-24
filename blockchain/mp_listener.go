package blockchain

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"path/filepath"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/genefriendway/onchain-handler/internal/utils/log"
)

// MPEventData represents the event data for a MembershipPurchased event
type MPEventData struct {
	User    common.Address
	Amount  *big.Int
	OrderID *big.Int
	TxHash  string
}

// MPEventListener listens for the MembershipPurchased event and sends the event data to a channel
func MPEventListener(
	client *ethclient.Client,
	contractAddr common.Address,
	eventChan chan<- MPEventData,
) {
	log.LG.Info("Starting event listener for MembershipPurchased...")

	abiFilePath, err := filepath.Abs("./contracts/abis/MembershipPurchase.abi.json")
	if err != nil {
		log.LG.Errorf("Failed to get absolute path for ABI file: %v", err)
		return
	}

	// Load the ABI from the hardcoded file path
	parsedABI, err := loadABI(abiFilePath)
	if err != nil {
		log.LG.Errorf("Failed to load ABI: %v", err)
		return
	}

	// Subscribe to filter logs
	logs, sub, err := subscribeToLogs(client, contractAddr)
	if err != nil {
		log.LG.Errorf("Failed to subscribe to logs: %v", err)
		return
	}

	// Listen for logs and handle them
	for {
		select {
		case err := <-sub.Err():
			log.LG.Errorf("Subscription error: %v", err)
			return

		case vLog := <-logs:
			eventData, err := parseEventLog(vLog, parsedABI)
			if err != nil {
				log.LG.Errorf("Failed to unpack log: %v", err)
				continue
			}

			// Send the event data to the event channel
			eventChan <- eventData
		}
	}
}

// loadABI loads and parses the ABI from a JSON file
func loadABI(path string) (abi.ABI, error) {
	file, err := os.Open(path)
	if err != nil {
		return abi.ABI{}, fmt.Errorf("failed to open ABI file: %w", err)
	}
	defer file.Close()

	// Read and parse the ABI JSON file
	var parsedABI abi.ABI
	err = json.NewDecoder(file).Decode(&parsedABI)
	if err != nil {
		return abi.ABI{}, fmt.Errorf("failed to decode ABI JSON: %w", err)
	}

	return parsedABI, nil
}

// subscribeToLogs creates a log subscription for the contract and returns the log channel and subscription
func subscribeToLogs(
	client *ethclient.Client,
	contractAddr common.Address,
) (chan types.Log, ethereum.Subscription, error) {
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddr},
	}

	logs := make(chan types.Log)
	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to subscribe to logs: %w", err)
	}

	return logs, sub, nil
}

// parseEventLog unpacks the log data into the MPEventData structure
func parseEventLog(
	vLog types.Log,
	parsedABI abi.ABI,
) (MPEventData, error) {
	event := struct {
		User    common.Address
		Amount  *big.Int
		OrderID *big.Int
	}{}

	// Unpack the log data into the event structure
	err := parsedABI.UnpackIntoInterface(&event, "MembershipPurchased", vLog.Data)
	if err != nil {
		return MPEventData{}, fmt.Errorf("failed to unpack log: %w", err)
	}

	// Extract indexed fields (user address and order ID)
	event.User = common.HexToAddress(vLog.Topics[1].Hex())
	event.OrderID = new(big.Int)
	event.OrderID.SetString(vLog.Topics[2].Hex()[2:], 16)

	// Prepare event data
	eventData := MPEventData{
		User:    event.User,
		Amount:  event.Amount,
		OrderID: event.OrderID,
		TxHash:  vLog.TxHash.Hex(),
	}

	return eventData, nil
}
