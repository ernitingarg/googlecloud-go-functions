package state

import (
	"context"
	"soteria-functions/entity"
)

type Repository interface {
	Save(ctx context.Context, balance *entity.Balance) error
	AddAmount(ctx context.Context, userID string, assetType entity.AssetType, increment int64) error
	Fetch(ctx context.Context, userID string, assetType entity.AssetType) (*entity.Balance, error)
}
