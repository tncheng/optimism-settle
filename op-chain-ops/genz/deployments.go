package genz

import (
	"github.com/holiman/uint256"

	"github.com/ethereum/go-ethereum/common"
)

type L1Deployment struct {
	// preinstalls maybe?
}

type Implementations struct {
	L1CrossDomainMessenger       common.Address `json:"L1CrossDomainMessenger"`
	L1ERC721Bridge               common.Address `json:"L1ERC721Bridge"`
	L1StandardBridge             common.Address `json:"L1StandardBridge"`
	L2OutputOracle               common.Address `json:"L2OutputOracle"`
	OptimismMintableERC20Factory common.Address `json:"OptimismMintableERC20Factory"`
	OptimismPortal               common.Address `json:"OptimismPortal"`
	SystemConfig                 common.Address `json:"SystemConfig"`

	DisputeGameFactory common.Address `json:"DisputeGameFactory"`
}

type SuperchainDeployment struct {
	Implementations

	// Safe that will own the Superchain contracts
	SystemOwnerSafe common.Address `json:"SystemOwnerSafe"`

	AddressManager common.Address `json:"AddressManager"`
	ProxyAdmin     common.Address `json:"ProxyAdmin"`

	ProtocolVersions      common.Address `json:"ProtocolVersions"`
	ProtocolVersionsProxy common.Address `json:"ProtocolVersionsProxy"`

	SuperchainConfig      common.Address `json:"SuperchainConfig"`
	SuperchainConfigProxy common.Address `json:"SuperchainConfigProxy"`
}

type L2Proxies struct {
	L1CrossDomainMessengerProxy       common.Address `json:"L1CrossDomainMessengerProxy"`
	L1ERC721BridgeProxy               common.Address `json:"L1ERC721BridgeProxy"`
	L1StandardBridgeProxy             common.Address `json:"L1StandardBridgeProxy"`
	L2OutputOracleProxy               common.Address `json:"L2OutputOracleProxy"`
	OptimismMintableERC20FactoryProxy common.Address `json:"OptimismMintableERC20FactoryProxy"`
	OptimismPortalProxy               common.Address `json:"OptimismPortalProxy"`
	SystemConfigProxy                 common.Address `json:"SystemConfigProxy"`

	// Fault proofs; some of these don't have to be deployed per chain
	AnchorStateRegistryProxy common.Address `json:"AnchorStateRegistryProxy"`
	DelayedWETHProxy         common.Address `json:"DelayedWETHProxy"`
	DisputeGameFactoryProxy  common.Address `json:"DisputeGameFactoryProxy"`
}

type L2Deployment struct {
	L2Proxies

	ProxyAdmin common.Address `json:"ProxyAdmin"`

	// Safe that will own the L2 chain contracts
	SystemOwnerSafe common.Address `json:"SystemOwnerSafe"`
}

type WorldDeployment struct {
	L1         *L1Deployment
	Superchain *SuperchainDeployment
	L2s        map[uint256.Int]*L2Deployment
}
