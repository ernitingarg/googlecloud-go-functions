package entity

import "github.com/pkg/errors"

type EthAccount struct {
	UserID     string `json:"user_id"`
	EncPrivKey string `json:"enc_priv_key"`
	Address    string `json:"address"`
}

type BtcAccount struct {
	UserID     string `json:"user_id"`
	EncPrivKey string `json:"enc_priv_key"`
	Address    string `json:"address"`
}

func NewBtcAccount(uid, encPrivKey, address string) *BtcAccount {
	return &BtcAccount{UserID: uid, EncPrivKey: encPrivKey, Address: address}
}

func NewEthAccount(uid, encPrivKey, address string) *EthAccount {
	return &EthAccount{UserID: uid, EncPrivKey: encPrivKey, Address: address}
}

func (b *BtcAccount) Validate() error {
	if b.UserID == "" {
		return errors.New("user id is empty")
	}
	if b.Address == "" {
		return errors.New("address id is empty")
	}
	return nil
}

func (e *EthAccount) Validate() error {
	if e.UserID == "" {
		return errors.New("user id is empty")
	}
	if e.Address == "" {
		return errors.New("address id is empty")
	}
	return nil
}
