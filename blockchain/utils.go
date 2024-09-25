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

// pollForLogsFromBlock polls logs from a specified block number onwards for the given contract.
func pollForLogsFromBlock(
	ctx context.Context,
	client *ethclient.Client, // Ethereum client
	contractAddr common.Address, // Contract address to filter logs
	fromBlock uint64, // Block number to start querying from
	endBlock uint64,
) ([]types.Log, error) {
	// Prepare filter query for the logs
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddr}, // Contract address to filter logs from
		FromBlock: big.NewInt(int64(fromBlock)),   // Start block for querying logs
		ToBlock:   big.NewInt(int64(endBlock)),    // End block for querying logs
	}

	// Poll for logs
	logs, err := client.FilterLogs(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve logs from block %d: %w", fromBlock, err)
	}

	return logs, nil
}

// getLatestBlockNumber retrieves the latest block number from the Ethereum client
func getLatestBlockNumber(ctx context.Context, client *ethclient.Client) (*big.Int, error) {
	header, err := client.HeaderByNumber(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch the latest block header: %w", err)
	}
	return header.Number, nil
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
