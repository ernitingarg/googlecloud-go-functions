package ethaccount

import (
	"context"
	"soteria-functions/entity"
)

type Repository interface {
	Save(ctx context.Context, btcAccount *entity.EthAccount) error
	FetchByAddress(ctx context.Context, address string) (*entity.EthAccount, error)
}
