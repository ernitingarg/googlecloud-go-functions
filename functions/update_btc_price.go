package functions

import (
	"context"
	"soteria-functions/client/exchange"
	"soteria-functions/env"
	"soteria-functions/infra/firestore"
	"soteria-functions/logger"
	"soteria-functions/usecases/btc"
)

var priceHistoryUsecase *btc.PriceHistoryUsecase

func init() {
	client := firestore.NewFirestoreClient(env.EnvVars.GCP.ProjectID)
	ftxClient := exchange.NewFtxClient()
	priceHistoryUsecase = btc.NewPriceHistoryUsecase(client, ftxClient)
}

const targetMarket = "BTC-PERP"

// CreateEthAccount eth address
func UpdateBtcPrice(ctx context.Context) error {
	err := priceHistoryUsecase.UpdateBtcPrice(ctx, targetMarket)

	if err != nil {
		logger.Log.Error(err.Error())
		return err
	}
	return nil
}
