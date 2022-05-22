package transaction

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
	TransactionCol = "transactions"
)

type ftransaction struct {
	UserID    string           `firestore:"user_id"`
	TxID      string           `firestore:"tx_id"`
	Value     int64            `firestore:"value"`
	Type      entity.AssetType `firestore:"type"`
	CreatedAt time.Time        `firestore:"created_at"`
	UpdatedAt time.Time        `firestore:"updated_at"`
}

func (r *repository) Save(ctx context.Context, btx *entity.BalanceTransaction) error {
	account := newFbalanceTransaction(btx.UserID, btx.TxID, btx.Value, btx.Type, btx.CreatedAt, btx.UpdatedAt)
	ref := r.createRef(btx.UserID, btx.TxID)
	_, err := ref.Create(ctx, account)
	if err != nil {
		return errors.Wrap(err, "failed balance_transactions save in firestore.")
	}
	return nil
}

func (r *txrepository) Save(_ctx context.Context, btx *entity.BalanceTransaction) error {
	account := newFbalanceTransaction(btx.UserID, btx.TxID, btx.Value, btx.Type, btx.CreatedAt, btx.UpdatedAt)
	ref := r.createRef(btx.UserID, btx.TxID)
	err := r.transaction.Create(ref, account)
	if err != nil {
		return errors.Wrap(err, "failed balance_transaction save in firestore.")
	}
	return nil
}

func newFbalanceTransaction(userID string, txid string, value int64, assetType entity.AssetType, createdAt, updatedAt time.Time) *ftransaction {
	return &ftransaction{UserID: userID, TxID: txid, Value: value, Type: assetType, CreatedAt: createdAt, UpdatedAt: updatedAt}
}

func createRef(client *firestore.Client, uid string, txid string) *firestore.DocumentRef {
	return client.Collection(balance.BalanceCol).Doc(uid).Collection(TransactionCol).Doc(string(txid))
}

func (r *repository) createRef(uid string, txid string) *firestore.DocumentRef {
	return createRef(r.client, uid, txid)
}

func (r *txrepository) createRef(uid string, txid string) *firestore.DocumentRef {
	return createRef(r.client, uid, txid)
}
