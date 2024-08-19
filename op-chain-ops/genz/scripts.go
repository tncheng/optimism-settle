package genz

type DeployScript struct {
	DeploySafe func(name string)

	// setupSuperchain
	DeployAddressManager        func()
	DeployProxyAdmin            func()
	TransferProxyAdminOwnership func()
	DeploySuperchainConfig      func()
	InitializeSuperchainConfig  func()
	DeployERC1967Proxy          func(name string)
	DeployProtocolVersions      func()
	InitializeProtocolVersions  func()

	// technically still part of setupOpChain, but we deploy them once, and share them, OPSM style.
	DeployImplementations func()

	// setupOpChain
	DeployProxies             func()
	InitializeImplementations func()

	SetAlphabetFaultGameImplementation           func(allowUpgrade bool)
	SetFastFaultGameImplementation               func(allowUpgrade bool)
	SetCannonFaultGameImplementation             func(allowUpgrade bool)
	SetPermissionedCannonFaultGameImplementation func(allowUpgrade bool)
	TransferDisputeGameFactoryOwnership          func()
	TransferDelayedWETHOwnership                 func()
}

type L2GenesisScript struct {
	RunWithAllUpgrades func() error
	SetPreinstalls     func() error
}
