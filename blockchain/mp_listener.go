package blockchain

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/genefriendway/onchain-handler/internal/utils/log"
)

type MPEventData struct {
	User    common.Address
	Amount  *big.Int
	OrderID *big.Int
	TxHash  string
}

// MPEventListener listens for the MembershipPurchased event and sends the event data to a channel.
func MPEventListener(
	client *ethclient.Client,
	contractAddr common.Address,
	parsedABI abi.ABI,
	eventChan chan<- MPEventData,
	errorChan chan<- error,
) {
	log.LG.Info("Starting event listener for MembershipPurchased...")

	// Subscribe to filter logs
	logs, sub, err := subscribeToLogs(client, contractAddr)
	if err != nil {
		errorChan <- err
		return
	}

	// Listen for logs and handle them
	for {
		select {
		case err := <-sub.Err():
			errorChan <- fmt.Errorf("subscription error: %w", err)
			return

		case vLog := <-logs:
			eventData, err := parseEventLog(vLog, parsedABI)
			if err != nil {
				errorChan <- err
				continue
			}

			// Send the event data to the event channel
			eventChan <- eventData
		}
	}
}

// subscribeToLogs creates a log subscription for the contract and returns the log channel and subscription.
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

// parseEventLog unpacks the log data into the MPEventData structure.
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
