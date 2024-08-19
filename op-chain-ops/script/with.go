package script

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

func checkABI(abiData *abi.ABI, methodSignature string) bool {
	for _, m := range abiData.Methods {
		if m.Sig == methodSignature {
			return true
		}
	}
	return false
}

func WithScript[B any](h *Host, name string, contract string) (b *B, cleanup func(), err error) {
	// load contract artifact
	artifact, err := h.af.ReadArtifact(name, contract)
	if err != nil {
		return nil, nil, fmt.Errorf("could not load script artifact: %w", err)
	}

	// TODO compute address of script contract to be deployed
	addr := common.Address{}

	// init bindings (with ABI check)
	bindings, err := MakeBindings[B](h.ScriptBackendFn(addr), func(abiDef string) bool {
		return checkABI(&artifact.ABI, abiDef)
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to make bindings: %w", err)
	}

	// TODO deploy the script contract

	// TODO cleanup func to remove script contract
	return bindings, nil, errors.New("TODO")
}

// WithPrecompileAtAddress turns a struct into a precompile,
// and inserts it as override at the given address in the host.
// A cleanup function is returned, to remove the precompile override again.
func WithPrecompileAtAddress[E any](h *Host, addr common.Address, elem E) (cleanup func(), err error) {
	if h.HasPrecompileOverride(addr) {
		return nil, fmt.Errorf("already have existing precompile override at %s", addr)
	}
	precompile, err := NewPrecompile[E](elem)
	if err != nil {
		return nil, fmt.Errorf("failed to construct precompile: %w", err)
	}
	h.SetPrecompile(addr, precompile)
	return func() {
		h.SetPrecompile(addr, nil)
	}, nil
}
