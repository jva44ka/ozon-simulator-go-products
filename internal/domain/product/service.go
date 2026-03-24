package product

import (
	"context"
	"fmt"
	"maps"
	"slices"

	domainErrors "github.com/jva44ka/ozon-simulator-go-products/internal/domain/errors"
	"github.com/jva44ka/ozon-simulator-go-products/internal/domain/reservation"
)

type Repository interface {
	GetProductBySku(ctx context.Context, sku uint64) (*Product, error)
	GetProductsBySkus(ctx context.Context, skus []uint64) ([]*Product, error)
	UpdateCount(ctx context.Context, products []*Product) error
}

type ReservationRepository interface {
	Insert(ctx context.Context, sku uint64, count uint32) (int64, error)
	GetByIds(ctx context.Context, ids []int64) ([]reservation.Reservation, error)
	DeleteByIds(ctx context.Context, ids []int64) error
}

type Service struct {
	repo         Repository
	reservations ReservationRepository
}

func NewService(repo Repository, reservations ReservationRepository) *Service {
	return &Service{
		repo:         repo,
		reservations: reservations,
	}
}

type UpdateCount struct {
	Sku   uint64
	Delta uint32
}

func (s *Service) GetProductBySku(ctx context.Context, sku uint64) (*Product, error) {
	p, err := s.repo.GetProductBySku(ctx, sku)
	if err != nil {
		return nil, fmt.Errorf("productRepository.GetProductBySku: %w", err)
	}
	if p == nil {
		return nil, domainErrors.NewProductNotFoundError(sku)
	}
	return p, nil
}

func (s *Service) IncreaseCount(ctx context.Context, products []UpdateCount) error {
	existingMap, err := s.validateProductsExist(ctx, products)
	if err != nil {
		return err
	}
	for _, p := range products {
		existingMap[p.Sku].Count += p.Delta
	}
	return s.repo.UpdateCount(ctx, slices.Collect(maps.Values(existingMap)))
}

func (s *Service) Reserve(ctx context.Context, products []UpdateCount) (map[uint64]int64, error) {
	existingMap, err := s.validateProductsExist(ctx, products)
	if err != nil {
		return nil, err
	}
	for _, p := range products {
		existing := existingMap[p.Sku]
		if existing.Count < p.Delta {
			return nil, domainErrors.NewInsufficientProductError(p.Sku, existing.Count, p.Delta)
		}
		existing.Count -= p.Delta
	}
	if err = s.repo.UpdateCount(ctx, slices.Collect(maps.Values(existingMap))); err != nil {
		return nil, fmt.Errorf("ProductService.Reserve: %w", err)
	}

	reservationIds := make(map[uint64]int64, len(products))
	for _, p := range products {
		id, err := s.reservations.Insert(ctx, p.Sku, p.Delta)
		if err != nil {
			return nil, fmt.Errorf("ProductService.Reserve: %w", err)
		}
		reservationIds[p.Sku] = id
	}
	return reservationIds, nil
}

func (s *Service) ReleaseReservations(ctx context.Context, ids []int64) error {
	reservations, err := s.reservations.GetByIds(ctx, ids)
	if err != nil {
		return fmt.Errorf("ProductService.ReleaseReservations: %w", err)
	}
	products := make([]UpdateCount, len(reservations))
	for i, r := range reservations {
		products[i] = UpdateCount{Sku: r.Sku, Delta: r.Count}
	}
	if err = s.ReleaseReservation(ctx, products); err != nil {
		return err
	}
	return s.reservations.DeleteByIds(ctx, ids)
}

func (s *Service) ConfirmReservations(ctx context.Context, ids []int64) error {
	return s.reservations.DeleteByIds(ctx, ids)
}

// ReleaseReservation возвращает count без затрагивания таблицы reservations.
// Используется джобой истечения резерваций.
func (s *Service) ReleaseReservation(ctx context.Context, products []UpdateCount) error {
	existingMap, err := s.validateProductsExist(ctx, products)
	if err != nil {
		return err
	}
	for _, p := range products {
		existingMap[p.Sku].Count += p.Delta
	}
	return s.repo.UpdateCount(ctx, slices.Collect(maps.Values(existingMap)))
}

func (s *Service) validateProductsExist(ctx context.Context, products []UpdateCount) (map[uint64]*Product, error) {
	skus := make([]uint64, 0, len(products))
	for _, p := range products {
		skus = append(skus, p.Sku)
	}

	existing, err := s.repo.GetProductsBySkus(ctx, skus)
	if err != nil {
		return nil, fmt.Errorf("ProductService.validateProductsExist: %w", err)
	}

	existingMap := make(map[uint64]*Product, len(existing))
	for _, p := range existing {
		existingMap[p.Sku] = p
	}

	for _, p := range products {
		if _, ok := existingMap[p.Sku]; !ok {
			return nil, domainErrors.NewProductNotFoundError(p.Sku)
		}
	}

	return existingMap, nil
}
