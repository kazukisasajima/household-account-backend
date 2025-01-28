package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"household-account-backend/adapter/controller/echo/handler"
	"household-account-backend/entity"
	"household-account-backend/adapter/controller/echo/presenter"
)

type MockCategoryUseCase struct {
	mock.Mock
}

func (m *MockCategoryUseCase) CreateCategory(category *entity.Category) (*entity.Category, error) {
	args := m.Called(category)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Category), args.Error(1)
}

func (m *MockCategoryUseCase) GetCategoryByID(categoryID int) (*entity.Category, error) {
	args := m.Called(categoryID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Category), args.Error(1)
}

func (m *MockCategoryUseCase) GetCategoriesByUserID(userID int) ([]entity.Category, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.Category), args.Error(1)
}

func (m *MockCategoryUseCase) UpdateCategory(category *entity.Category) (*entity.Category, error) {
	args := m.Called(category)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Category), args.Error(1)
}

func (m *MockCategoryUseCase) DeleteCategory(categoryID int) error {
	args := m.Called(categoryID)
	return args.Error(0)
}

func TestCreateCategory(t *testing.T) {
	e := echo.New()
	mockUseCase := new(MockCategoryUseCase)
	h := handler.NewCategoryHandler(mockUseCase)

	requestBody := presenter.CreateCategoryJSONRequestBody{
		UserId: 1,
		Name:   "Groceries",
		Type:   "expense",
	}
	jsonBody, _ := json.Marshal(requestBody)

	req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewReader(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockCategory := &entity.Category{
		ID:     1,
		UserID: 1,
		Name:   "Groceries",
		Type:   "expense",
	}

	mockUseCase.On("CreateCategory", mock.AnythingOfType("*entity.Category")).Return(mockCategory, nil)

	if assert.NoError(t, h.CreateCategory(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		var response presenter.CategoryResponse
		json.Unmarshal(rec.Body.Bytes(), &response)
		assert.Equal(t, 1, response.Id)
		assert.Equal(t, "Groceries", response.Name)
		assert.Equal(t, "expense", string(response.Type))
	}
}

func TestGetCategoryByID(t *testing.T) {
	e := echo.New()
	mockUseCase := new(MockCategoryUseCase)
	h := handler.NewCategoryHandler(mockUseCase)

	req := httptest.NewRequest(http.MethodGet, "/categories/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	mockCategory := &entity.Category{
		ID:     1,
		UserID: 1,
		Name:   "Groceries",
		Type:   "expense",
	}

	mockUseCase.On("GetCategoryByID", 1).Return(mockCategory, nil)

	if assert.NoError(t, h.GetCategoryByID(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		var response presenter.CategoryResponse
		json.Unmarshal(rec.Body.Bytes(), &response)
		assert.Equal(t, 1, response.Id)
		assert.Equal(t, "Groceries", response.Name)
		assert.Equal(t, "expense", string(response.Type))
	}
}

func TestGetCategoriesByUserID(t *testing.T) {
	e := echo.New()
	mockUseCase := new(MockCategoryUseCase)
	h := handler.NewCategoryHandler(mockUseCase)

	req := httptest.NewRequest(http.MethodGet, "/categories?user_id=1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockCategories := []entity.Category{
		{ID: 1, UserID: 1, Name: "Groceries", Type: "expense"},
		{ID: 2, UserID: 1, Name: "Salary", Type: "income"},
	}

	mockUseCase.On("GetCategoriesByUserID", 1).Return(mockCategories, nil)

	if assert.NoError(t, h.GetCategoriesByUserID(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		var response []presenter.CategoryResponse
		json.Unmarshal(rec.Body.Bytes(), &response)
		assert.Len(t, response, 2)
		assert.Equal(t, "Groceries", response[0].Name)
		assert.Equal(t, "Salary", response[1].Name)
	}
}

func TestUpdateCategory(t *testing.T) {
	e := echo.New()
	mockUseCase := new(MockCategoryUseCase)
	h := handler.NewCategoryHandler(mockUseCase)

	requestBody := presenter.CategoryUpdateRequestBody{
		Name: "Updated Name",
		Type: "expense",
	}
	jsonBody, _ := json.Marshal(requestBody)

	req := httptest.NewRequest(http.MethodPut, "/categories/1", bytes.NewReader(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	mockCategory := &entity.Category{
		ID:   1,
		Name: "Updated Name",
		Type: "expense",
	}

	mockUseCase.On("UpdateCategory", mock.AnythingOfType("*entity.Category")).Return(mockCategory, nil)

	if assert.NoError(t, h.UpdateCategory(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		var response presenter.CategoryResponse
		json.Unmarshal(rec.Body.Bytes(), &response)
		assert.Equal(t, 1, response.Id)
		assert.Equal(t, "Updated Name", response.Name)
		assert.Equal(t, "expense", string(response.Type))
	}
}

func TestDeleteCategory(t *testing.T) {
	e := echo.New()
	mockUseCase := new(MockCategoryUseCase)
	h := handler.NewCategoryHandler(mockUseCase)

	req := httptest.NewRequest(http.MethodDelete, "/categories/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	mockUseCase.On("DeleteCategory", 1).Return(nil)

	if assert.NoError(t, h.DeleteCategory(c)) {
		assert.Equal(t, http.StatusNoContent, rec.Code)
	}
}
