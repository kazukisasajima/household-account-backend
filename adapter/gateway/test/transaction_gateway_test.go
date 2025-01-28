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

type TransactionRepositorySuite struct {
	tester.DBSQLiteSuite
	repository gateway.TransactionRepository
}

func TestTransactionRepositorySuite(t *testing.T) {
	suite.Run(t, new(TransactionRepositorySuite))
}

func (suite *TransactionRepositorySuite) SetupSuite() {
	suite.DBSQLiteSuite.SetupSuite()
	suite.repository = gateway.NewTransactionRepository(suite.DB)
}

func (suite *TransactionRepositorySuite) MockDB() sqlmock.Sqlmock {
	mock, mockGormDB := tester.MockDB()
	suite.repository = gateway.NewTransactionRepository(mockGormDB)
	return mock
}

func (suite *TransactionRepositorySuite) AfterTest(suiteName, testName string) {
	suite.repository = gateway.NewTransactionRepository(suite.DB)
}

func (suite *TransactionRepositorySuite) TestTransactionRepositoryCRUD() {
	transaction := &entity.Transaction{
		UserID:     1,
		CategoryID: 1,
		Date:       "2024-01-01",
		Amount:     100.00,
		Content:    "Groceries",
	}
	createdTransaction, err := suite.repository.CreateTransaction(transaction)
	suite.Assert().Nil(err)
	suite.Assert().NotZero(createdTransaction.ID)
	suite.Assert().Equal(100.00, createdTransaction.Amount)
	suite.Assert().Equal("Groceries", createdTransaction.Content)

	getTransaction, err := suite.repository.GetTransactionByID(createdTransaction.ID)
	suite.Assert().Nil(err)
	suite.Assert().Equal(100.00, getTransaction.Amount)
	suite.Assert().Equal("Groceries", getTransaction.Content)

	getTransaction.Amount = 150.00
	updatedTransaction, err := suite.repository.UpdateTransaction(getTransaction)
	suite.Assert().Nil(err)
	suite.Assert().Equal(150.00, updatedTransaction.Amount)
	suite.Assert().Equal("Groceries", updatedTransaction.Content)

	err = suite.repository.DeleteTransaction(createdTransaction.ID)
	suite.Assert().Nil(err)
	deletedTransaction, err := suite.repository.GetTransactionByID(createdTransaction.ID)
	suite.Assert().Nil(deletedTransaction)
	suite.Assert().Equal("record not found", err.Error())
}

func (suite *TransactionRepositorySuite) TestTransactionCreateFailure() {
	mockDB := suite.MockDB()
	mockDB.ExpectBegin()
	mockDB.ExpectExec(regexp.QuoteMeta("INSERT INTO `transactions` (`user_id`,`category_id`,`date`,`amount`,`content`) VALUES (?,?,?,?,?)")).
		WithArgs(1, 1, "2024-01-01", 100.00, "Groceries").
		WillReturnError(errors.New("create error"))
	mockDB.ExpectRollback()

	transaction := &entity.Transaction{
		UserID:     1,
		CategoryID: 1,
		Date:       "2024-01-01",
		Amount:     100.00,
		Content:    "Groceries",
	}

	createdTransaction, err := suite.repository.CreateTransaction(transaction)
	suite.Assert().Nil(createdTransaction)
	suite.Assert().NotNil(err)
	suite.Assert().Equal("create error", err.Error())
}

func (suite *TransactionRepositorySuite) TestTransactionGetFailure() {
	mockDB := suite.MockDB()
	mockDB.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `transactions` WHERE `transactions`.`id` = ? ORDER BY `transactions`.`id` LIMIT ?")).
		WithArgs(1, 1).
		WillReturnError(errors.New("get error"))

	transaction, err := suite.repository.GetTransactionByID(1)
	suite.Assert().Nil(transaction)
	suite.Assert().NotNil(err)
	suite.Assert().Equal("get error", err.Error())
}

func (suite *TransactionRepositorySuite) TestTransactionUpdateFailure() {
	mockDB := suite.MockDB()
	mockDB.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `transactions` WHERE `transactions`.`id` = ? ORDER BY `transactions`.`id` LIMIT ?")).
		WithArgs(1, 1).
		WillReturnError(errors.New("update error"))

	transaction := &entity.Transaction{
		ID:         1,
		UserID:     1,
		CategoryID: 1,
		Date:       "2024-01-01",
		Amount:     100.00,
		Content:    "Groceries",
	}

	transaction, err := suite.repository.UpdateTransaction(transaction)
	suite.Assert().Nil(transaction)
	suite.Assert().NotNil(err)
	suite.Assert().Equal("update error", err.Error())
}

func (suite *TransactionRepositorySuite) TestTransactionDeleteFailure() {
	mockDB := suite.MockDB()
	mockDB.ExpectBegin()
	mockDB.ExpectExec(regexp.QuoteMeta("DELETE FROM `transactions` WHERE `transactions`.`id` = ?")).WithArgs(1).
		WillReturnError(errors.New("delete error"))
	mockDB.ExpectRollback()

	err := suite.repository.DeleteTransaction(1)
	suite.Assert().NotNil(err)
	suite.Assert().Equal("delete error", err.Error())
}
