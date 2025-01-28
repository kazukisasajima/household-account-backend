package handler

import (
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"

	"household-account-backend/adapter/controller/echo/presenter"
	"household-account-backend/entity"
	"household-account-backend/usecase"
)

type CategoryHandler struct {
	categoryUseCase usecase.CategoryUseCase
}

func NewCategoryHandler(categoryUseCase usecase.CategoryUseCase) *CategoryHandler {
	return &CategoryHandler{
		categoryUseCase: categoryUseCase,
	}
}

func categoryToResponse(category *entity.Category) *presenter.CategoryResponse {
	return &presenter.CategoryResponse{
		Id:     category.ID,
		Name:   category.Name,
		Type:   presenter.CategoryRequestType(category.Type),
	}
}

func (h *CategoryHandler) CreateCategory(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userId := int(claims["user_id"].(float64))

	var requestBody presenter.CreateCategoryJSONRequestBody
	if err := c.Bind(&requestBody); err != nil {
		return c.JSON(http.StatusBadRequest, &presenter.ErrorResponse{Message: "Invalid request body"})
	}

	category := &entity.Category{
		UserID: userId,
		Name:   requestBody.Name,
		Type:   string(requestBody.Type),
	}
	
	createdCategory, err := h.categoryUseCase.CreateCategory(category)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &presenter.ErrorResponse{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, categoryToResponse(createdCategory))
}

func (h *CategoryHandler) GetCategoriesByUserID(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userId := int(claims["user_id"].(float64))

	createdCategory, err := h.categoryUseCase.GetCategoriesByUserID(userId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &presenter.ErrorResponse{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, createdCategory)
}

func (h *CategoryHandler) GetCategoryByID(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userId := int(claims["user_id"].(float64))

	idParam := c.Param("id")
	categoryId, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &presenter.ErrorResponse{Message: "Invalid category ID"})
	}

	category, err := h.categoryUseCase.GetCategoryByID(userId, categoryId)
	if err != nil {
		return c.JSON(http.StatusNotFound, &presenter.ErrorResponse{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, categoryToResponse(category))
}

func (h *CategoryHandler) UpdateCategory(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userId := int(claims["user_id"].(float64))

	idParam := c.Param("id")
	categoryId, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &presenter.ErrorResponse{Message: "Invalid category ID"})
	}

	var requestBody presenter.UpdateCategoryByIdJSONRequestBody
	if err := c.Bind(&requestBody); err != nil {
		return c.JSON(http.StatusBadRequest, &presenter.ErrorResponse{Message: "Invalid request body"})
	}

	category := &entity.Category{
		ID:   categoryId,
		UserID: userId,
		Name: requestBody.Name,
		Type: string(requestBody.Type),
	}

	updatedCategory, err := h.categoryUseCase.UpdateCategory(category)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &presenter.ErrorResponse{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, categoryToResponse(updatedCategory))
}

func (h *CategoryHandler) DeleteCategory(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userId := int(claims["user_id"].(float64))

	idParam := c.Param("id")
	categoryId, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &presenter.ErrorResponse{Message: "Invalid category ID"})
	}

	if err := h.categoryUseCase.DeleteCategory(userId, categoryId); err != nil {
		return c.JSON(http.StatusInternalServerError, &presenter.ErrorResponse{Message: err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}
