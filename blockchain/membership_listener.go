package blockchain

import (
	"context"
	"fmt"
	"math/big"
	"path/filepath"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/genefriendway/onchain-handler/internal/interfaces"
	"github.com/genefriendway/onchain-handler/internal/model"
	"github.com/genefriendway/onchain-handler/internal/utils/log"
)

// MembershipEventData represents the event data for a MembershipPurchased event.
type MembershipEventData struct {
	User     common.Address
	Amount   *big.Int
	OrderID  uint64
	TxHash   string
	Duration uint // Duration as an integer representing the type (0 for 1 year, 1 for 3 years)
}

// MembershipEventListener listens for MembershipPurchased events and processes them.
type MembershipEventListener struct {
	ETHClient       *ethclient.Client
	ContractAddress common.Address // Changed to common.Address for better type safety
	EventChan       chan MembershipEventData
	Repo            interfaces.MembershipRepository
	ParsedABI       abi.ABI // Pre-loaded ABI for efficiency
}

// NewMembershipEventListener initializes the listener and preloads the ABI.
func NewMembershipEventListener(
	client *ethclient.Client,
	contractAddr string,
	repo interfaces.MembershipRepository,
) (MembershipEventListener, error) {
	// Load the ABI once during initialization.
	abiFilePath, err := filepath.Abs("./contracts/abis/MembershipPurchase.abi.json")
	if err != nil {
		return MembershipEventListener{}, fmt.Errorf("failed to get ABI file path: %w", err)
	}

	parsedABI, err := loadABI(abiFilePath)
	if err != nil {
		return MembershipEventListener{}, fmt.Errorf("failed to load ABI: %w", err)
	}

	// Initialize the event channel with a buffer of 25.
	eventChan := make(chan MembershipEventData, 25)

	return MembershipEventListener{
		ETHClient:       client,
		ContractAddress: common.HexToAddress(contractAddr), // Convert string to common.Address
		EventChan:       eventChan,
		Repo:            repo,
		ParsedABI:       parsedABI, // Preloaded ABI
	}, nil
}

// RunListener starts the listener and processes incoming events.
func (listener MembershipEventListener) RunListener() {
	// Start listening for events.
	go listener.listen()

	// Handle incoming events from the channel.
	go func() {
		for event := range listener.EventChan {
			// Map the duration to days.
			var durationDays int
			switch event.Duration {
			case 0:
				durationDays = 365
			case 1:
				durationDays = 1095
			default:
				log.LG.Errorf("Unexpected duration: %d", event.Duration)
				continue // Skip this event or handle it as needed.
			}

			// Calculate the end duration based on the duration in days.
			endDuration := time.Now().AddDate(0, 0, durationDays)

			// Create the MembershipEvents model.
			membershipEvent := model.MembershipEvents{
				UserAddress:     event.User.Hex(),
				OrderID:         event.OrderID,
				TransactionHash: event.TxHash,
				Amount:          event.Amount.String(),
				Status:          1, // Assuming 1 means active or successful
				EndDuration:     endDuration,
			}

			// Save the event to the repository.
			if err := listener.Repo.CreateMembershipEventHistory(context.Background(), membershipEvent); err != nil {
				log.LG.Errorf("Error saving event to database: %v", err)
			} else {
				log.LG.Infof("Event processed: %+v", membershipEvent)
			}
		}
	}()

	select {}
}

// listen polls the blockchain for logs and parses them.
func (listener MembershipEventListener) listen() {
	log.LG.Info("Starting event listener for MembershipPurchased...")

	// Define how far back to start polling.
	blockOffset := int64(10)

	// Poll for logs using the passed context instead of creating a new one.
	logs, err := pollForLogs(listener.ETHClient, listener.ContractAddress, blockOffset)
	if err != nil {
		log.LG.Errorf("Failed to poll logs from contract %s: %v", listener.ContractAddress.Hex(), err)
		return
	}

	// Process each log entry.
	for _, vLog := range logs {
		eventData, err := listener.parseEventLog(vLog, listener.ParsedABI)
		if err != nil {
			log.LG.Errorf("Failed to parse log for TxHash %s: %v", vLog.TxHash.Hex(), err)
			continue
		}

		// Send the parsed event data to the event channel.
		listener.EventChan <- eventData
	}
}

// parseEventLog unpacks the log data into the MembershipEventData structure.
func (listener MembershipEventListener) parseEventLog(
	vLog types.Log,
	parsedABI abi.ABI,
) (MembershipEventData, error) {
	event := struct {
		User     common.Address
		Amount   *big.Int
		OrderID  uint64
		Duration uint // Duration passed as part of the event
	}{}

	// Unpack the log data into the event structure.
	err := parsedABI.UnpackIntoInterface(&event, "MembershipPurchased", vLog.Data)
	if err != nil {
		return MembershipEventData{}, fmt.Errorf("failed to unpack log for TxHash %s: %w", vLog.TxHash.Hex(), err)
	}

	// Extract indexed fields (user address and order ID).
	event.User = common.HexToAddress(vLog.Topics[1].Hex())

	// Parse the OrderID from hex to uint64.
	orderID, err := parseHexToUint64(vLog.Topics[2].Hex())
	if err != nil {
		return MembershipEventData{}, fmt.Errorf("failed to parse order ID for TxHash %s: %w", vLog.TxHash.Hex(), err)
	}

	// Prepare event data, including the parsed duration.
	eventData := MembershipEventData{
		User:     event.User,
		Amount:   event.Amount,
		OrderID:  orderID,
		Duration: event.Duration,
		TxHash:   vLog.TxHash.Hex(),
	}

	return eventData, nil
}
