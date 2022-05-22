package price_history

import (
	"context"
	"soteria-functions/entity"
)

type Repository interface {
	Save(ctx context.Context, priceHistory *entity.PriceHistory) error
}
