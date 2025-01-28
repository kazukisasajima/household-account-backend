package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/oapi-codegen/runtime/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"household-account-backend/adapter/controller/echo/handler"
	"household-account-backend/adapter/controller/echo/presenter"
	"household-account-backend/entity"
)

type MockUserUseCase struct {
	mock.Mock
}

func (m *MockUserUseCase) Signup(user *entity.User) (*entity.User, error) {
	args := m.Called(user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserUseCase) Login(credentials *entity.Credentials) (string, error) {
	args := m.Called(credentials)
	return args.String(0), args.Error(1)
}

func (m *MockUserUseCase) GetCurrentUser(userID int) (*entity.User, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserUseCase) UpdateUser(user *entity.User) (*entity.User, error) {
	args := m.Called(user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserUseCase) DeleteUser(userID int) error {
	args := m.Called(userID)
	return args.Error(0)
}

func TestSignup(t *testing.T) {
	e := echo.New()
	mockUseCase := new(MockUserUseCase)
	h := handler.NewUserHandler(mockUseCase)

	requestBody := presenter.CreateUserJSONRequestBody{
		Email:    "test@example.com",
		Password: "password123",
		Name:     "John",
	}
	jsonBody, _ := json.Marshal(requestBody)

	req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewReader(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockUser := &entity.User{
		ID:       1,
		Email:    "test@example.com",
		Password: "hashed_password",
		Name:     "John",
	}

	mockUseCase.On("Signup", mock.AnythingOfType("*entity.User")).Return(mockUser, nil)

	if assert.NoError(t, h.Signup(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		var response presenter.UserResponse
		json.Unmarshal(rec.Body.Bytes(), &response)
		assert.Equal(t, 1, response.Id)
		assert.Equal(t, types.Email("test@example.com"), response.Email)
		assert.Equal(t, "John", response.Name)
	}
}

func TestLogin(t *testing.T) {
	e := echo.New()
	mockUseCase := new(MockUserUseCase)
	h := handler.NewUserHandler(mockUseCase)

	requestBody := entity.Credentials{
		Email:    "test@example.com",
		Password: "password123",
	}
	jsonBody, _ := json.Marshal(requestBody)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	token := "dummy_jwt_token"
	mockUseCase.On("Login", mock.AnythingOfType("*entity.Credentials")).Return(token, nil)

	if assert.NoError(t, h.Login(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		cookie := rec.Result().Cookies()
		assert.Len(t, cookie, 1)
		assert.Equal(t, "auth_token", cookie[0].Name)
		assert.Equal(t, token, cookie[0].Value)
	}
}

func TestGetCurrentUser(t *testing.T) {
	e := echo.New()
	mockUseCase := new(MockUserUseCase)
	h := handler.NewUserHandler(mockUseCase)

	req := httptest.NewRequest(http.MethodGet, "/user", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	userID := 1
	mockUser := &entity.User{
		ID:    userID,
		Email: "test@example.com",
		Name:  "John",
	}

	mockUseCase.On("GetCurrentUser", userID).Return(mockUser, nil)
	c.Set("user", &jwt.Token{
		Claims: jwt.MapClaims{
			"user_id": float64(userID),
		},
	})

	if assert.NoError(t, h.GetCurrentUser(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		var response presenter.UserResponse
		json.Unmarshal(rec.Body.Bytes(), &response)
		assert.Equal(t, userID, response.Id)
		assert.Equal(t, types.Email("test@example.com"), response.Email)
		assert.Equal(t, "John", response.Name)
	}
}

func TestLogout(t *testing.T) {
	e := echo.New()
	mockUseCase := new(MockUserUseCase)
	h := handler.NewUserHandler(mockUseCase)

	req := httptest.NewRequest(http.MethodPost, "/logout", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, h.Logout(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		cookie := rec.Result().Cookies()
		assert.Len(t, cookie, 1)
		assert.Equal(t, "auth_token", cookie[0].Name)
		assert.Empty(t, cookie[0].Value)
		assert.True(t, cookie[0].Expires.Before(time.Now()))
	}
}

func TestDeleteUser(t *testing.T) {
	e := echo.New()
	mockUseCase := new(MockUserUseCase)
	h := handler.NewUserHandler(mockUseCase)

	req := httptest.NewRequest(http.MethodDelete, "/user", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	userID := 1
	mockUseCase.On("DeleteUser", userID).Return(nil)

	c.Set("user", &jwt.Token{
		Claims: jwt.MapClaims{
			"user_id": float64(userID),
		},
	})

	if assert.NoError(t, h.DeleteUser(c)) {
		assert.Equal(t, http.StatusNoContent, rec.Code)
	}
}
