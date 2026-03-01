package app

import (
	"context"
	"fmt"

	desc "github.com/jva44ka/ozon-simulator-go-products/internal/app/gen/ozon-simulator-go-products/api/proto"
	"github.com/jva44ka/ozon-simulator-go-products/internal/domain/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ desc.ProductsServer = (*GrpcService)(nil)

type GrpcService struct {
	desc.UnimplementedProductsServer
	ProductService service.ProductService
}

func NewGrpcService(productService service.ProductService) *GrpcService {
	return &GrpcService{ProductService: productService}
}

func (s *GrpcService) GetProduct(ctx context.Context, request *desc.GetProductRequest) (resp *desc.GetProductResponse, err error) {
	if request.Sku < 1 {
		return &desc.GetProductResponse{}, status.Errorf(codes.InvalidArgument, "sku must be more than zero")
	}

	product, err := s.ProductService.GetProductBySku(ctx, request.Sku)
	if err != nil {
		return &desc.GetProductResponse{}, fmt.Errorf("error getting product: %w", err)
	}

	response := &desc.GetProductResponse{
		Sku:   product.Sku,
		Name:  product.Name,
		Price: float32(product.Price),
		Count: product.Count,
	}

	return response, nil
}
