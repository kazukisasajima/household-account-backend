package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"household-account-backend/adapter/controller/echo/handler"
	"household-account-backend/adapter/controller/echo/presenter"
	"household-account-backend/entity"
)

type MockMonthlySummaryUseCase struct {
	mock.Mock
}

func (m *MockMonthlySummaryUseCase) CreateMonthlySummary(summary *entity.MonthlySummary) (*entity.MonthlySummary, error) {
	args := m.Called(summary)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.MonthlySummary), args.Error(1)
}

func (m *MockMonthlySummaryUseCase) GetMonthlySummaryByID(summaryID int) (*entity.MonthlySummary, error) {
	args := m.Called(summaryID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.MonthlySummary), args.Error(1)
}

func (m *MockMonthlySummaryUseCase) GetMonthlySummariesByUserID(userID int) ([]entity.MonthlySummary, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.MonthlySummary), args.Error(1)
}

func (m *MockMonthlySummaryUseCase) UpdateMonthlySummary(summary *entity.MonthlySummary) (*entity.MonthlySummary, error) {
	args := m.Called(summary)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.MonthlySummary), args.Error(1)
}

func (m *MockMonthlySummaryUseCase) DeleteMonthlySummary(summaryID int) error {
	args := m.Called(summaryID)
	return args.Error(0)
}

func TestCreateMonthlySummary(t *testing.T) {
	e := echo.New()
	mockUseCase := new(MockMonthlySummaryUseCase)
	h := handler.NewMonthlySummaryHandler(mockUseCase)

	requestBody := presenter.CreateMonthlySummaryJSONRequestBody{
		UserId:    1,
		YearMonth: "2023-12",
		Income:    1000.0,
		Expense:   500.0,
		Balance:   500.0,
	}
	jsonBody, _ := json.Marshal(requestBody)

	req := httptest.NewRequest(http.MethodPost, "/monthly-summaries", bytes.NewReader(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockSummary := &entity.MonthlySummary{
		ID:        1,
		UserID:    1,
		YearMonth: "2023-12",
		Income:    1000.0,
		Expense:   500.0,
		Balance:   500.0,
	}

	mockUseCase.On("CreateMonthlySummary", mock.AnythingOfType("*entity.MonthlySummary")).Return(mockSummary, nil)

	if assert.NoError(t, h.CreateMonthlySummary(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		var response presenter.MonthlySummaryResponse
		json.Unmarshal(rec.Body.Bytes(), &response)
		assert.Equal(t, 1, response.Id)
		assert.Equal(t, "2023-12", response.YearMonth)
		assert.Equal(t, float32(1000.0), response.Income)
		assert.Equal(t, float32(500.0), response.Expense)
		assert.Equal(t, float32(500.0), response.Balance)
	}
}

func TestGetMonthlySummaryByID(t *testing.T) {
	e := echo.New()
	mockUseCase := new(MockMonthlySummaryUseCase)
	h := handler.NewMonthlySummaryHandler(mockUseCase)

	summaryID := 1
	req := httptest.NewRequest(http.MethodGet, "/monthly-summaries/"+strconv.Itoa(summaryID), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(summaryID))

	mockSummary := &entity.MonthlySummary{
		ID:        summaryID,
		UserID:    1,
		YearMonth: "2023-12",
		Income:    1000.0,
		Expense:   500.0,
		Balance:   500.0,
	}

	mockUseCase.On("GetMonthlySummaryByID", summaryID).Return(mockSummary, nil)

	if assert.NoError(t, h.GetMonthlySummaryByID(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		var response presenter.MonthlySummaryResponse
		json.Unmarshal(rec.Body.Bytes(), &response)
		assert.Equal(t, summaryID, response.Id)
		assert.Equal(t, "2023-12", response.YearMonth)
	}
}

func TestGetMonthlySummariesByUserID(t *testing.T) {
	e := echo.New()
	mockUseCase := new(MockMonthlySummaryUseCase)
	h := handler.NewMonthlySummaryHandler(mockUseCase)

	userID := 1
	req := httptest.NewRequest(http.MethodGet, "/monthly-summaries?user_id="+strconv.Itoa(userID), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockSummaries := []entity.MonthlySummary{
		{
			ID:        1,
			UserID:    userID,
			YearMonth: "2023-12",
			Income:    1000.0,
			Expense:   500.0,
			Balance:   500.0,
		},
		{
			ID:        2,
			UserID:    userID,
			YearMonth: "2024-01",
			Income:    1500.0,
			Expense:   700.0,
			Balance:   800.0,
		},
	}

	mockUseCase.On("GetMonthlySummariesByUserID", userID).Return(mockSummaries, nil)

	if assert.NoError(t, h.GetMonthlySummariesByUserID(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		var response []presenter.MonthlySummaryResponse
		json.Unmarshal(rec.Body.Bytes(), &response)
		assert.Len(t, response, 2)
		assert.Equal(t, "2023-12", response[0].YearMonth)
		assert.Equal(t, "2024-01", response[1].YearMonth)
	}
}

func TestUpdateMonthlySummary(t *testing.T) {
	e := echo.New()
	mockUseCase := new(MockMonthlySummaryUseCase)
	h := handler.NewMonthlySummaryHandler(mockUseCase)

	summaryID := 1
	requestBody := presenter.UpdateMonthlySummaryByIdJSONRequestBody{
		Income:  1200.0,
		Expense: 600.0,
		Balance: 600.0,
	}
	jsonBody, _ := json.Marshal(requestBody)

	req := httptest.NewRequest(http.MethodPut, "/monthly-summaries/"+strconv.Itoa(summaryID), bytes.NewReader(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(summaryID))

	mockSummary := &entity.MonthlySummary{
		ID:        summaryID,
		Income:    1200.0,
		Expense:   600.0,
		Balance:   600.0,
	}

	mockUseCase.On("UpdateMonthlySummary", mock.AnythingOfType("*entity.MonthlySummary")).Return(mockSummary, nil)

	if assert.NoError(t, h.UpdateMonthlySummary(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		var response presenter.MonthlySummaryResponse
		json.Unmarshal(rec.Body.Bytes(), &response)
		assert.Equal(t, summaryID, response.Id)
		assert.Equal(t, float32(1200.0), response.Income)
		assert.Equal(t, float32(600.0), response.Expense)
		assert.Equal(t, float32(600.0), response.Balance)
	}
}

func TestDeleteMonthlySummary(t *testing.T) {
	e := echo.New()
	mockUseCase := new(MockMonthlySummaryUseCase)
	h := handler.NewMonthlySummaryHandler(mockUseCase)

	summaryID := 1
	req := httptest.NewRequest(http.MethodDelete, "/monthly-summaries/"+strconv.Itoa(summaryID), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(summaryID))

	mockUseCase.On("DeleteMonthlySummary", summaryID).Return(nil)

	if assert.NoError(t, h.DeleteMonthlySummary(c)) {
		assert.Equal(t, http.StatusNoContent, rec.Code)
	}
}
