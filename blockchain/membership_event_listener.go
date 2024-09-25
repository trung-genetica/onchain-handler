package blockchain

import (
	"context"
	"fmt"
	"math/big"
	"path/filepath"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/genefriendway/onchain-handler/internal/interfaces"
)

// MembershipEventData represents the event data for a MembershipPurchased event.
type MembershipEventData struct {
	User     common.Address
	Amount   *big.Int
	OrderID  uint64
	TxHash   string
	Duration uint8 // Duration as an integer representing the type (0 for 1 year, 1 for 3 years)
}

// MembershipEventListener listens for MembershipPurchased events.
type MembershipEventListener struct {
	*BaseEventListener
	Repo interfaces.MembershipRepository
}

// NewMembershipEventListener initializes the membership event listener.
func NewMembershipEventListener(
	client *ethclient.Client,
	contractAddr string,
	repo interfaces.MembershipRepository,
	lastBlockRepo interfaces.BlockStateRepository,
) (*MembershipEventListener, error) {
	abiFilePath, err := filepath.Abs("./contracts/abis/MembershipPurchase.abi.json")
	if err != nil {
		return nil, fmt.Errorf("failed to get ABI file path: %w", err)
	}

	parsedABI, err := loadABI(abiFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load ABI: %w", err)
	}

	baseListener := NewBaseEventListener(client, contractAddr, parsedABI, lastBlockRepo)
	return &MembershipEventListener{
		BaseEventListener: baseListener,
		Repo:              repo,
	}, nil
}

// parseAndProcessMembershipEvent handles MembershipPurchased event-specific logic.
func (listener *MembershipEventListener) parseAndProcessMembershipEvent(vLog types.Log) (interface{}, error) {
	event := struct {
		User     common.Address
		Amount   *big.Int
		OrderID  uint64
		Duration uint8
	}{}

	// Unpack the log data into the event structure.
	err := listener.ParsedABI.UnpackIntoInterface(&event, "MembershipPurchased", vLog.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to unpack log for TxHash %s: %w", vLog.TxHash.Hex(), err)
	}

	// Extract indexed fields (user address and order ID).
	event.User = common.HexToAddress(vLog.Topics[1].Hex())
	orderID, err := parseHexToUint64(vLog.Topics[2].Hex())
	if err != nil {
		return nil, fmt.Errorf("failed to parse order ID: %w", err)
	}

	// Create event data.
	eventData := &MembershipEventData{
		User:     event.User,
		Amount:   event.Amount,
		OrderID:  orderID,
		Duration: event.Duration,
		TxHash:   vLog.TxHash.Hex(),
	}

	return eventData, nil
}

// RunListener starts the listener with specific event processing logic.
func (listener *MembershipEventListener) RunListener(ctx context.Context) error {
	// Pass the specific event parsing function.
	return listener.BaseEventListener.RunListener(ctx, listener.parseAndProcessMembershipEvent)
}
