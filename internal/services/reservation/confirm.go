package reservation

import (
	"context"
	"fmt"
	"maps"
	"slices"

	"github.com/jackc/pgx/v5"
	"github.com/jva44ka/ozon-simulator-go-products/internal/models"
)

func (s *Service) Confirm(ctx context.Context, ids []int64) error {
	reservations, err := s.db.Reservations().GetByIds(ctx, ids)
	if err != nil {
		return fmt.Errorf("ReservationService.Confirm: %w", err)
	}

	reservationSumsBySku := make(map[uint64]uint32, len(reservations))
	for _, reservation := range reservations {
		reservationSumsBySku[reservation.Sku] += reservation.Count
	}

	skus := slices.Collect(maps.Keys(reservationSumsBySku))
	products, err := s.db.Products().GetBySkus(ctx, skus)
	if err != nil {
		return fmt.Errorf("ReservationService.Confirm: %w", err)
	}

	productMap := make(map[uint64]*models.Product, len(products))
	for _, product := range products {
		productMap[product.Sku] = product
	}

	for _, product := range products {
		delta := reservationSumsBySku[product.Sku]
		productMap[product.Sku].Count -= delta
		productMap[product.Sku].ReservedCount -= delta
	}

	return s.db.InTransaction(ctx, func(tx pgx.Tx) error {
		if err = s.db.Products().WithTx(tx).Update(ctx, slices.Collect(maps.Values(productMap))); err != nil {
			return fmt.Errorf("Confirm: %w", err)
		}
		return s.db.Reservations().WithTx(tx).DeleteByIds(ctx, ids)
	})
}
