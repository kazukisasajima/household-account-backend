package handler

import (
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"

	"household-account-backend/adapter/controller/echo/presenter"
	"household-account-backend/entity"
	"household-account-backend/pkg/logger"
	"household-account-backend/usecase"
)

type MonthlySummaryHandler struct {
	monthlySummaryUseCase usecase.MonthlySummaryUseCase
}

func NewMonthlySummaryHandler(monthlySummaryUseCase usecase.MonthlySummaryUseCase) *MonthlySummaryHandler {
	return &MonthlySummaryHandler{
		monthlySummaryUseCase: monthlySummaryUseCase,
	}
}

func monthlySummaryToResponse(summary *entity.MonthlySummary) *presenter.MonthlySummaryResponse {
	return &presenter.MonthlySummaryResponse{
		Id:         summary.ID,
		YearMonth:  summary.YearMonth,
		Income:     float32(summary.Income),
		Expense:    float32(summary.Expense),
		Balance:    float32(summary.Balance),
	}
}

func (h *MonthlySummaryHandler) CreateMonthlySummary(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userId := int(claims["user_id"].(float64))

	var requestBody presenter.CreateMonthlySummaryJSONRequestBody
	if err := c.Bind(&requestBody); err != nil {
		logger.Warn(err.Error())
		return c.JSON(http.StatusBadRequest, &presenter.ErrorResponse{Message: err.Error()})
	}

	// Balance = Income - Expense
	// requestBody.Balance = requestBody.Income - requestBody.Expense
	// usecaseで計算するように変更?
	summary := &entity.MonthlySummary{
		UserID:    userId,
		YearMonth: requestBody.YearMonth,
		Income:    float64(requestBody.Income),
		Expense:   float64(requestBody.Expense),
		Balance:   float64(requestBody.Balance),
	}

	createdSummary, err := h.monthlySummaryUseCase.CreateMonthlySummary(summary)
	if err != nil {
		logger.Error(err.Error())
		return c.JSON(http.StatusInternalServerError, &presenter.ErrorResponse{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, monthlySummaryToResponse(createdSummary))
}

func (h *MonthlySummaryHandler) GetMonthlySummaryByID(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userId := int(claims["user_id"].(float64))

	monthlySummaryId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, &presenter.ErrorResponse{Message: "Invalid ID"})
	}

	summary, err := h.monthlySummaryUseCase.GetMonthlySummaryByID(userId, monthlySummaryId)
	if err != nil {
		logger.Error(err.Error())
		return c.JSON(http.StatusNotFound, &presenter.ErrorResponse{Message: "Monthly summary not found"})
	}

	return c.JSON(http.StatusOK, monthlySummaryToResponse(summary))
}

func (h *MonthlySummaryHandler) GetMonthlySummariesByUserID(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userId := int(claims["user_id"].(float64))

	summaries, err := h.monthlySummaryUseCase.GetMonthlySummariesByUserID(userId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &presenter.ErrorResponse{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, summaries)
}

func (h *MonthlySummaryHandler) UpdateMonthlySummary(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userId := int(claims["user_id"].(float64))

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, &presenter.ErrorResponse{Message: "Invalid ID"})
	}

	var requestBody presenter.UpdateMonthlySummaryByIdJSONRequestBody
	if err := c.Bind(&requestBody); err != nil {
		logger.Warn(err.Error())
		return c.JSON(http.StatusBadRequest, &presenter.ErrorResponse{Message: "Invalid request format"})
	}

	summary := &entity.MonthlySummary{
		ID:        id,
		UserID:    userId,
		Income:    float64(requestBody.Income),
		Expense:   float64(requestBody.Expense),
		Balance:   float64(requestBody.Balance),
		YearMonth: requestBody.YearMonth,
	}

	updatedSummary, err := h.monthlySummaryUseCase.UpdateMonthlySummary(summary)
	if err != nil {
		logger.Error(err.Error())
		return c.JSON(http.StatusInternalServerError, &presenter.ErrorResponse{Message: "Failed to update monthly summary"})
	}

	return c.JSON(http.StatusOK, monthlySummaryToResponse(updatedSummary))
}

func (h *MonthlySummaryHandler) DeleteMonthlySummary(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userId := int(claims["user_id"].(float64))

	monthlySummaryId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, &presenter.ErrorResponse{Message: "Invalid ID"})
	}

	if err := h.monthlySummaryUseCase.DeleteMonthlySummary(userId, monthlySummaryId); err != nil {
		logger.Error(err.Error())
		return c.JSON(http.StatusInternalServerError, &presenter.ErrorResponse{Message: "Failed to delete monthly summary"})
	}

	return c.NoContent(http.StatusNoContent)
}
