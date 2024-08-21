package genz

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/params"

	"github.com/ethereum-optimism/optimism/op-chain-ops/genesis"
	"github.com/ethereum-optimism/optimism/op-chain-ops/genz/devkeys"
)

type InteropDevRecipe struct {
	L1ChainID        uint64
	L2ChainIDs       []uint64
	GenesisTimestamp uint64
}

func (r *InteropDevRecipe) Build(addrs devkeys.DevAddresses) (*WorldConfig, error) {
	// L1 genesis
	l1Cfg := &L1Config{
		ChainID: new(big.Int).SetUint64(r.L1ChainID),
		DevL1DeployConfig: genesis.DevL1DeployConfig{
			L1BlockTime:             6,
			L1GenesisBlockTimestamp: hexutil.Uint64(r.GenesisTimestamp),
			L1GenesisBlockGasLimit:  30_000_000,
		},
	}
	superchainAddrs := devkeys.Scope(addrs, devkeys.SuperchainKeyDomain, l1Cfg.ChainID)

	superchainDeployer, err := superchainAddrs(devkeys.DeployerRole)
	if err != nil {
		return nil, err
	}
	finalSystemOwner, err := superchainAddrs(devkeys.FinalSystemOwnerRole)
	if err != nil {
		return nil, err
	}
	superchainProxyAdmin, err := superchainAddrs(devkeys.ProxyAdminOwnerRole)
	if err != nil {
		return nil, err
	}
	superchainConfigGuardian, err := superchainAddrs(devkeys.SuperchainConfigGuardianRole)
	if err != nil {
		return nil, err
	}
	superchainCfg := &SuperchainConfig{
		FinalSystemOwner: finalSystemOwner,
		ProxyAdminOwner:  superchainProxyAdmin,
		Deployer:         superchainDeployer,
		SuperchainL1DeployConfig: genesis.SuperchainL1DeployConfig{
			RequiredProtocolVersion:    params.OPStackSupport,
			RecommendedProtocolVersion: params.OPStackSupport,
			SuperchainConfigGuardian:   superchainConfigGuardian,
		},
	}
	world := &WorldConfig{
		L1:         l1Cfg,
		Superchain: superchainCfg,
		L2s:        make(map[string]*L2Config),
	}
	for _, l2ChainID := range r.L2ChainIDs {
		l2Cfg, err := InteropL2DevConfig(r.L1ChainID, l2ChainID, addrs)
		if err != nil {
			return nil, fmt.Errorf("failed to generate L2 config for chain %d: %w", l2ChainID, err)
		}
		world.L2s[fmt.Sprintf("%d", l2ChainID)] = l2Cfg
	}
	return world, nil
}

