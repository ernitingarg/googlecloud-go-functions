package btcaccount

import (
	"context"
	"fmt"
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

type txrepository struct {
	client      *firestore.Client
	transaction *firestore.Transaction
}

func NewTxFirestoreRepository(client *firestore.Client, transaction *firestore.Transaction) Repository {
	return &txrepository{client: client, transaction: transaction}
}

const (
	BtcAccountCol = "btc_accounts"
)

type fbtcAccount struct {
	UserID     string `firestore:"user_id"`
	EncPrivKey string `firestore:"enc_priv_key"`
	Address    string `firestore:"address"`
}

func newFbtcAccount(userID, encprivKey, address string) *fbtcAccount {
	return &fbtcAccount{UserID: userID, EncPrivKey: encprivKey, Address: address}
}

func colRef(client *firestore.Client) *firestore.CollectionRef {
	return client.Collection(BtcAccountCol)
}

func docRef(client *firestore.Client, uid string) *firestore.DocumentRef {
	return client.Collection(BtcAccountCol).Doc(uid)
}

func (r *repository) Save(ctx context.Context, btcAccount *entity.BtcAccount) error {
	account := newFbtcAccount(btcAccount.UserID, btcAccount.EncPrivKey, btcAccount.Address)
	ref := docRef(r.client, account.UserID)
	_, err := ref.Create(ctx, account)
	if err != nil {
		return errors.Wrap(err, "failed btcaccount save in firestore.")
	}
	return nil
}

func (r *txrepository) Save(ctx context.Context, btcAccount *entity.BtcAccount) error {
	account := newFbtcAccount(btcAccount.UserID, btcAccount.EncPrivKey, btcAccount.Address)
	ref := docRef(r.client, account.UserID)
	err := r.transaction.Create(ref, account)
	if err != nil {
		return errors.Wrap(err, "failed btcaccount save in firestore.")
	}
	return nil
}

func (r *repository) FetchByAddress(ctx context.Context, address string) (*entity.BtcAccount, error) {
	query := colRef(r.client).Where("address", "==", address).Limit(1)
	docs, err := query.Documents(ctx).GetAll()

	if err != nil {
		return nil, errors.Wrap(err, "failed error in FetchByAddress")
	}
	if len(docs) == 0 {
		return nil, errors.New("entity.BtcAccount is empty")
	}

	doc := docs[0]
	var faccount fbtcAccount
	doc.DataTo(&faccount)
	account := entity.NewBtcAccount(faccount.UserID, faccount.EncPrivKey, faccount.Address)
	if err := account.Validate(); err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("invalid account. %+v", account))
	}

	return account, nil
}

func (r *txrepository) FetchByAddress(ctx context.Context, address string) (*entity.BtcAccount, error) {
	panic("no implement")
}
