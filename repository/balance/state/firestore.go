package state

import (
	"context"
	"soteria-functions/entity"
	"soteria-functions/repository/balance"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/pkg/errors"
)

type repository struct {
	client *firestore.Client
}

func NewFirestoreRepository(client *firestore.Client) Repository {
	return &repository{client: client}
}

type txrepository struct {
	client      *firestore.Client
	transaction *firestore.Transaction
}

func NewTxFirestoreRepository(client *firestore.Client, transaction *firestore.Transaction) Repository {
	return &txrepository{client: client, transaction: transaction}
}

const (
	StateCol = "states"
)

type fbalance struct {
	UserID    string           `firestore:"user_id"`
	Amount    int64            `firestore:"amount"`
	Type      entity.AssetType `firestore:"type"`
	CreatedAt time.Time        `firestore:"created_at"`
	UpdatedAt time.Time        `firestore:"updated_at"`
}

func (r *repository) Save(ctx context.Context, balance *entity.Balance) error {
	account := newFbalance(balance.UserID, balance.Amount, balance.Type, balance.CreatedAt, balance.UpdatedAt)
	ref := r.createRef(balance.UserID, balance.Type)
	_, err := ref.Create(ctx, account)
	if err != nil {
		return errors.Wrap(err, "failed btcaccount save in firestore.")
	}
	return nil
}

func (r *txrepository) Save(_ctx context.Context, balance *entity.Balance) error {
	account := newFbalance(balance.UserID, balance.Amount, balance.Type, balance.CreatedAt, balance.UpdatedAt)
	ref := r.createRef(balance.UserID, balance.Type)
	err := r.transaction.Create(ref, account)
	if err != nil {
		return errors.Wrap(err, "failed btcaccount save in firestore.")
	}
	return nil
}

func (r *repository) AddAmount(ctx context.Context, userID string, assetType entity.AssetType, increment int64) error {
	ref := r.createRef(userID, assetType)
	now := time.Now()
	_, err := ref.Update(ctx, []firestore.Update{
		{Path: "amount", Value: firestore.Increment(increment)},
		{Path: "updated_at", Value: now},
	})
	if err != nil {
		return errors.Wrap(err, "failed AddAmount save in firestore.")
	}
	return nil
}

func (r *txrepository) AddAmount(_ctx context.Context, userID string, assetType entity.AssetType, increment int64) error {
	ref := r.createRef(userID, assetType)
	now := time.Now()
	err := r.transaction.Update(ref, []firestore.Update{
		{Path: "amount", Value: firestore.Increment(increment)},
		{Path: "updated_at", Value: now},
	})
	if err != nil {
		return errors.Wrap(err, "failed AddAmount save in firestore.")
	}
	return nil
}

func (r *repository) Fetch(ctx context.Context, userID string, assetType entity.AssetType) (*entity.Balance, error) {
	ref := r.createRef(userID, assetType)
	doc, err := ref.Get(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed AddAmount save in firestore.")
	}

	var fb fbalance
	doc.DataTo(&fb)
	balance := entity.NewBalance(fb.UserID, fb.Amount, fb.Type, fb.CreatedAt, fb.UpdatedAt)

	if err = balance.Validate(); err != nil {
		return nil, errors.Wrap(err, "invalid balance")
	}
	return balance, nil
}

func (r *txrepository) Fetch(ctx context.Context, userID string, assetType entity.AssetType) (*entity.Balance, error) {
	ref := r.createRef(userID, assetType)
	doc, err := r.transaction.Get(ref)
	if err != nil {
		return nil, errors.Wrap(err, "failed AddAmount save in firestore.")
	}

	var fb fbalance
	doc.DataTo(&fb)
	balance := entity.NewBalance(fb.UserID, fb.Amount, fb.Type, fb.CreatedAt, fb.UpdatedAt)

	if err = balance.Validate(); err != nil {
		return nil, errors.Wrap(err, "invalid balance")
	}
	return balance, nil
}

func newFbalance(userID string, amount int64, assetType entity.AssetType, createdAt, updatedAt time.Time) *fbalance {
	return &fbalance{UserID: userID, Amount: amount, Type: assetType, CreatedAt: createdAt, UpdatedAt: updatedAt}
}

func createRef(client *firestore.Client, uid string, assetType entity.AssetType) *firestore.DocumentRef {
	return client.Collection(balance.BalanceCol).Doc(uid).Collection(StateCol).Doc(string(assetType))
}

func (r *repository) createRef(uid string, assetType entity.AssetType) *firestore.DocumentRef {
	return createRef(r.client, uid, assetType)
}

func (r *txrepository) createRef(uid string, assetType entity.AssetType) *firestore.DocumentRef {
	return createRef(r.client, uid, assetType)
}
