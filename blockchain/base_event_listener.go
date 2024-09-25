package blockchain

import (
	"context"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/genefriendway/onchain-handler/internal/interfaces"
	"github.com/genefriendway/onchain-handler/internal/utils/log"
)

const (
	DefaultEventChannelBufferSize = 25   // Buffer size for event channel
	DefaultBlockOffset            = 10   // Default block offset if last processed block is missing
	MaxBlockRange                 = 2048 // Maximum number of blocks to query at once
)

// BaseEventListener represents the shared behavior of any blockchain event listener.
type BaseEventListener struct {
	ETHClient       *ethclient.Client
	ContractAddress common.Address
	EventChan       chan interface{}
	ParsedABI       abi.ABI
	LastBlockRepo   interfaces.BlockStateRepository // Repository to store last processed block
}

// NewBaseEventListener initializes a base listener.
func NewBaseEventListener(
	client *ethclient.Client,
	contractAddr string,
	parsedABI abi.ABI,
	lastBlockRepo interfaces.BlockStateRepository,
) *BaseEventListener {
	// Create a base listener with a buffered event channel.
	eventChan := make(chan interface{}, DefaultEventChannelBufferSize)

	return &BaseEventListener{
		ETHClient:       client,
		ContractAddress: common.HexToAddress(contractAddr),
		EventChan:       eventChan,
		ParsedABI:       parsedABI,
		LastBlockRepo:   lastBlockRepo,
	}
}

// RunListener starts the listener and processes incoming events.
func (listener *BaseEventListener) RunListener(ctx context.Context, parseAndProcessFunc func(types.Log) (interface{}, error)) error {
	// Start listening for events.
	go listener.listen(ctx, parseAndProcessFunc)

	// Handle incoming events from the channel.
	go listener.processEvents(ctx)

	// Wait for the context cancellation to stop the listener.
	<-ctx.Done()
	log.LG.Info("Event listener stopped.")
	return nil
}

// listen polls the blockchain for logs and parses them.
func (listener *BaseEventListener) listen(ctx context.Context, parseAndProcessFunc func(types.Log) (interface{}, error)) {
	log.LG.Info("Starting event listener...")

	// Get the latest block number to query up to
	latestBlock, err := getLatestBlockNumber(ctx, listener.ETHClient)
	if err != nil {
		log.LG.Errorf("Failed to retrieve latest block number from blockchain: %v", err)
		return
	}

	// Get the last processed block number.
	lastBlock, err := listener.LastBlockRepo.GetLastProcessedBlock(ctx)
	if err != nil || lastBlock == 0 {
		log.LG.Warnf("Failed to get last processed block or it was zero: %v", err)
		lastBlock = latestBlock.Uint64() - DefaultBlockOffset // Default to the offset before the latest block
	}

	// Poll for logs in chunks to avoid querying too many blocks at once.
	currentBlock := lastBlock

	for currentBlock < latestBlock.Uint64() {
		// Determine the upper bound for this chunk.
		endBlock := currentBlock + MaxBlockRange
		if endBlock > latestBlock.Uint64() {
			endBlock = latestBlock.Uint64()
		}

		// Poll logs for this block range.
		logs, err := pollForLogsFromBlock(ctx, listener.ETHClient, listener.ContractAddress, currentBlock, endBlock)
		if err != nil {
			log.LG.Errorf("Failed to poll logs from block %d to %d: %v", currentBlock, endBlock, err)
			return
		}

		// Process each log entry.
		for _, vLog := range logs {
			eventData, err := parseAndProcessFunc(vLog)
			if err != nil {
				log.LG.Errorf("Failed to parse log: %v", err)
				continue
			}

			// Send the parsed event data to the event channel.
			select {
			case listener.EventChan <- eventData:
				// Store the latest block number after processing each event.
				if err := listener.LastBlockRepo.UpdateLastProcessedBlock(ctx, vLog.BlockNumber); err != nil {
					log.LG.Errorf("Failed to update last processed block: %v", err)
				}
			case <-ctx.Done():
				return // Exit if the context is canceled.
			}
		}

		// Move to the next block range.
		currentBlock = endBlock + 1
	}
}

// processEvents handles events from the event channel. Override in specific listener if needed.
func (listener *BaseEventListener) processEvents(ctx context.Context) {
	for {
		select {
		case event := <-listener.EventChan:
			log.LG.Infof("Event processed: %+v", event)

		case <-ctx.Done():
			return // Exit if the context is canceled.
		}
	}
}
