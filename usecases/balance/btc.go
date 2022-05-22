package balance

import (
	"context"
	"fmt"
	"soteria-functions/entity"
	"soteria-functions/repository/balance/state"
	"soteria-functions/repository/balance/transaction"
	"soteria-functions/repository/btcaccount"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/pkg/errors"
)

type BtcBalanceUsecase struct {
	client     *firestore.Client
	repository btcaccount.Repository
}

func NewBtcBalanceUsecase(client *firestore.Client) *BtcBalanceUsecase {
	repository := btcaccount.NewFirestoreRepository(client)
	return &BtcBalanceUsecase{client: client, repository: repository}
}

func (b *BtcBalanceUsecase) AddBalance(ctx context.Context, address, txid string, amount int64) error {
	if amount == 0 {
		return errors.New(fmt.Sprintf("ammount is zero. address: %+v", address))
	}
	account, err := b.repository.FetchByAddress(ctx, address)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed fetch by address. address: %+v", address))
	}

	f := func(ctx context.Context, tx *firestore.Transaction) error {
		balanceRepository := state.NewTxFirestoreRepository(b.client, tx)
		transactionRepository := transaction.NewTxFirestoreRepository(b.client, tx)
		now := time.Now()

		err := balanceRepository.AddAmount(ctx, account.UserID, entity.BTC, amount)
		if err != nil {
			return errors.Wrapf(err, "failed balanceRepository.AddAmount. %+v, %+v, %+v", account.UserID, entity.BTC, amount)
		}

		transaction := entity.NewBalanceTransaction(account.UserID, txid, amount, entity.BTC, now, now)
		err = transactionRepository.Save(ctx, transaction)
		if err != nil {
			return errors.Wrapf(err, "failed transactionRepository.Save. %+v", transaction)
		}

		return nil
	}

	if err := b.client.RunTransaction(ctx, f); err != nil {
		return errors.Wrap(err, "failed AddBalance")
	}
	return nil
}
