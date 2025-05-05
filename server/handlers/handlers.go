package handlers

import (
	calculationservice "calculator/internal/calculationService"
	"net/http"

	"github.com/labstack/echo/v4"
)

type CalculationHandler struct {
	service calculationservice.CalculationService
}

func NewCalculationHandler(s calculationservice.CalculationService) *CalculationHandler {
	return &CalculationHandler{service: s}
}

func (h *CalculationHandler) GetCalculations(ctx echo.Context) error {
	calculations, err := h.service.GetAllCalculations()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "could not get calculations"})

	}
	return ctx.JSON(http.StatusOK, calculations)
}

func (h *CalculationHandler) PostCalculations(ctx echo.Context) error {
	var req calculationservice.CalculationRequest

	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	calc, err := h.service.CreateCalculation(req.Expression)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Could not create calculation"})
	}
	return ctx.JSON(http.StatusOK, calc)
}

func (h *CalculationHandler) PatchCalculations(ctx echo.Context) error {
	id := ctx.Param("id")

	var req calculationservice.CalculationRequest
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	updCalculation, err := h.service.UpdateCalculation(id, req.Expression)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]string{"error": "could not update calculate"})
	}

	return ctx.JSON(http.StatusOK, updCalculation)
}

func (h *CalculationHandler) DeleteCalculations(ctx echo.Context) error {
	id := ctx.Param("id")

	if err := h.service.DeleteCalculation(id).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "could not delete calculation"})

	}

	return ctx.NoContent(http.StatusNoContent)
}
