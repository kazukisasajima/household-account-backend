package gateway_test

import (
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"

	"household-account-backend/adapter/gateway"
	"household-account-backend/entity"
	"household-account-backend/pkg/tester"
)

type CategoryRepositorySuite struct {
	tester.DBSQLiteSuite
	repository gateway.CategoryRepository
}

func TestCategoryRepositorySuite(t *testing.T) {
	suite.Run(t, new(CategoryRepositorySuite))
}

func (suite *CategoryRepositorySuite) SetupSuite() {
	suite.DBSQLiteSuite.SetupSuite()
	suite.repository = gateway.NewCategoryRepository(suite.DB)
}

func (suite *CategoryRepositorySuite) MockDB() sqlmock.Sqlmock {
	mock, mockGormDB := tester.MockDB()
	suite.repository = gateway.NewCategoryRepository(mockGormDB)
	return mock
}

func (suite *CategoryRepositorySuite) AfterTest(suiteName, testName string) {
	suite.repository = gateway.NewCategoryRepository(suite.DB)
}

func (suite *CategoryRepositorySuite) TestCategoryRepositoryCRUD() {
	category := &entity.Category{
		UserID: 1,
		Name:   "Food",
		Type:   "expense",
	}
	createdCategory, err := suite.repository.CreateCategory(category)
	suite.Assert().Nil(err)
	suite.Assert().NotZero(createdCategory.ID)
	suite.Assert().Equal("Food", createdCategory.Name)
	suite.Assert().Equal("expense", createdCategory.Type)

	getCategory, err := suite.repository.GetCategoryByID(createdCategory.ID)
	suite.Assert().Nil(err)
	suite.Assert().Equal("Food", getCategory.Name)
	suite.Assert().Equal("expense", getCategory.Type)

	getCategory.Name = "Groceries"
	updatedCategory, err := suite.repository.UpdateCategory(getCategory)
	suite.Assert().Nil(err)
	suite.Assert().Equal("Groceries", updatedCategory.Name)
	suite.Assert().Equal("expense", updatedCategory.Type)

	err = suite.repository.DeleteCategory(createdCategory.ID)
	suite.Assert().Nil(err)
	deletedCategory, err := suite.repository.GetCategoryByID(createdCategory.ID)
	suite.Assert().Nil(deletedCategory)
	suite.Assert().Equal("record not found", err.Error())
}

func (suite *CategoryRepositorySuite) TestCategoryCreateFailure() {
	mockDB := suite.MockDB()
	mockDB.ExpectBegin()
	mockDB.ExpectExec(regexp.QuoteMeta("INSERT INTO `categories` (`user_id`,`name`,`type`) VALUES (?,?,?)")).
		WithArgs(1, "Food", "expense").
		WillReturnError(errors.New("create error"))
	mockDB.ExpectRollback()

	category := &entity.Category{
		UserID: 1,
		Name:   "Food",
		Type:   "expense",
	}

	createdCategory, err := suite.repository.CreateCategory(category)
	suite.Assert().Nil(createdCategory)
	suite.Assert().NotNil(err)
	suite.Assert().Equal("create error", err.Error())
}

func (suite *CategoryRepositorySuite) TestCategoryGetFailure() {
	mockDB := suite.MockDB()
	mockDB.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `categories` WHERE `categories`.`id` = ? ORDER BY `categories`.`id` LIMIT ?")).
		WithArgs(1, 1).
		WillReturnError(errors.New("get error"))

	category, err := suite.repository.GetCategoryByID(1)
	suite.Assert().Nil(category)
	suite.Assert().NotNil(err)
	suite.Assert().Equal("get error", err.Error())
}

func (suite *CategoryRepositorySuite) TestCategoryUpdateFailure() {
	mockDB := suite.MockDB()
	mockDB.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `categories` WHERE `categories`.`id` = ? ORDER BY `categories`.`id` LIMIT ?")).
		WithArgs(1, 1).
		WillReturnError(errors.New("update error"))

	category := &entity.Category{
		ID:   1,
		Name: "Groceries",
		Type: "expense",
	}

	updatedCategory, err := suite.repository.UpdateCategory(category)
	suite.Assert().Nil(updatedCategory)
	suite.Assert().NotNil(err)
	suite.Assert().Equal("update error", err.Error())
}

func (suite *CategoryRepositorySuite) TestCategoryDeleteFailure() {
	mockDB := suite.MockDB()
	mockDB.ExpectBegin()
	mockDB.ExpectExec(regexp.QuoteMeta("DELETE FROM `categories` WHERE `categories`.`id` = ?")).WithArgs(1).
		WillReturnError(errors.New("delete error"))
	mockDB.ExpectRollback()

	err := suite.repository.DeleteCategory(1)
	suite.Assert().NotNil(err)
	suite.Assert().Equal("delete error", err.Error())
}
