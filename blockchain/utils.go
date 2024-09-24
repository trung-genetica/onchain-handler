package blockchain

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"strconv"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// loadABI loads and parses the ABI from a JSON file
func loadABI(path string) (abi.ABI, error) {
	file, err := os.Open(path)
	if err != nil {
		return abi.ABI{}, fmt.Errorf("failed to open ABI file at %s: %w", path, err)
	}
	defer file.Close()

	// Read and parse the ABI JSON file
	var parsedABI abi.ABI
	if err := json.NewDecoder(file).Decode(&parsedABI); err != nil {
		return abi.ABI{}, fmt.Errorf("failed to decode ABI JSON from file %s: %w", path, err)
	}

	return parsedABI, nil
}

// pollForLogs polls logs for the contract from a given block offset
func pollForLogs(
	client *ethclient.Client, // Ethereum client
	contractAddr common.Address, // Contract address to filter logs
	blockOffset int64, // How many blocks in the past to start querying
) ([]types.Log, error) {
	// Get the latest block number
	latestBlock, err := getLatestBlockNumber(client)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve latest block number: %w", err)
	}

	// Define the fromBlock based on the offset
	fromBlock := calculateFromBlock(latestBlock, blockOffset)

	// Prepare filter query for the logs
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddr},
		FromBlock: fromBlock,
	}

	// Poll for logs
	logs, err := client.FilterLogs(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve contract logs: %w", err)
	}

	return logs, nil
}

// getLatestBlockNumber retrieves the latest block number from the Ethereum client
func getLatestBlockNumber(client *ethclient.Client) (*big.Int, error) {
	header, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch the latest block header: %w", err)
	}
	return header.Number, nil
}

// calculateFromBlock calculates the starting block for logs, ensuring it doesn't go below block 0
func calculateFromBlock(latestBlock *big.Int, blockOffset int64) *big.Int {
	fromBlock := new(big.Int).Sub(latestBlock, big.NewInt(blockOffset))
	if fromBlock.Sign() < 0 {
		return big.NewInt(0) // Ensure fromBlock doesn't go below 0
	}
	return fromBlock
}

// parseHexToUint64 converts a hex string to uint64
func parseHexToUint64(hexStr string) (uint64, error) {
	// Strip the "0x" prefix if present
	if len(hexStr) > 2 && hexStr[:2] == "0x" {
		hexStr = hexStr[2:]
	}
	// Convert hex string to uint64
	orderID, err := strconv.ParseUint(hexStr, 16, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse hex string %s to uint64: %w", hexStr, err)
	}
	return orderID, nil
}
