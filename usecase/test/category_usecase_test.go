package usecase_test

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"household-account-backend/entity"
	"household-account-backend/usecase"
)

type mockCategoryRepository struct {
	mock.Mock
}

func NewMockCategoryRepository() *mockCategoryRepository {
	return new(mockCategoryRepository)
}

func (m *mockCategoryRepository) CreateCategory(category *entity.Category) (*entity.Category, error) {
	args := m.Called(category)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Category), args.Error(1)
}

func (m *mockCategoryRepository) GetCategoryByID(categoryID int) (*entity.Category, error) {
	args := m.Called(categoryID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Category), args.Error(1)
}

func (m *mockCategoryRepository) GetCategoriesByUserID(userID int) ([]entity.Category, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.Category), args.Error(1)
}

func (m *mockCategoryRepository) UpdateCategory(category *entity.Category) (*entity.Category, error) {
	args := m.Called(category)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Category), args.Error(1)
}

func (m *mockCategoryRepository) DeleteCategory(categoryID int) error {
	args := m.Called(categoryID)
	return args.Error(0)
}

type CategoryUseCaseSuite struct {
	suite.Suite
	categoryUseCase usecase.CategoryUseCase
}

func TestCategoryUseCaseSuite(t *testing.T) {
	suite.Run(t, new(CategoryUseCaseSuite))
}

func (suite *CategoryUseCaseSuite) SetupTest() {
	mockRepo := NewMockCategoryRepository()
	suite.categoryUseCase = usecase.NewCategoryUseCase(mockRepo)
}

func (suite *CategoryUseCaseSuite) TestCreateCategory() {
    category := &entity.Category{
        UserID: 1,
        Name:   "Groceries",
        Type:   "expense",
    }

    mockRepo := NewMockCategoryRepository()
    suite.categoryUseCase = usecase.NewCategoryUseCase(mockRepo)
    mockRepo.On("CreateCategory", category).Return(category, nil)

    createdCategory, err := suite.categoryUseCase.CreateCategory(category)
    suite.Assert().Nil(err)
    suite.Assert().Equal("Groceries", createdCategory.Name)
    suite.Assert().Equal("expense", createdCategory.Type)
}

func (suite *CategoryUseCaseSuite) TestGetCategoryByID() {
    category := &entity.Category{
        ID:     1,
        UserID: 1,
        Name:   "Groceries",
        Type:   "expense",
    }

    mockRepo := NewMockCategoryRepository()
    suite.categoryUseCase = usecase.NewCategoryUseCase(mockRepo)
    mockRepo.On("GetCategoryByID", category.ID).Return(category, nil)

    retrievedCategory, err := suite.categoryUseCase.GetCategoryByID(category.ID)
    suite.Assert().Nil(err)
    suite.Assert().Equal("Groceries", retrievedCategory.Name)
    suite.Assert().Equal("expense", retrievedCategory.Type)
}

func (suite *CategoryUseCaseSuite) TestGetCategoriesByUserID() {
    categories := []entity.Category{
        {
            ID:     1,
            UserID: 1,
            Name:   "Groceries",
            Type:   "expense",
        },
        {
            ID:     2,
            UserID: 1,
            Name:   "Salary",
            Type:   "income",
        },
    }

    mockRepo := NewMockCategoryRepository()
    suite.categoryUseCase = usecase.NewCategoryUseCase(mockRepo)
    mockRepo.On("GetCategoriesByUserID", 1).Return(categories, nil)

    retrievedCategories, err := suite.categoryUseCase.GetCategoriesByUserID(1)
    suite.Assert().Nil(err)
    suite.Assert().Equal(2, len(retrievedCategories))
    suite.Assert().Equal("Groceries", retrievedCategories[0].Name)
    suite.Assert().Equal("Salary", retrievedCategories[1].Name)
}

func (suite *CategoryUseCaseSuite) TestUpdateCategory() {
    category := &entity.Category{
        ID:     1,
        UserID: 1,
        Name:   "Groceries",
        Type:   "expense",
    }

    mockRepo := NewMockCategoryRepository()
    suite.categoryUseCase = usecase.NewCategoryUseCase(mockRepo)
    mockRepo.On("UpdateCategory", category).Return(category, nil)

    updatedCategory, err := suite.categoryUseCase.UpdateCategory(category)
    suite.Assert().Nil(err)
    suite.Assert().Equal("Groceries", updatedCategory.Name)
    suite.Assert().Equal("expense", updatedCategory.Type)
}

func (suite *CategoryUseCaseSuite) TestDeleteCategory() {
    mockRepo := NewMockCategoryRepository()
    suite.categoryUseCase = usecase.NewCategoryUseCase(mockRepo)
    mockRepo.On("DeleteCategory", 1).Return(nil)

    err := suite.categoryUseCase.DeleteCategory(1)
    suite.Assert().Nil(err)
}
