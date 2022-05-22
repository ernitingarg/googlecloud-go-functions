package btcaccount

import (
	"context"
	"soteria-functions/entity"
)

type Repository interface {
	Save(ctx context.Context, btcAccount *entity.BtcAccount) error
	FetchByAddress(ctx context.Context, address string) (*entity.BtcAccount, error)
}
