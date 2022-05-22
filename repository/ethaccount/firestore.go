package ethaccount

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
	EthAccountCol = "eth_accounts"
)

type fethAccount struct {
	UserID     string `firestore:"user_id"`
	EncPrivKey string `firestore:"enc_priv_key"`
	Address    string `firestore:"address"`
}

func newFethAccount(userID, encprivKey, address string) *fethAccount {
	return &fethAccount{UserID: userID, EncPrivKey: encprivKey, Address: address}
}

func colRef(client *firestore.Client) *firestore.CollectionRef {
	return client.Collection(EthAccountCol)
}

func docRef(client *firestore.Client, uid string) *firestore.DocumentRef {
	return client.Collection(EthAccountCol).Doc(uid)
}

func (r *repository) Save(ctx context.Context, ethAccount *entity.EthAccount) error {
	account := newFethAccount(ethAccount.UserID, ethAccount.EncPrivKey, ethAccount.Address)
	ref := docRef(r.client, account.UserID)
	_, err := ref.Create(ctx, account)
	if err != nil {
		return errors.Wrap(err, "failed ethaccount save in firestore.")
	}
	return nil
}

func (r *txrepository) Save(ctx context.Context, ethAccount *entity.EthAccount) error {
	account := newFethAccount(ethAccount.UserID, ethAccount.EncPrivKey, ethAccount.Address)
	ref := docRef(r.client, account.UserID)
	err := r.transaction.Create(ref, account)
	if err != nil {
		return errors.Wrap(err, "failed ethaccount save in firestore.")
	}
	return nil
}

func (r *repository) FetchByAddress(ctx context.Context, address string) (*entity.EthAccount, error) {
	query := colRef(r.client).Where("address", "==", address).Limit(1)
	docs, err := query.Documents(ctx).GetAll()

	if err != nil {
		return nil, errors.Wrap(err, "failed error in FetchByAddress")
	}
	if len(docs) == 0 {
		return nil, errors.New("entity.EthAccount is empty")
	}

	doc := docs[0]
	var faccount fethAccount
	doc.DataTo(&faccount)
	account := entity.NewEthAccount(faccount.UserID, faccount.EncPrivKey, faccount.Address)
	if err := account.Validate(); err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("invalid account. %+v", account))
	}

	return account, nil
}

func (r *txrepository) FetchByAddress(ctx context.Context, address string) (*entity.EthAccount, error) {
	panic("no implement")
}
