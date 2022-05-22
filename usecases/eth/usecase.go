package eth

import (
	"context"
	"soteria-functions/client/gcp"
	"soteria-functions/client/wallet"
	"soteria-functions/entity"
	"soteria-functions/env"
	"soteria-functions/repository/balance/state"
	"soteria-functions/repository/ethaccount"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/pkg/errors"
)

type EthUsecase struct {
	client    *firestore.Client
	kmsClient gcp.KmsClient
	chainID   entity.EthChainID
}

func NewEthUsecase(client *firestore.Client) *EthUsecase {
	kmsClient := gcp.NewKmsClient(env.EnvVars.GCP.ProjectID, env.EnvVars.GCP.KmsLocationID)
	return &EthUsecase{client: client, kmsClient: kmsClient, chainID: env.EnvVars.ETH.ChainID}
}

func (e *EthUsecase) CreateEthAccount(uid string) (*entity.EthAccount, error) {
	wallet := wallet.NewEthWallet(e.chainID)

	privKey, err := wallet.PrivateKey()
	if err != nil {
		return nil, errors.Wrap(err, "failed get private key")
	}

	encPrivKey, err := e.kmsClient.Encryption(entity.EthKeyRing.String(), entity.PrivateKeyName.String(), privKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed encrypt private key")
	}

	address, err := wallet.Address()
	if err != nil {
		return nil, errors.Wrap(err, "failed create eth account")
	}
	return &entity.EthAccount{UserID: uid, EncPrivKey: encPrivKey, Address: address}, nil
}

func (e *EthUsecase) SaveEthAccount(ctx context.Context, account *entity.EthAccount) error {
	f := func(ctx context.Context, tx *firestore.Transaction) error {
		ethRepository := ethaccount.NewTxFirestoreRepository(e.client, tx)
		balanceRepository := state.NewTxFirestoreRepository(e.client, tx)

		err := ethRepository.Save(ctx, account)
		if err != nil {
			return errors.Wrapf(err, "failed ethRepository.Save. %+v", account)
		}

		now := time.Now()
		balance := entity.NewBalance(account.UserID, 0, entity.USDS, now, now)
		err = balanceRepository.Save(ctx, balance)
		if err != nil {
			return errors.Wrapf(err, "failed balanceRepository.Save. %+v", balance)
		}
		return nil
	}

	if err := e.client.RunTransaction(ctx, f); err != nil {
		return errors.Wrap(err, "failed ethaccount.Repository.Save.")
	}
	return nil
}
