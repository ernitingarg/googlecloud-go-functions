package wallet

import (
	"log"
	"soteria-functions/entity"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/pkg/errors"
)

type btcWallet struct {
	key *hdkeychain.ExtendedKey
	net *chaincfg.Params
}

func NewBtcWallet(chainID entity.BtcChainID) Wallet {
	var net *chaincfg.Params
	switch chainID {
	case entity.MainBtc:
		net = &chaincfg.MainNetParams
	case entity.TestBtc:
		net = &chaincfg.TestNet3Params
	default:
		net = &chaincfg.MainNetParams
	}

	seed, err := hdkeychain.GenerateSeed(hdkeychain.RecommendedSeedLen)
	if err != nil {
		log.Fatalln(errors.Wrap(err, "failed generate seed"))
	}

	key, err := hdkeychain.NewMaster(seed, net)
	if err != nil {
		log.Fatalln(errors.Wrap(err, "failed new master btcWallet"))
	}

	return &btcWallet{key: key, net: net}
}

func (w *btcWallet) PrivateKey() (string, error) {
	return w.key.String(), nil
}

func (w *btcWallet) Address() (string, error) {
	address, err := w.key.Address(w.net)
	if err != nil {
		return "", errors.Wrap(err, "failed btcWallet.Address")
	}
	return address.EncodeAddress(), nil
}
