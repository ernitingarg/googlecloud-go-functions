package btc

import (
	"context"
	"soteria-functions/client/exchange"
	"soteria-functions/entity"
	"soteria-functions/repository/price_history"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/pkg/errors"
)

type PriceHistoryUsecase struct {
	client     *firestore.Client
	ftxClient  *exchange.FtxClient
	repository price_history.Repository
}

func NewPriceHistoryUsecase(client *firestore.Client, ftxClient *exchange.FtxClient) *PriceHistoryUsecase {
	repository := price_history.NewFirestoreRepository(client)
	return &PriceHistoryUsecase{client: client, ftxClient: ftxClient, repository: repository}
}

func (phu *PriceHistoryUsecase) UpdateBtcPrice(ctx context.Context, targetMarket string) error {

	price, err := phu.ftxClient.GetPrice(ctx, targetMarket)
	if err != nil {
		return errors.Wrap(err, "PriceHistoryUsecase: failed to GetPrice in UpdateBtcPrice")
	}

	now := time.Now().Unix()
	priceHistory := entity.NewPriceHistory("BTC-USD", price, "FTX", uint32(now))
	err = phu.repository.Save(ctx, priceHistory)

	if err != nil {
		return errors.Wrapf(err, "PriceHistoryUsecase: failed to save PriceHistory in firestore. %+v", priceHistory)
	}

	return nil
}
