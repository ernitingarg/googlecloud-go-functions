package price_history

import (
	"context"
	"soteria-functions/entity"

	"cloud.google.com/go/firestore"
	"github.com/pkg/errors"
)

type repository struct {
	client *firestore.Client
}

func NewFirestoreRepository(client *firestore.Client) Repository {
	return &repository{client: client}
}

const (
	PriceHistoryCol = "price_histories"
)

type fPriceHistory struct {
	CurrencyPair string  `firestore:"currency_pair"`
	Rate         float64 `firestore:"rate"`
	Source       string  `firestore:"source"`
	Timestamp    uint32  `firestore:"timestamp"`
}

func newFPriceHistory(p *entity.PriceHistory) *fPriceHistory {
	return &fPriceHistory{CurrencyPair: p.CurrencyPair, Rate: p.Rate, Source: p.Source, Timestamp: p.Timestamp}
}

func colRef(client *firestore.Client) *firestore.CollectionRef {
	return client.Collection(PriceHistoryCol)
}

func newDocRef(client *firestore.Client) *firestore.DocumentRef {
	return client.Collection(PriceHistoryCol).NewDoc()
}

func (r *repository) Save(ctx context.Context, ph *entity.PriceHistory) error {
	fph := newFPriceHistory(ph)
	ref := newDocRef(r.client)
	_, err := ref.Create(ctx, fph)
	if err != nil {
		return errors.Wrapf(err, "failed to save price history in firestore. data: %+v", ph)
	}
	return nil
}
