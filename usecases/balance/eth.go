package balance

import (
	"context"
	"fmt"
	"soteria-functions/entity"
	"soteria-functions/repository/balance/state"
	"soteria-functions/repository/balance/transaction"
	"soteria-functions/repository/ethaccount"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/pkg/errors"
)

type EthBalanceUsecase struct {
	client     *firestore.Client
	repository ethaccount.Repository
}

func NewEthBalanceUsecase(client *firestore.Client) *EthBalanceUsecase {
	repository := ethaccount.NewFirestoreRepository(client)
	return &EthBalanceUsecase{client: client, repository: repository}
}

func (b *EthBalanceUsecase) AddUsdsBalance(ctx context.Context, address, txid string, amount int64) error {
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

		err := balanceRepository.AddAmount(ctx, account.UserID, entity.USDS, amount)
		if err != nil {
			return errors.Wrapf(err, "failed balanceRepository.AddAmount. %+v, %+v", account.UserID, amount)
		}

		transaction := entity.NewBalanceTransaction(account.UserID, txid, amount, entity.USDS, now, now)
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
