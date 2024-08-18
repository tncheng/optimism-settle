package genz

import (
	"github.com/holiman/uint256"

	"github.com/ethereum/go-ethereum/common"
)

type L1Deployment struct {
	// other preinstalls maybe?
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
	L1CrossDomainMessengerProxy       common.Address
	L1ERC721BridgeProxy               common.Address
	L1StandardBridgeProxy             common.Address
	L2OutputOracleProxy               common.Address
	OptimismMintableERC20FactoryProxy common.Address
	OptimismPortalProxy               common.Address
	SystemConfigProxy                 common.Address

	// Fault proofs; some of these don't have to be deployed per chain
	AnchorStateRegistryProxy common.Address
	DelayedWETHProxy         common.Address
	DisputeGameFactoryProxy  common.Address
}

type L2Deployment struct {
	L2Proxies

	ProxyAdmin common.Address

	// Safe that will own the L2 chain contracts
	SystemOwnerSafe common.Address `json:"SystemOwnerSafe"`
}

type WorldDeployment struct {
	L1         *L1Deployment
	Superchain *SuperchainDeployment
	L2s        map[uint256.Int]*L2Deployment
}
