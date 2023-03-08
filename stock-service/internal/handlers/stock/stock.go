package stock

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"stock-service/internal/domain/dto"
	"stock-service/internal/services"
	"strconv"
)

const (
	productIDPathArg = "product_id"
	limitQueryParam  = "limit"
	offsetQueryParam = "offset"
)

type handler struct {
	service services.Service
}

func NewHandler(s services.Service) *handler {
	return &handler{service: s}
}

func (h *handler) CreateProduct(ctx echo.Context) error {
	var body CreateProductRequest

	err := ctx.Bind(&body)
	if err != nil {
		return ctx.JSON(http.StatusUnprocessableEntity, ResponseError{Message: err.Error()})
	}

	orderDTO := dto.ProductRequestDTO{
		Quantity: body.Quantity,
		Name:     body.Name,
	}

	cct := ctx.Request().Context()

	id, err := h.service.Create(cct, orderDTO)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
	}

	return ctx.JSON(http.StatusCreated, ResponseCreated{ID: id})
}

func (h *handler) DetailProduct(ctx echo.Context) error {
	productID := ctx.Param(productIDPathArg)
	id, err := strconv.Atoi(productID)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, ResponseError{Message: fmt.Errorf("bad id parameter %v", err).Error()})
	}

	ccx := ctx.Request().Context()

	product, err := h.service.Detail(ccx, int64(id))
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
	}

	return ctx.JSON(http.StatusOK, ResponseProduct{
		ID:       product.ID,
		Quantity: product.Quantity,
		Name:     product.Name,
	})
}

func (h *handler) ListProduct(ctx echo.Context) error {

	limit, err := queryParamsToUInt64(ctx.QueryParam(limitQueryParam), 10)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, ResponseError{Message: fmt.Errorf("wrong limit parameter: %v", err).Error()})
	}

	offset, err := queryParamsToUInt64(ctx.QueryParam(offsetQueryParam), 0)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, ResponseError{Message: fmt.Errorf("wrong offst parameter: %v", err).Error()})
	}

	filter := dto.FilterProductDTO{
		Quantity: nil,
		Limit:    limit,
		Offset:   offset,
	}

	ccx := ctx.Request().Context()

	orders, err := h.service.List(ccx, filter)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
	}

	result := make([]*ResponseProduct, 0, len(orders))

	for _, v := range orders {
		result = append(result, &ResponseProduct{
			ID:       v.ID,
			Quantity: v.Quantity,
			Name:     v.Name,
		})
	}

	res := ResponseProducts{
		Total:   len(orders),
		Results: result,
	}

	return ctx.JSON(http.StatusOK, res)
}

func (h *handler) FillProducts(ctx echo.Context) error {
	var body FillProductsRequest

	err := ctx.Bind(&body)
	if err != nil {
		return ctx.JSON(http.StatusUnprocessableEntity, ResponseError{Message: err.Error()})
	}

	data := make([]dto.FillProductRequestDTO, 0, len(body.Data))
	for _, v := range body.Data {
		data = append(data, dto.FillProductRequestDTO{
			ID:       v.ID,
			Quantity: v.Quantity,
			Name:     v.Name,
		})
	}

	var productsDTO = dto.FillRequestDTO{
		Data: data,
	}
	cct := ctx.Request().Context()

	err = h.service.Fill(cct, productsDTO)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
	}

	return ctx.JSON(http.StatusOK, nil)
}
