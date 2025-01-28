package gateway_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"

	"household-account-backend/adapter/gateway"
	"household-account-backend/entity"
	"household-account-backend/pkg/tester"
)

type MonthlySummaryRepositorySuite struct {
	tester.DBSQLiteSuite
	repository gateway.MonthlySummaryRepository
}

func TestMonthlySummaryRepositorySuite(t *testing.T) {
	suite.Run(t, new(MonthlySummaryRepositorySuite))
}

func (suite *MonthlySummaryRepositorySuite) SetupSuite() {
	suite.DBSQLiteSuite.SetupSuite()
	suite.repository = gateway.NewMonthlySummaryRepository(suite.DB)
}

func (suite *MonthlySummaryRepositorySuite) MockDB() sqlmock.Sqlmock {
	mock, mockGormDB := tester.MockDB()
	suite.repository = gateway.NewMonthlySummaryRepository(mockGormDB)
	return mock
}

func (suite *MonthlySummaryRepositorySuite) AfterTest(suiteName, testName string) {
	suite.repository = gateway.NewMonthlySummaryRepository(suite.DB)
}

func (suite *MonthlySummaryRepositorySuite) TestMonthlySummaryRepositoryCRUD() {
	summary := &entity.MonthlySummary{
		UserID:    1,
		YearMonth: "2024-01",
		Income:    5000.00,
		Expense:   3000.00,
		Balance:   2000.00,
	}
	createdSummary, err := suite.repository.CreateMonthlySummary(summary)
	suite.Assert().Nil(err)
	suite.Assert().NotZero(createdSummary.ID)
	suite.Assert().Equal(5000.00, createdSummary.Income)
	suite.Assert().Equal(2000.00, createdSummary.Balance)

	getSummary, err := suite.repository.GetMonthlySummaryByID(createdSummary.ID)
	suite.Assert().Nil(err)
	suite.Assert().Equal(5000.00, getSummary.Income)
	suite.Assert().Equal(2000.00, getSummary.Balance)

	getSummary.Expense = 3500.00
	updatedSummary, err := suite.repository.UpdateMonthlySummary(getSummary)
	suite.Assert().Nil(err)
	suite.Assert().Equal(3500.00, updatedSummary.Expense)

	err = suite.repository.DeleteMonthlySummary(createdSummary.ID)
	suite.Assert().Nil(err)
	deletedSummary, err := suite.repository.GetMonthlySummaryByID(createdSummary.ID)
	suite.Assert().Nil(deletedSummary)
	suite.Assert().Equal("record not found", err.Error())
}
