package domain

import (
	"context"

	"github.com/jackc/pgx/v5"
	models2 "github.com/jva44ka/ozon-simulator-go-products/internal/models"
)

type ProductReadRepository interface {
	GetProductBySku(ctx context.Context, sku uint64) (*models2.Product, error)
	GetProductsBySkus(ctx context.Context, skus []uint64) ([]*models2.Product, error)
	WithTx(tx pgx.Tx) ProductWriteRepository
}

type ProductWriteRepository interface {
	UpdateCount(ctx context.Context, products []*models2.Product) error
}

type ReservationReadRepository interface {
	GetByIds(ctx context.Context, ids []int64) ([]models2.Reservation, error)
	WithTx(tx pgx.Tx) ReservationWriteRepository
}

type ReservationWriteRepository interface {
	Insert(ctx context.Context, sku uint64, count uint32) (models2.Reservation, error)
	DeleteByIds(ctx context.Context, ids []int64) error
}

type DBManager interface {
	Products() ProductReadRepository
	Reservations() ReservationReadRepository
	InTransaction(ctx context.Context, fn func(tx pgx.Tx) error) error
}
