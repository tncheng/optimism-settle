package genz

type DeployScript struct {
	DeploySafe func(name string) error

	// setupSuperchain
	SetupSuperchain func() error

	// technically still part of setupOpChain, but we deploy them once, and share them, OPSM style.
	DeployImplementations func() error

	// setupOpChain
	DeployProxies             func() error
	InitializeImplementations func() error

	// FP functions
	SetAlphabetFaultGameImplementation           func(allowUpgrade bool) error
	SetFastFaultGameImplementation               func(allowUpgrade bool) error
	SetCannonFaultGameImplementation             func(allowUpgrade bool) error
	SetPermissionedCannonFaultGameImplementation func(allowUpgrade bool) error
	TransferDisputeGameFactoryOwnership          func() error
	TransferDelayedWETHOwnership                 func() error
}

type L2GenesisScript struct {
	RunWithAllUpgrades func() error
	SetPreinstalls     func() error
}