func InteropL2DevConfig(l1ChainID, l2ChainID uint64, addrs devkeys.DevAddresses) (*L2Config, error) {
	// Padded chain ID, hex encoded, prefixed with 0xff like inboxes, then 0x02 to signify devnet.
	batchInboxAddress := common.HexToAddress(fmt.Sprintf("0xff02%016x", l2ChainID))
	operatorAddrs := devkeys.Scope(addrs, devkeys.L2OperatorKeyDomain, new(big.Int).SetUint64(l2ChainID))

	deployer, err := operatorAddrs(devkeys.DeployerRole)
	if err != nil {
		return nil, err
	}
	proxyAdminOwner, err := operatorAddrs(devkeys.ProxyAdminOwnerRole)
	if err != nil {
		return nil, err
	}
	finalSystemOwner, err := operatorAddrs(devkeys.FinalSystemOwnerRole)
	if err != nil {
		return nil, err
	}
	baseFeeVaultRecipient, err := operatorAddrs(devkeys.BaseFeeVaultRecipientRole)
	if err != nil {
		return nil, err
	}
	l1FeeVaultRecipient, err := operatorAddrs(devkeys.L1FeeVaultRecipientRole)
	if err != nil {
		return nil, err
	}
	sequencerFeeVaultRecipient, err := operatorAddrs(devkeys.SequencerFeeVaultRecipientRole)
	if err != nil {
		return nil, err
	}
	sequencerP2P, err := operatorAddrs(devkeys.SequencerP2PRole)
	if err != nil {
		return nil, err
	}
	batcher, err := operatorAddrs(devkeys.BatcherRole)
	if err != nil {
		return nil, err
	}

	return &L2Config{
		Deployer: deployer,
		L2InitializationConfig: genesis.L2InitializationConfig{
			DevDeployConfig: genesis.DevDeployConfig{
				FundDevAccounts: true,
			},
			L2GenesisBlockDeployConfig: genesis.L2GenesisBlockDeployConfig{
				L2GenesisBlockGasLimit:      30_000_000,
				L2GenesisBlockBaseFeePerGas: (*hexutil.Big)(big.NewInt(params.InitialBaseFee)),
			},
			OwnershipDeployConfig: genesis.OwnershipDeployConfig{
				ProxyAdminOwner:  proxyAdminOwner,
				FinalSystemOwner: finalSystemOwner,
			},
			L2VaultsDeployConfig: genesis.L2VaultsDeployConfig{
				BaseFeeVaultRecipient:                    baseFeeVaultRecipient,
				L1FeeVaultRecipient:                      l1FeeVaultRecipient,
				SequencerFeeVaultRecipient:               sequencerFeeVaultRecipient,
				BaseFeeVaultMinimumWithdrawalAmount:      (*hexutil.Big)(Ether(10)),
				L1FeeVaultMinimumWithdrawalAmount:        (*hexutil.Big)(Ether(10)),
				SequencerFeeVaultMinimumWithdrawalAmount: (*hexutil.Big)(Ether(10)),
				BaseFeeVaultWithdrawalNetwork:            "remote",
				L1FeeVaultWithdrawalNetwork:              "remote",
				SequencerFeeVaultWithdrawalNetwork:       "remote",
			},
			GovernanceDeployConfig: genesis.GovernanceDeployConfig{
				EnableGovernance: false,
			},
			GasPriceOracleDeployConfig: genesis.GasPriceOracleDeployConfig{
				GasPriceOracleBaseFeeScalar:     1368,
				GasPriceOracleBlobBaseFeeScalar: 810949,
			},
			GasTokenDeployConfig: genesis.GasTokenDeployConfig{
				UseCustomGasToken: false,
			},
			OperatorDeployConfig: genesis.OperatorDeployConfig{
				P2PSequencerAddress: sequencerP2P,
				BatchSenderAddress:  batcher,
			},
			EIP1559DeployConfig: genesis.EIP1559DeployConfig{
				EIP1559Elasticity:        6,
				EIP1559Denominator:       50,
				EIP1559DenominatorCanyon: 250,
			},
			UpgradeScheduleDeployConfig: genesis.UpgradeScheduleDeployConfig{
				L2GenesisRegolithTimeOffset: new(hexutil.Uint64),
				L2GenesisCanyonTimeOffset:   new(hexutil.Uint64),
				L2GenesisDeltaTimeOffset:    new(hexutil.Uint64),
				L2GenesisEcotoneTimeOffset:  new(hexutil.Uint64),
				L2GenesisFjordTimeOffset:    new(hexutil.Uint64),
				L2GenesisGraniteTimeOffset:  new(hexutil.Uint64),
				L2GenesisInteropTimeOffset:  new(hexutil.Uint64),
				L1CancunTimeOffset:          new(hexutil.Uint64),
				UseInterop:                  true,
			},
			L2CoreDeployConfig: genesis.L2CoreDeployConfig{
				L1ChainID:                 l1ChainID,
				L2ChainID:                 l2ChainID,
				L2BlockTime:               2,
				FinalizationPeriodSeconds: 2, // instant output finalization
				MaxSequencerDrift:         300,
				SequencerWindowSize:       200,
				ChannelTimeoutBedrock:     120,
				ChannelTimeoutGranite:     50, // becoming a constant soon
				BatchInboxAddress:         batchInboxAddress,
				SystemConfigStartBlock:    0,
			},
			AltDADeployConfig: genesis.AltDADeployConfig{
				UseAltDA: false,
			},
		},
		FaultProofDeployConfig: genesis.FaultProofDeployConfig{
			UseFaultProofs:                  true,
			FaultGameAbsolutePrestate:       common.HexToHash("0x03c7ae758795765c6664a5d39bf63841c71ff191e9189522bad8ebff5d4eca98"),
			FaultGameMaxDepth:               50,
			FaultGameClockExtension:         0,
			FaultGameMaxClockDuration:       1200,
			FaultGameGenesisBlock:           0,
			FaultGameGenesisOutputRoot:      common.HexToHash("0xDEADBEEFDEADBEEFDEADBEEFDEADBEEFDEADBEEFDEADBEEFDEADBEEFDEADBEEF"),
			FaultGameSplitDepth:             14,
			FaultGameWithdrawalDelay:        604800,
			PreimageOracleMinProposalSize:   10000,
			PreimageOracleChallengePeriod:   120,
			ProofMaturityDelaySeconds:       12,
			DisputeGameFinalityDelaySeconds: 6,
			RespectedGameType:               254, // "fast" game type
		},
	}, nil
}

var etherScalar = new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)

// Ether converts a uint64 Ether amount into a *big.Int amount in wei units, for allocating test balances.
func Ether(v uint64) *big.Int {
	return new(big.Int).Mul(new(big.Int).SetUint64(v), etherScalar)
}
