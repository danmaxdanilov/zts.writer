package dtos

import (
	"github.com/danmaxdanilov/zts.writer/internal/products/dtos"
)

type GetProductByIdQueryResponse struct {
	Product *dtos.ProductDto `json:"product"`
}
