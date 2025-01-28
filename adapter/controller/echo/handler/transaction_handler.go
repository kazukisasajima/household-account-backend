package handler

import (
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/oapi-codegen/runtime/types"

	"household-account-backend/adapter/controller/echo/presenter"
	"household-account-backend/entity"
	"household-account-backend/pkg/logger"
	"household-account-backend/usecase"
)

type TransactionHandler struct {
	transactionUseCase usecase.TransactionUseCase
}

func NewTransactionHandler(transactionUseCase usecase.TransactionUseCase) *TransactionHandler {
	return &TransactionHandler{
		transactionUseCase: transactionUseCase,
	}
}

func transactionToResponse(transaction *entity.Transaction) *presenter.TransactionResponse {
	return &presenter.TransactionResponse{
		Id:         transaction.ID,
		UserId:     transaction.UserID,
		CategoryId: transaction.CategoryID,
		Date:       types.Date{Time: transaction.Date},
		Amount:     float32(transaction.Amount),
		Content:    &transaction.Content,
	}
}

func (h *TransactionHandler) CreateTransaction(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userId := int(claims["user_id"].(float64))

	var requestBody presenter.CreateTransactionJSONRequestBody
	if err := c.Bind(&requestBody); err != nil {
		logger.Warn(err.Error())
		return c.JSON(http.StatusBadRequest, &presenter.ErrorResponse{Message: err.Error()})
	}

	transaction := &entity.Transaction{
		UserID:     userId,
		CategoryID: requestBody.CategoryId,
		Date:       requestBody.Date.Time,
		Amount:     float32(requestBody.Amount),
		Content:    *requestBody.Content,
	}

	createdTransaction, err := h.transactionUseCase.CreateTransaction(transaction)
	if err != nil {
		logger.Error(err.Error())
		return c.JSON(http.StatusInternalServerError, &presenter.ErrorResponse{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, transactionToResponse(createdTransaction))
}

func (h *TransactionHandler) GetTransactionsByUserID(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userId := int(claims["user_id"].(float64))

	transactions, err := h.transactionUseCase.GetTransactionsByUserID(userId)
	if err != nil {
		logger.Error(err.Error())
		return c.JSON(http.StatusInternalServerError, &presenter.ErrorResponse{Message: "Failed to retrieve transactions"})
	}

	var response []presenter.TransactionResponse
	for _, transaction := range transactions {
		response = append(response, *transactionToResponse(&transaction))
	}

	return c.JSON(http.StatusOK, response)
}

func (h *TransactionHandler) GetTransactionByID(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userId := int(claims["user_id"].(float64))

	transactionId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, &presenter.ErrorResponse{Message: "Invalid ID"})
	}

	transaction, err := h.transactionUseCase.GetTransactionByID(userId, transactionId)
	if err != nil {
		logger.Error(err.Error())
		return c.JSON(http.StatusNotFound, &presenter.ErrorResponse{Message: "Transaction not found"})
	}

	return c.JSON(http.StatusOK, transactionToResponse(transaction))
}

func (h *TransactionHandler) UpdateTransaction(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userId := int(claims["user_id"].(float64))

	transactionId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, &presenter.ErrorResponse{Message: "Invalid ID"})
	}

	var requestBody presenter.TransactionUpdateRequestBody
	if err := c.Bind(&requestBody); err != nil {
		logger.Warn(err.Error())
		return c.JSON(http.StatusBadRequest, &presenter.ErrorResponse{Message: "Invalid request format"})
	}

	transaction := &entity.Transaction{
		ID:         transactionId,
		UserID:     userId,
		CategoryID: requestBody.CategoryId,
		Date:       requestBody.Date.Time,
		Amount:     float32(requestBody.Amount),
		Content:    *requestBody.Content,
	}

	updatedTransaction, err := h.transactionUseCase.UpdateTransaction(transaction)
	if err != nil {
		logger.Error(err.Error())
		return c.JSON(http.StatusInternalServerError, &presenter.ErrorResponse{Message: "Failed to update transaction"})
	}

	return c.JSON(http.StatusOK, transactionToResponse(updatedTransaction))
}

func (h *TransactionHandler) DeleteTransaction(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userId := int(claims["user_id"].(float64))
	
	transactionId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, &presenter.ErrorResponse{Message: "Invalid ID"})
	}

	if err := h.transactionUseCase.DeleteTransaction(userId, transactionId); err != nil {
		logger.Error(err.Error())
		return c.JSON(http.StatusInternalServerError, &presenter.ErrorResponse{Message: "Failed to delete transaction"})
	}

	return c.NoContent(http.StatusNoContent)
}
