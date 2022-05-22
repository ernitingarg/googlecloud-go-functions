package entity

import (
	"time"

	"github.com/pkg/errors"
)

type AssetType string

const (
	BTC  AssetType = "btc"
	USDS AssetType = "usds"
)

func (t *AssetType) Validate() error {
	if *t != BTC && *t != USDS {
		return errors.New("invalid asset type")
	}
	return nil
}

// balances/{uid}/states/{btc|usds}
type Balance struct {
	UserID    string    `json:"user_id"`
	Amount    int64     `json:"amount"`
	Type      AssetType `json:"type"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewBalance(uid string, amount int64, assetType AssetType, createdAt time.Time, updatedAt time.Time) *Balance {
	return &Balance{UserID: uid, Amount: amount, Type: assetType, CreatedAt: createdAt, UpdatedAt: updatedAt}
}

func (b *Balance) Validate() error {
	if b.UserID == "" {
		return errors.New("balance.UserID is empty")
	}
	if b.Type == "" {
		return errors.New("balance.Type is empty")
	}
	if b.CreatedAt.IsZero() {
		return errors.New("balance.CreatedAt is zero")
	}
	if b.UpdatedAt.IsZero() {
		return errors.New("balance.UpdatedAt is zero")
	}
	return nil
}

// balances/{uid}/transactions/{txid}
type BalanceTransaction struct {
	UserID    string    `json:"user_id"`
	TxID      string    `json:"tx_id"`
	Value     int64     `json:"value"`
	Type      AssetType `json:"type"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewBalanceTransaction(uid, txid string, value int64, assetType AssetType, createdAt time.Time, updatedAt time.Time) *BalanceTransaction {
	return &BalanceTransaction{UserID: uid, TxID: txid, Value: value, Type: assetType, CreatedAt: createdAt, UpdatedAt: updatedAt}
}
