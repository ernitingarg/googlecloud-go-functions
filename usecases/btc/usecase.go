package btc

import (
	"context"
	"soteria-functions/client/gcp"
	"soteria-functions/client/wallet"
	"soteria-functions/entity"
	"soteria-functions/env"
	"soteria-functions/repository/balance/state"
	"soteria-functions/repository/btcaccount"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/pkg/errors"
)

type BtcUsecase struct {
	client    *firestore.Client
	kmsClient gcp.KmsClient
	chainID   entity.BtcChainID
}

func NewBtcUsecase(client *firestore.Client) *BtcUsecase {
	kmsClient := gcp.NewKmsClient(env.EnvVars.GCP.ProjectID, env.EnvVars.GCP.KmsLocationID)
	return &BtcUsecase{client: client, kmsClient: kmsClient, chainID: env.EnvVars.BTC.ChainID}
}

func (b *BtcUsecase) CreateBtcAccount(uid string) (*entity.BtcAccount, error) {
	wallet := wallet.NewBtcWallet(b.chainID)

	privKey, err := wallet.PrivateKey()
	if err != nil {
		return nil, errors.Wrap(err, "failed get private key")
	}
	encPrivKey, err := b.kmsClient.Encryption(entity.BtcKeyRing.String(), entity.PrivateKeyName.String(), privKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed encrypt private key")
	}

	address, err := wallet.Address()
	if err != nil {
		return nil, errors.Wrap(err, "failed create btc account")
	}
	return &entity.BtcAccount{UserID: uid, EncPrivKey: encPrivKey, Address: address}, nil
}

func (b *BtcUsecase) SaveBtcAccount(ctx context.Context, account *entity.BtcAccount) error {
	f := func(ctx context.Context, tx *firestore.Transaction) error {
		btcRepository := btcaccount.NewTxFirestoreRepository(b.client, tx)
		balanceRepository := state.NewTxFirestoreRepository(b.client, tx)

		err := btcRepository.Save(ctx, account)
		if err != nil {
			return errors.Wrapf(err, "failed btcRepository.Save. %+v", account)
		}

		now := time.Now()
		balance := entity.NewBalance(account.UserID, 0, entity.BTC, now, now)
		err = balanceRepository.Save(ctx, balance)
		if err != nil {
			return errors.Wrapf(err, "failed balanceRepository.Save. %+v", balance)
		}
		return nil
	}

	if err := b.client.RunTransaction(ctx, f); err != nil {
		return errors.Wrap(err, "failed btcaccount.Repository.Save.")
	}
	return nil
}
