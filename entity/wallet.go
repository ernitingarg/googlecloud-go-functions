package entity

type BtcWallet struct {
	PrivateKey string
	Address    string
}

type BtcChainID string

const (
	MainBtc BtcChainID = "main"
	TestBtc BtcChainID = "test"
)

func (v BtcChainID) String() string {
	return string(v)
}

type EthWallet struct {
	PrivateKey string
	Address    string
}

type EthChainID string

const (
	MainEth EthChainID = "main"
	TestEth EthChainID = "test"
)

func (v EthChainID) String() string {
	return string(v)
}
