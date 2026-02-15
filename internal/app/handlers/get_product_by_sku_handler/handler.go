package get_product_by_sku_handler

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/jva44ka/ozon-simulator-go-products/internal/domain/model"
	httpPkg "github.com/jva44ka/ozon-simulator-go-products/pkg/http/json"
)

type ProductService interface {
	GetProductBySku(ctx context.Context, sku uint64) (*model.Product, error)
}

type GetProductsBySkuHandler struct {
	ProductService ProductService
}

func NewGetProductsBySkuHandler(ProductService ProductService) *GetProductsBySkuHandler {
	return &GetProductsBySkuHandler{ProductService: ProductService}
}

// @Summary      Get product by SKU
// @Description  Возвращает товар по SKU
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        sku  path  int  true  "SKU товара"
// @Success      200  {object}  GetProductsResponse
// @Failure      404  {object}  httpPkg.ErrorResponse
// @Router       /product/{sku} [get]
func (h GetProductsBySkuHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sku, err := parseSku(r)
	if err != nil {
		httpPkg.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	product, err := h.ProductService.GetProductBySku(r.Context(), uint64(sku))
	if err != nil {
		if errors.Is(err, model.ErrProductNotFound) {
			httpPkg.WriteErrorResponse(w, http.StatusNotFound, "Product not found")
			return
		}

		httpPkg.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := GetProductsResponse{
		Product: GetProductProductResponse{
			Sku:   product.Sku,
			Name:  product.Name,
			Price: product.Price,
		}}

	httpPkg.WriteSuccessResponse(w, response)
	return
}

func parseSku(r *http.Request) (int, error) {
	skuRaw := r.PathValue("sku")
	sku, err := strconv.Atoi(skuRaw)
	if err != nil {
		return 0, errors.New("sku must be a number")
	}

	if sku < 1 {
		return 0, errors.New("sku must be more than zero")
	}

	return sku, nil
}
