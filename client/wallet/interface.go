package wallet

type Wallet interface {
	PrivateKey() (string, error)
	Address() (string, error)
}
