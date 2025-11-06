package get_product_by_sku_handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/jva44ka/ozon-simulator-go-products/internal/domain/model"
	httpPkg "github.com/jva44ka/ozon-simulator-go-products/pkg/http"
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
// @Failure      404  {object}  http.ErrorResponse
// @Router       /products/{sku} [get]
func (h GetProductsBySkuHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	skuRaw := r.PathValue("sku")
	sku, err := strconv.Atoi(skuRaw)
	if err != nil {
		if err = httpPkg.NewErrorResponse(w, http.StatusBadRequest, "sku must be more than zero"); err != nil {
			fmt.Println("json.Encode failed ", err)

			return
		}

		return
	}

	if sku < 1 {
		if err = httpPkg.NewErrorResponse(w, http.StatusBadRequest, "sku must be more than zero"); err != nil {
			fmt.Println("json.Encode failed ", err)

			return
		}

		return
	}

	Product, err := h.ProductService.GetProductBySku(r.Context(), uint64(sku))
	if err != nil {
		if err = httpPkg.NewErrorResponse(w, http.StatusInternalServerError, err.Error()); err != nil {

			return
		}

		return
	}

	response := GetProductsResponse{
		Product: GetProductProductResponse{
			Sku:   Product.Sku,
			Name:  Product.Name,
			Price: Product.Price,
		}}

	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(&response); err != nil {
		fmt.Println("success status failed")
		return
	}

	return
}
