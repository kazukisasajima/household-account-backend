package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/oapi-codegen/runtime/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"household-account-backend/adapter/controller/echo/handler"
	"household-account-backend/adapter/controller/echo/presenter"
	"household-account-backend/entity"
)

type MockTransactionUseCase struct {
	mock.Mock
}

func (m *MockTransactionUseCase) CreateTransaction(transaction *entity.Transaction) (*entity.Transaction, error) {
	args := m.Called(transaction)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Transaction), args.Error(1)
}

func (m *MockTransactionUseCase) GetTransactionByID(transactionID int) (*entity.Transaction, error) {
	args := m.Called(transactionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Transaction), args.Error(1)
}

func (m *MockTransactionUseCase) GetTransactionsByUserID(userID int) ([]entity.Transaction, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.Transaction), args.Error(1)
}

func (m *MockTransactionUseCase) UpdateTransaction(transaction *entity.Transaction) (*entity.Transaction, error) {
	args := m.Called(transaction)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Transaction), args.Error(1)
}

func (m *MockTransactionUseCase) DeleteTransaction(transactionID int) error {
	args := m.Called(transactionID)
	return args.Error(0)
}

func TestCreateTransaction(t *testing.T) {
	e := echo.New()
	mockUseCase := new(MockTransactionUseCase)
	h := handler.NewTransactionHandler(mockUseCase)

	requestBody := presenter.CreateTransactionJSONRequestBody{
		UserId:     1,
		CategoryId: 1,
		Date:       types.Date{Time: time.Now()},
		Amount:     150.75,
		Content:    pointerToString("Groceries"),
	}
	jsonBody, _ := json.Marshal(requestBody)

	req := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewReader(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockTransaction := &entity.Transaction{
		ID:         1,
		UserID:     1,
		CategoryID: 1,
		Date:       requestBody.Date.Time,
		Amount:     float32(requestBody.Amount),
		Content:    "Groceries",
	}

	mockUseCase.On("CreateTransaction", mock.AnythingOfType("*entity.Transaction")).Return(mockTransaction, nil)

	if assert.NoError(t, h.CreateTransaction(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		var response presenter.TransactionResponse
		json.Unmarshal(rec.Body.Bytes(), &response)
		assert.Equal(t, 1, response.Id)
		assert.Equal(t, "Groceries", *response.Content)
		assert.Equal(t, float32(150.75), response.Amount)
	}
}

func TestGetTransactionByID(t *testing.T) {
	e := echo.New()
	mockUseCase := new(MockTransactionUseCase)
	h := handler.NewTransactionHandler(mockUseCase)

	req := httptest.NewRequest(http.MethodGet, "/transactions/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	mockTransaction := &entity.Transaction{
		ID:         1,
		UserID:     1,
		CategoryID: 1,
		Date:       time.Now(),
		Amount:     100.50,
		Content:    "Groceries",
	}

	mockUseCase.On("GetTransactionByID", 1).Return(mockTransaction, nil)

	if assert.NoError(t, h.GetTransactionByID(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		var response presenter.TransactionResponse
		json.Unmarshal(rec.Body.Bytes(), &response)
		assert.Equal(t, 1, response.Id)
		assert.Equal(t, "Groceries", *response.Content)
	}
}

func TestGetTransactionsByUserID(t *testing.T) {
	e := echo.New()
	mockUseCase := new(MockTransactionUseCase)
	h := handler.NewTransactionHandler(mockUseCase)

	req := httptest.NewRequest(http.MethodGet, "/transactions?user_id=1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockTransactions := []entity.Transaction{
		{ID: 1, UserID: 1, CategoryID: 1, Date: time.Now(), Amount: 100.50, Content: "Groceries"},
		{ID: 2, UserID: 1, CategoryID: 2, Date: time.Now(), Amount: 200.00, Content: "Rent"},
	}

	mockUseCase.On("GetTransactionsByUserID", 1).Return(mockTransactions, nil)

	if assert.NoError(t, h.GetTransactionsByUserID(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		var response []presenter.TransactionResponse
		json.Unmarshal(rec.Body.Bytes(), &response)
		assert.Equal(t, 2, len(response))
		assert.Equal(t, "Groceries", *response[0].Content)
		assert.Equal(t, "Rent", *response[1].Content)
	}
}

func TestUpdateTransaction(t *testing.T) {
	e := echo.New()
	mockUseCase := new(MockTransactionUseCase)
	h := handler.NewTransactionHandler(mockUseCase)

	requestBody := presenter.TransactionUpdateRequestBody{
		UserId:     1,
		CategoryId: 1,
		Date:       types.Date{Time: time.Now()},
		Amount:     200.00,
		Content:    pointerToString("Updated Groceries"),
	}
	jsonBody, _ := json.Marshal(requestBody)

	req := httptest.NewRequest(http.MethodPut, "/transactions/1", bytes.NewReader(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	mockTransaction := &entity.Transaction{
		ID:         1,
		UserID:     1,
		CategoryID: 1,
		Date:       requestBody.Date.Time,
		Amount:     200.00,
		Content:    "Updated Groceries",
	}

	mockUseCase.On("UpdateTransaction", mock.AnythingOfType("*entity.Transaction")).Return(mockTransaction, nil)

	if assert.NoError(t, h.UpdateTransaction(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		var response presenter.TransactionResponse
		json.Unmarshal(rec.Body.Bytes(), &response)
		assert.Equal(t, "Updated Groceries", *response.Content)
	}
}

func TestDeleteTransaction(t *testing.T) {
	e := echo.New()
	mockUseCase := new(MockTransactionUseCase)
	h := handler.NewTransactionHandler(mockUseCase)

	req := httptest.NewRequest(http.MethodDelete, "/transactions/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	mockUseCase.On("DeleteTransaction", 1).Return(nil)

	if assert.NoError(t, h.DeleteTransaction(c)) {
		assert.Equal(t, http.StatusNoContent, rec.Code)
	}
}

func pointerToString(s string) *string {
	return &s
}
