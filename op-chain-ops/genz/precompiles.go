package genz

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/ethereum-optimism/optimism/op-chain-ops/script"
)

var deployConfigAddr = common.Address(crypto.Keccak256([]byte("optimism.deployconfig"))[12:])
var deploymentRegistryAddr = common.Address(crypto.Keccak256([]byte("optimism.deploymentregistry"))[12:])

func WithPrecompileAtAddress[E any](h *script.Host, addr common.Address, elem E) (cleanup func(), err error) {
	precompile, err := script.NewPrecompile[E](elem)
	if err != nil {
		return nil, fmt.Errorf("failed to construct precompile: %w", err)
	}
	_ = precompile
	// set code to []byte{0}
	// override precompile
	return // TODO
}

func WithPrecompile(h *script.Host, precompile any) (addr common.Address, cleanup func()) {
	// create tmp addr
	// set code to []byte{0}
	// override precompile
	return // TODO
}
