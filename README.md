# onchain-handler
## Install golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
## abigen
abigen --abi=./contracts/abis/LifePointToken.abi.json --pkg=lptoken --out=./contracts/lptoken/LPToken.go