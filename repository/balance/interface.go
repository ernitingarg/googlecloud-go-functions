package balance

import (
	"context"
	"soteria-functions/entity"
)

type Repository interface {
	Save(ctx context.Context, btcAccount *entity.BtcAccount) error
}
