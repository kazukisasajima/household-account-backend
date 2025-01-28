package usecase_test

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"household-account-backend/entity"
	"household-account-backend/usecase"
)

type mockMonthlySummaryRepository struct {
	mock.Mock
}

func NewMockMonthlySummaryRepository() *mockMonthlySummaryRepository {
	return new(mockMonthlySummaryRepository)
}

func (m *mockMonthlySummaryRepository) CreateMonthlySummary(summary *entity.MonthlySummary) (*entity.MonthlySummary, error) {
	args := m.Called(summary)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.MonthlySummary), args.Error(1)
}

func (m *mockMonthlySummaryRepository) GetMonthlySummaryByID(summaryID int) (*entity.MonthlySummary, error) {
	args := m.Called(summaryID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.MonthlySummary), args.Error(1)
}

func (m *mockMonthlySummaryRepository) GetMonthlySummariesByUserID(userID int) ([]entity.MonthlySummary, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.MonthlySummary), args.Error(1)
}

func (m *mockMonthlySummaryRepository) UpdateMonthlySummary(summary *entity.MonthlySummary) (*entity.MonthlySummary, error) {
	args := m.Called(summary)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.MonthlySummary), args.Error(1)
}

func (m *mockMonthlySummaryRepository) DeleteMonthlySummary(summaryID int) error {
	args := m.Called(summaryID)
	return args.Error(0)
}

type MonthlySummaryUseCaseSuite struct {
	suite.Suite
	monthlySummaryUseCase usecase.MonthlySummaryUseCase
}

func TestMonthlySummaryUseCaseSuite(t *testing.T) {
	suite.Run(t, new(MonthlySummaryUseCaseSuite))
}

func (suite *MonthlySummaryUseCaseSuite) SetupTest() {
	mockRepo := NewMockMonthlySummaryRepository()
	suite.monthlySummaryUseCase = usecase.NewMonthlySummaryUseCase(mockRepo)
}

func (suite *MonthlySummaryUseCaseSuite) TestCreateMonthlySummary() {
	summary := &entity.MonthlySummary{
		UserID:    1,
		YearMonth: "2025-01",
		Income:    5000.00,
		Expense:   3000.00,
		Balance:   2000.00,
	}

	mockRepo := NewMockMonthlySummaryRepository()
	suite.monthlySummaryUseCase = usecase.NewMonthlySummaryUseCase(mockRepo)
	mockRepo.On("CreateMonthlySummary", summary).Return(summary, nil)

	createdSummary, err := suite.monthlySummaryUseCase.CreateMonthlySummary(summary)
	suite.Assert().Nil(err)
	suite.Assert().Equal("2025-01", createdSummary.YearMonth)
	suite.Assert().Equal(2000.00, createdSummary.Balance)
}

func (suite *MonthlySummaryUseCaseSuite) TestGetMonthlySummaryByID() {
	summary := &entity.MonthlySummary{
		ID:        1,
		UserID:    1,
		YearMonth: "2025-01",
		Income:    5000.00,
		Expense:   3000.00,
		Balance:   2000.00,
	}

	mockRepo := NewMockMonthlySummaryRepository()
	suite.monthlySummaryUseCase = usecase.NewMonthlySummaryUseCase(mockRepo)
	mockRepo.On("GetMonthlySummaryByID", summary.ID).Return(summary, nil)

	retrievedSummary, err := suite.monthlySummaryUseCase.GetMonthlySummaryByID(summary.ID)
	suite.Assert().Nil(err)
	suite.Assert().Equal("2025-01", retrievedSummary.YearMonth)
	suite.Assert().Equal(2000.00, retrievedSummary.Balance)
}

func (suite *MonthlySummaryUseCaseSuite) TestGetMonthlySummariesByUserID() {
	summaries := []entity.MonthlySummary{
		{
			ID:        1,
			UserID:    1,
			YearMonth: "2025-01",
			Income:    5000.00,
			Expense:   3000.00,
			Balance:   2000.00,
		},
		{
			ID:        2,
			UserID:    1,
			YearMonth: "2025-02",
			Income:    5500.00,
			Expense:   3500.00,
			Balance:   2000.00,
		},
	}

	mockRepo := NewMockMonthlySummaryRepository()
	suite.monthlySummaryUseCase = usecase.NewMonthlySummaryUseCase(mockRepo)
	mockRepo.On("GetMonthlySummariesByUserID", 1).Return(summaries, nil)

	retrievedSummaries, err := suite.monthlySummaryUseCase.GetMonthlySummariesByUserID(1)
	suite.Assert().Nil(err)
	suite.Assert().Equal(2, len(retrievedSummaries))
	suite.Assert().Equal("2025-01", retrievedSummaries[0].YearMonth)
	suite.Assert().Equal("2025-02", retrievedSummaries[1].YearMonth)
}

func (suite *MonthlySummaryUseCaseSuite) TestUpdateMonthlySummary() {
	summary := &entity.MonthlySummary{
		ID:        1,
		UserID:    1,
		YearMonth: "2025-01",
		Income:    6000.00,
		Expense:   4000.00,
		Balance:   2000.00,
	}

	mockRepo := NewMockMonthlySummaryRepository()
	suite.monthlySummaryUseCase = usecase.NewMonthlySummaryUseCase(mockRepo)
	mockRepo.On("UpdateMonthlySummary", summary).Return(summary, nil)

	updatedSummary, err := suite.monthlySummaryUseCase.UpdateMonthlySummary(summary)
	suite.Assert().Nil(err)
	suite.Assert().Equal(6000.00, updatedSummary.Income)
	suite.Assert().Equal(4000.00, updatedSummary.Expense)
	suite.Assert().Equal(2000.00, updatedSummary.Balance)
}

func (suite *MonthlySummaryUseCaseSuite) TestDeleteMonthlySummary() {
	mockRepo := NewMockMonthlySummaryRepository()
	suite.monthlySummaryUseCase = usecase.NewMonthlySummaryUseCase(mockRepo)
	mockRepo.On("DeleteMonthlySummary", 1).Return(nil)

	err := suite.monthlySummaryUseCase.DeleteMonthlySummary(1)
	suite.Assert().Nil(err)
}
