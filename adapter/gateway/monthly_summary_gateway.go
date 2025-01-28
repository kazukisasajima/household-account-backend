package gateway

import (
	"gorm.io/gorm"

	"household-account-backend/entity"
)

type MonthlySummaryRepository interface {
	CreateMonthlySummary(summary *entity.MonthlySummary) (*entity.MonthlySummary, error)
	GetMonthlySummaryByID(userID int, summaryID int) (*entity.MonthlySummary, error)
	GetMonthlySummariesByUserID(userID int) ([]entity.MonthlySummary, error)
	UpdateMonthlySummary(summary *entity.MonthlySummary) (*entity.MonthlySummary, error)
	DeleteMonthlySummary(userID int, summaryID int) error
}

type monthlySummaryRepository struct {
	db *gorm.DB
}

func NewMonthlySummaryRepository(db *gorm.DB) MonthlySummaryRepository {
	return &monthlySummaryRepository{db}
}

func (msr *monthlySummaryRepository) CreateMonthlySummary(summary *entity.MonthlySummary) (*entity.MonthlySummary, error) {
	if err := msr.db.Create(summary).Error; err != nil {
		return nil, err
	}
	return summary, nil
}

func (msr *monthlySummaryRepository) GetMonthlySummaryByID(userID int, summaryID int) (*entity.MonthlySummary, error) {
	summary := &entity.MonthlySummary{}
	if err := msr.db.Where("id = ? AND user_id = ?", summaryID, userID).First(summary).Error; err != nil {
		return nil, err
	}
	return summary, nil
}

func (msr *monthlySummaryRepository) GetMonthlySummariesByUserID(userID int) ([]entity.MonthlySummary, error) {
	var summaries []entity.MonthlySummary
	if err := msr.db.Where("user_id = ?", userID).Find(&summaries).Error; err != nil {
		return nil, err
	}
	return summaries, nil
}

func (msr *monthlySummaryRepository) UpdateMonthlySummary(summary *entity.MonthlySummary) (*entity.MonthlySummary, error) {
	if err := msr.db.Save(summary).Error; err != nil {
		return nil, err
	}
	return summary, nil
}

func (msr *monthlySummaryRepository) DeleteMonthlySummary(userID int, summaryID int) error {
	if err := msr.db.Where("id = ? AND user_id = ?", summaryID, userID).Delete(&entity.MonthlySummary{}).Error; err != nil {
		return err
	}
	return nil
}
