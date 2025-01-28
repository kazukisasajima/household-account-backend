package usecase

import (
	"household-account-backend/adapter/gateway"
	"household-account-backend/entity"
)

type CategoryUseCase interface {
	CreateCategory(category *entity.Category) (*entity.Category, error)
	GetCategoryByID(userID int, categoryID int) (*entity.Category, error)
	GetCategoriesByUserID(userID int) ([]entity.Category, error)
	UpdateCategory(category *entity.Category) (*entity.Category, error)
	DeleteCategory(userID int, categoryID int) error
}

type categoryUseCase struct {
	categoryRepository gateway.CategoryRepository
}

func NewCategoryUseCase(categoryRepository gateway.CategoryRepository) CategoryUseCase {
	return &categoryUseCase{
		categoryRepository: categoryRepository,
	}
}

func (cu *categoryUseCase) CreateCategory(category *entity.Category) (*entity.Category, error) {
	return cu.categoryRepository.CreateCategory(category)
}

func (cu *categoryUseCase) GetCategoryByID(userID int, categoryID int) (*entity.Category, error) {
	return cu.categoryRepository.GetCategoryByID(userID, categoryID)
}

func (cu *categoryUseCase) GetCategoriesByUserID(userID int) ([]entity.Category, error) {
	return cu.categoryRepository.GetCategoriesByUserID(userID)
}

func (cu *categoryUseCase) UpdateCategory(category *entity.Category) (*entity.Category, error) {
	return cu.categoryRepository.UpdateCategory(category)
}

func (cu *categoryUseCase) DeleteCategory(userID int, categoryID int) error {
	return cu.categoryRepository.DeleteCategory(userID, categoryID)
}
