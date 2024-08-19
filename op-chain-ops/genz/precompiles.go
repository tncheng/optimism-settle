package genz

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var deployConfigAddr = common.Address(crypto.Keccak256([]byte("optimism.deployconfig"))[12:])
var deploymentRegistryAddr = common.Address(crypto.Keccak256([]byte("optimism.deploymentregistry"))[12:])
