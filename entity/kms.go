package entity

type KeyRingName string

func (v KeyRingName) String() string {
	return string(v)
}

const (
	BtcKeyRing KeyRingName = "btc"
	EthKeyRing KeyRingName = "eth"
)

type KeyName string

func (v KeyName) String() string {
	return string(v)
}

const (
	PrivateKeyName KeyName = "for_private_key"
)
