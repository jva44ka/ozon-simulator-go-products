package data

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jva44ka/ozon-simulator-go-products/internal/domain"
)

type Transactor struct {
	pool               *pgxpool.Pool
	productMetrics     RepositoryMetrics
	reservationMetrics ReservationMetrics
}

func NewTransactor(pool *pgxpool.Pool, productMetrics RepositoryMetrics, reservationMetrics ReservationMetrics) *Transactor {
	return &Transactor{
		pool:               pool,
		productMetrics:     productMetrics,
		reservationMetrics: reservationMetrics,
	}
}

type repositories struct {
	products     domain.ProductRepository
	reservations domain.ReservationRepository
}

func (r repositories) Products() domain.ProductRepository     { return r.products }
func (r repositories) Reservations() domain.ReservationRepository { return r.reservations }

func (t *Transactor) InTransaction(ctx context.Context, fn func(repos domain.Repositories) error) error {
	return pgx.BeginTxFunc(ctx, t.pool, pgx.TxOptions{}, func(tx pgx.Tx) error {
		return fn(repositories{
			products:     NewProductTxPgxRepository(tx, t.productMetrics),
			reservations: NewReservationTxPgxRepository(tx, t.reservationMetrics),
		})
	})
}
