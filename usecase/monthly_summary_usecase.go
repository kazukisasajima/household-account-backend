package usecase

import (
	"household-account-backend/adapter/gateway"
	"household-account-backend/entity"
)

type MonthlySummaryUseCase interface {
	CreateMonthlySummary(summary *entity.MonthlySummary) (*entity.MonthlySummary, error)
	GetMonthlySummaryByID(userID int, summaryID int) (*entity.MonthlySummary, error)
	GetMonthlySummariesByUserID(userID int) ([]entity.MonthlySummary, error)
	UpdateMonthlySummary(summary *entity.MonthlySummary) (*entity.MonthlySummary, error)
	DeleteMonthlySummary(userID int, summaryID int) error
}

type monthlySummaryUseCase struct {
	monthlySummaryRepository gateway.MonthlySummaryRepository
}

func NewMonthlySummaryUseCase(monthlySummaryRepository gateway.MonthlySummaryRepository) MonthlySummaryUseCase {
	return &monthlySummaryUseCase{
		monthlySummaryRepository: monthlySummaryRepository,
	}
}

func (msu *monthlySummaryUseCase) CreateMonthlySummary(summary *entity.MonthlySummary) (*entity.MonthlySummary, error) {
	return msu.monthlySummaryRepository.CreateMonthlySummary(summary)
}

func (msu *monthlySummaryUseCase) GetMonthlySummaryByID(userID int, summaryID int) (*entity.MonthlySummary, error) {
	return msu.monthlySummaryRepository.GetMonthlySummaryByID(userID, summaryID)
}

func (msu *monthlySummaryUseCase) GetMonthlySummariesByUserID(userID int) ([]entity.MonthlySummary, error) {
	return msu.monthlySummaryRepository.GetMonthlySummariesByUserID(userID)
}

func (msu *monthlySummaryUseCase) UpdateMonthlySummary(summary *entity.MonthlySummary) (*entity.MonthlySummary, error) {
	return msu.monthlySummaryRepository.UpdateMonthlySummary(summary)
}

func (msu *monthlySummaryUseCase) DeleteMonthlySummary(userID int, summaryID int) error {
	return msu.monthlySummaryRepository.DeleteMonthlySummary(userID, summaryID)
}
