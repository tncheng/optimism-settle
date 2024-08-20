package genz

import (
	"math/big"

	"github.com/holiman/uint256"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ethereum-optimism/optimism/op-chain-ops/genesis"
)

type L1Config struct {
	ChainID  *big.Int
	Deployer common.Address
	genesis.DevL1DeployConfig
}

type SuperchainConfig struct {
	// TODO address manager owner, proxy admin owner?

	ProxyAdminOwner common.Address
	Deployer        common.Address

	genesis.SuperchainL1DeployConfig
}

type L2Config struct {
	Deployer common.Address // account used to deploy contracts to L2
	genesis.L2InitializationConfig
	genesis.OutputOracleDeployConfig
	genesis.FaultProofDeployConfig
}

type WorldConfig struct {
	L1         *L1Config
	Superchain *SuperchainConfig
	L2s        map[uint256.Int]*L2Config
}
