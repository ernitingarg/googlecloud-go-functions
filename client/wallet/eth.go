package wallet

import (
	"log"
	"soteria-functions/entity"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/ethereum/go-ethereum/accounts"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
	"github.com/pkg/errors"
)

type ethWallet struct {
	wallet  *hdwallet.Wallet
	account *accounts.Account
	net     *chaincfg.Params
}

func NewEthWallet(param entity.EthChainID) Wallet {
	seed, err := hdwallet.NewSeed()
	if err != nil {
		log.Fatalln(errors.Wrap(err, "failed generate seed"))
	}

	wallet, err := hdwallet.NewFromSeed(seed)
	if err != nil {
		log.Fatalln(errors.Wrap(err, "failed new master ethWallet"))
	}

	path := hdwallet.MustParseDerivationPath("m/44'/60'/0'/0/0")
	account, err := wallet.Derive(path, false)
	if err != nil {
		log.Fatalln(errors.Wrap(err, "failed empty accounts"))
	}

	var net *chaincfg.Params
	switch param {
	case entity.MainEth:
		net = &chaincfg.MainNetParams
	case entity.TestEth:
		net = &chaincfg.TestNet3Params
	default:
		net = &chaincfg.MainNetParams
	}

	return &ethWallet{wallet: wallet, account: &account, net: net}
}

func (w *ethWallet) PrivateKey() (string, error) {
	return w.wallet.PrivateKeyHex(*w.account)
}

func (w *ethWallet) Address() (string, error) {
	return w.account.Address.Hex(), nil
}
