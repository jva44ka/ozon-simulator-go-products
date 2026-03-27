package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jva44ka/ozon-simulator-go-products/internal/domain/product"
	"github.com/jva44ka/ozon-simulator-go-products/internal/domain/reservation"
)

type Transactor struct {
	pool               *pgxpool.Pool
	productMetrics     product.RepositoryMetrics
	reservationMetrics reservation.Metrics
}

type Repositories struct {
	Products     product.ProductRepository
	Reservations product.ReservationRepository
}

func NewTransactor(pool *pgxpool.Pool, productMetrics product.RepositoryMetrics, reservationMetrics reservation.Metrics) *Transactor {
	return &Transactor{
		pool:               pool,
		productMetrics:     productMetrics,
		reservationMetrics: reservationMetrics,
	}
}

func (t *Transactor) InTransaction(ctx context.Context, fn func(repos product.Repositories) error) error {
	return pgx.BeginTxFunc(ctx, t.pool, pgx.TxOptions{}, func(tx pgx.Tx) error {
		repos := product.Repositories{
			Products:     product.NewTxPgxRepository(tx, t.productMetrics),
			Reservations: reservation.NewTxPgxRepository(tx, t.reservationMetrics),
		}
		return fn(repos)
	})
}
