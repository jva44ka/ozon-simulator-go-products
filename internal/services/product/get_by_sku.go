package product

import (
	"context"
	"fmt"

	"github.com/jva44ka/ozon-simulator-go-products/internal/errors"
	"github.com/jva44ka/ozon-simulator-go-products/internal/models"
)

func (s *Service) GetBySku(ctx context.Context, sku uint64) (*models.Product, error) {
	product, err := s.db.Products().GetBySku(ctx, sku)
	if err != nil {
		return nil, fmt.Errorf("productRepository.GetBySku: %w", err)
	}
	if product == nil {
		return nil, errors.NewProductNotFoundError(sku)
	}

	reservations, err := s.db.Reservations().GetSumBySkus(ctx, []uint64{sku})
	if err != nil {
		return nil, fmt.Errorf("productRepository.GetBySku: %w", err)
	}

	if reserved := reservations[sku]; reserved > 0 {
		if product.Count >= reserved {
			product.Count -= reserved
			//TODO: резерваций меньше чем оставшихся продуктов
		} else {
			product.Count = 0
			//TODO: резерваций больше чем оставшихся продуктов, алертим
		}
	}

	return product, nil
}
