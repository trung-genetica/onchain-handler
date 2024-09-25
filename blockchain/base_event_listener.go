package blockchain

import (
	"context"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/genefriendway/onchain-handler/internal/interfaces"
	"github.com/genefriendway/onchain-handler/internal/utils/log"
)

const (
	DefaultEventChannelBufferSize = 100             // Buffer size for event channel
	DefaultBlockOffset            = 10              // Default block offset if last processed block is missing
	MaxBlockRange                 = 2048            // Maximum number of blocks to query at once
	MaxRetries                    = 3               // Maximum number of retries when polling fails
	RetryDelay                    = 2 * time.Second // Delay between retries
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
	var wg sync.WaitGroup
	wg.Add(2) // Two goroutines: listen and processEvents

	go func() {
		defer wg.Done()
		listener.listen(ctx, parseAndProcessFunc)
	}()

	go func() {
		defer wg.Done()
		listener.processEvents(ctx)
	}()

	<-ctx.Done()
	log.LG.Info("Event listener stopped.")

	// Wait for the goroutines to finish
	wg.Wait()

	// Ensure the channel is closed when the listener stops
	close(listener.EventChan)
	return nil
}

// listen polls the blockchain for logs and parses them.
func (listener *BaseEventListener) listen(ctx context.Context, parseAndProcessFunc func(types.Log) (interface{}, error)) {
	log.LG.Info("Starting event listener...")

	latestBlock, err := getLatestBlockNumber(ctx, listener.ETHClient)
	if err != nil {
		log.LG.Errorf("Failed to retrieve latest block number from blockchain: %v", err)
		return
	}

	lastBlock, err := listener.LastBlockRepo.GetLastProcessedBlock(ctx)
	if err != nil || lastBlock == 0 {
		log.LG.Warnf("Failed to get last processed block or it was zero: %v", err)
		if latestBlock.Uint64() > DefaultBlockOffset {
			lastBlock = latestBlock.Uint64() - DefaultBlockOffset
		} else {
			lastBlock = 0
		}
	}

	currentBlock := lastBlock

	for currentBlock < latestBlock.Uint64() {
		endBlock := currentBlock + MaxBlockRange
		if endBlock > latestBlock.Uint64() {
			endBlock = latestBlock.Uint64()
		}

		var logs []types.Log
		for retries := 0; retries < MaxRetries; retries++ {
			logs, err = pollForLogsFromBlock(ctx, listener.ETHClient, listener.ContractAddress, currentBlock, endBlock)
			if err != nil {
				log.LG.Warnf("Failed to poll logs from block %d to %d: %v. Retrying...", currentBlock, endBlock, err)
				time.Sleep(RetryDelay)
				continue
			}
			break
		}
		if err != nil {
			log.LG.Errorf("Max retries reached. Failed to poll logs from block %d to %d: %v", currentBlock, endBlock, err)
			return
		}

		for _, vLog := range logs {
			eventData, err := parseAndProcessFunc(vLog)
			if err != nil {
				log.LG.Errorf("Failed to parse log: %v", err)
				continue
			}

			select {
			case listener.EventChan <- eventData:
				log.LG.Infof("Event successfully sent to channel: %+v", eventData)

				if err := listener.LastBlockRepo.UpdateLastProcessedBlock(ctx, vLog.BlockNumber); err != nil {
					log.LG.Errorf("Failed to update last processed block: %v", err)
				}
			case <-ctx.Done():
				log.LG.Info("Context canceled, stopping log processing.")
				return
			default:
				log.LG.Warnf("Event channel is full, dropping event: %+v", eventData)
			}
		}

		currentBlock = endBlock + 1
	}
}

// processEvents handles events from the event channel.
func (listener *BaseEventListener) processEvents(ctx context.Context) {
	for {
		select {
		case event := <-listener.EventChan:
			log.LG.Infof("Event processed: %+v", event)

		case <-ctx.Done():
			log.LG.Info("Context canceled, stopping event processing.")
			return
		}
	}
}
