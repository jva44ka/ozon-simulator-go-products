package reservation

import (
	"github.com/jva44ka/ozon-simulator-go-products/internal/services"
)

type Service struct {
	db services.DBManager
}

func NewService(db services.DBManager) *Service {
	return &Service{db: db}
}

type ReserveItem struct {
	Sku   uint64
	Delta uint32
}
