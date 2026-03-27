package domain

import (
	"context"

	"github.com/jva44ka/ozon-simulator-go-products/internal/domain/models"
)

type ProductRepository interface {
	GetProductBySku(ctx context.Context, sku uint64) (*models.Product, error)
	GetProductsBySkus(ctx context.Context, skus []uint64) ([]*models.Product, error)
	UpdateCount(ctx context.Context, products []*models.Product) error
}

type ReservationRepository interface {
	Insert(ctx context.Context, sku uint64, count uint32) (models.Reservation, error)
	GetByIds(ctx context.Context, ids []int64) ([]models.Reservation, error)
	DeleteByIds(ctx context.Context, ids []int64) error
}

type Repositories interface {
	Products() ProductRepository
	Reservations() ReservationRepository
}

type Transactor interface {
	InTransaction(context.Context, func(repos Repositories) error) error
}
