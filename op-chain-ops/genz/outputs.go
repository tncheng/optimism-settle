package genz

import (
	"github.com/holiman/uint256"

	"github.com/ethereum/go-ethereum/core"

	"github.com/ethereum-optimism/optimism/op-node/rollup"
)

type L1Output struct {
	Genesis *core.Genesis
}

type L2Output struct {
	Genesis   *core.Genesis
	RollupCfg *rollup.Config
}

type WorldOutput struct {
	L1  *L1Output
	L2s map[uint256.Int]*L2Output
}
