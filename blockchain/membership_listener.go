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

// MembershipEventData represents the event data for a MembershipPurchased event
type MembershipEventData struct {
	User     common.Address
	Amount   *big.Int
	OrderID  uint64
	TxHash   string
	Duration uint
}

// MembershipEventListener listens for the MembershipPurchased event and sends the event data to a channel
type MembershipEventListener struct {
	ETHClient        *ethclient.Client
	ContractAddrress string
	EventChan        chan MembershipEventData
	Repo             interfaces.MembershipRepository
}

func NewMembershipEventListener(
	client *ethclient.Client,
	contractAddr string,
	repo interfaces.MembershipRepository,
) MembershipEventListener {
	// Initialize the event channel with a buffer of 25
	eventChan := make(chan MembershipEventData, 25)

	return MembershipEventListener{
		ETHClient:        client,
		ContractAddrress: contractAddr,
		EventChan:        eventChan, // Channel is initialized here
		Repo:             repo,
	}
}

func (listener MembershipEventListener) RunListener(ctx context.Context) {
	// Start listening for events
	go listener.listen()

	// Handle incoming events from the channel
	go func() {
		for event := range listener.EventChan {
			// Convert the event to model.MembershipEvents
			var durationDays int
			switch event.Duration {
			case 0:
				durationDays = 365
			case 1:
				durationDays = 1095
			default:
				log.LG.Errorf("Unexpected duration: %d", event.Duration)
				continue // Skip this event or handle it as needed
			}

			// Calculate end date based on the duration
			endDuration := time.Now().AddDate(0, 0, durationDays)

			// Create the MembershipEvents model
			membershipEvent := model.MembershipEvents{
				UserAddress:     event.User.Hex(),
				OrderID:         event.OrderID,
				TransactionHash: event.TxHash,
				Amount:          event.Amount.String(),
				Status:          1,
				EndDuration:     endDuration,
			}

			// Save the event to the repository
			if err := listener.Repo.CreateMembershipEventHistory(context.Background(), membershipEvent); err != nil {
				log.LG.Errorf("Error saving event to database: %v", err)
			} else {
				log.LG.Infof("Event processed: %+v", membershipEvent)
			}
		}
	}()

	// Keep the listener running until the context is canceled
	<-ctx.Done()
	log.LG.Info("Event listener stopped.")
}

func (listener MembershipEventListener) listen() {
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
	logs, sub, err := subscribeToLogs(listener.ETHClient, common.HexToAddress(listener.ContractAddrress))
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
			eventData, err := listener.parseEventLog(vLog, parsedABI)
			if err != nil {
				log.LG.Errorf("Failed to unpack log: %v", err)
				continue
			}

			// Send the event data to the event channel
			listener.EventChan <- eventData
		}
	}
}

// parseEventLog unpacks the log data into the MembershipEventData structure
func (listener MembershipEventListener) parseEventLog(
	vLog types.Log,
	parsedABI abi.ABI,
) (MembershipEventData, error) {
	event := struct {
		User     common.Address
		Amount   *big.Int
		OrderID  uint64
		Duration uint
	}{}

	// Unpack the log data into the event structure
	err := parsedABI.UnpackIntoInterface(&event, "MembershipPurchased", vLog.Data)
	if err != nil {
		return MembershipEventData{}, fmt.Errorf("failed to unpack log: %w", err)
	}

	// Extract indexed fields (user address and order ID)
	event.User = common.HexToAddress(vLog.Topics[1].Hex())

	// Parse the OrderID from hex to uint64
	orderID, err := parseHexToUint64(vLog.Topics[2].Hex())
	if err != nil {
		return MembershipEventData{}, fmt.Errorf("failed to parse order ID: %w", err)
	}

	// Prepare event data, including the parsed duration
	eventData := MembershipEventData{
		User:     event.User,
		Amount:   event.Amount,
		OrderID:  orderID,
		Duration: event.Duration,
		TxHash:   vLog.TxHash.Hex(),
	}

	return eventData, nil
}
