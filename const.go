package cfi

const (
	SwapGroupFX SwapGroup = `FX`
)

type SwapGroup string

const (
	CollectiveGroupETF CollectiveGroup = `ETF`
)

type CollectiveGroup string

const (
	EquityGroupCommon    EquityGroup = `Common`
	EquityGroupPreferred EquityGroup = `Preferred`
)

type EquityGroup string
