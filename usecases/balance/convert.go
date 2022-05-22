package balance

import (
	"context"
	"fmt"
	"soteria-functions/client/exchange"
	"soteria-functions/constant"
	"soteria-functions/entity"
	"soteria-functions/repository/balance/state"
	"soteria-functions/repository/balance/transaction"
	"soteria-functions/util"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/pkg/errors"
)

type ConvertUsecase struct {
	client     *firestore.Client
	okexClient *exchange.OkexClient
}

func NewConvertUsecase(client *firestore.Client, okexClient *exchange.OkexClient) *ConvertUsecase {
	return &ConvertUsecase{client: client, okexClient: okexClient}
}

func (c *ConvertUsecase) Convert(ctx context.Context, uid string, fromAmount int64, toAmount int64, fromType entity.AssetType, toType entity.AssetType) error {

	f := func(ctx context.Context, tx *firestore.Transaction) error {
		balanceRepository := state.NewTxFirestoreRepository(c.client, tx)
		transactionRepository := transaction.NewTxFirestoreRepository(c.client, tx)

		fromBalance, err := balanceRepository.Fetch(ctx, uid, fromType)
		if err != nil {
			return errors.Wrapf(err, "failed balanceRepository.Fetch for fromBalance. %+v, %+v", uid, fromType)
		}

		// Validation for fromBalance and toBalance
		if fromBalance.Amount < fromAmount {
			return errors.New(fmt.Sprintf("not enough amount. type: %+v. balance: %+v, require amount: %+v", fromType, fromBalance.Amount, fromAmount))
		}

		err = balanceRepository.AddAmount(ctx, uid, fromType, -fromAmount)
		if err != nil {
			return errors.Wrapf(err, "failed balanceRepository.AddAmount for fromAsset. %+v, %+v, %+v", uid, fromType, fromAmount)
		}

		err = balanceRepository.AddAmount(ctx, uid, toType, toAmount)
		if err != nil {
			return errors.Wrapf(err, "failed balanceRepository.AddAmount for toAsset. %+v, %+v, %+v", uid, toType, toAmount)
		}

		txid, err := util.MakeRandomStr(32)
		if err != nil {
			return errors.Wrapf(err, "failed MakeRandomStr to create txid. %+v, %+v, %+v", uid, fromType, toType)
		}

		now := time.Now()
		fromTransaction := entity.NewBalanceTransaction(uid, txid, fromAmount, fromType, now, now)
		err = transactionRepository.Save(ctx, fromTransaction)
		if err != nil {
			return errors.Wrapf(err, "failed transactionRepository.Save for fromAsset. %+v", fromTransaction)
		}

		txid, err = util.MakeRandomStr(32)
		if err != nil {
			return errors.Wrapf(err, "failed MakeRandomStr to create txid. %+v, %+v, %+v", uid, fromType, toType)
		}
		toTransaction := entity.NewBalanceTransaction(uid, txid, toAmount, toType, now, now)
		err = transactionRepository.Save(ctx, toTransaction)
		if err != nil {
			return errors.Wrapf(err, "failed transactionRepository.Save for toAsset. %+v", fromTransaction)
		}

		return nil
	}

	if err := c.client.RunTransaction(ctx, f); err != nil {
		return errors.Wrap(err, "failed ConvertToken")
	}
	return nil
}

const BUFFER_PRICE = 100.0

func (c *ConvertUsecase) ValidatePrice(ctx context.Context, fromAmount int64, toAmount int64, fromType entity.AssetType, toType entity.AssetType) error {
	if fromType == toType {
		return errors.New("fromType and toType is same")
	}

	okexPrice, err := c.okexClient.GetPrice(ctx, exchange.BTC_USD)
	if err != nil {
		return errors.Wrap(err, "failed getPrice from okexClient")
	}

	fFromAmount := float64(fromAmount)
	fToAmount := float64(toAmount)
	if fromType == entity.USDS && toType == entity.BTC {
		price := fFromAmount / (fToAmount / constant.BTC_BID)
		if price <= (okexPrice - BUFFER_PRICE) {
			return errors.New(fmt.Sprintf("target price <= okexprice. target price: %f. okex price: %f", price, okexPrice))
		}
	} else if fromType == entity.BTC && toType == entity.USDS {
		price := fToAmount / (fFromAmount / constant.BTC_BID)
		if (price - BUFFER_PRICE) >= okexPrice {
			return errors.New(fmt.Sprintf("target price >= okexprice. target price: %f. okex price: %f", price, okexPrice))
		}
	} else {
		return errors.New(fmt.Sprintf("invalid types. fromType: %s toType: %s", fromType, toType))
	}

	return nil
}
