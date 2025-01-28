package gateway

import (
	"github.com/jinzhu/copier"
	"gorm.io/gorm"

	"household-account-backend/entity"
)

type CategoryRepository interface {
	CreateCategory(category *entity.Category) (*entity.Category, error)
	GetCategoryByID(userID int, categoryID int) (*entity.Category, error)
	GetCategoriesByUserID(userID int) ([]entity.Category, error)
	UpdateCategory(category *entity.Category) (*entity.Category, error)
	DeleteCategory(userID int, categoryID int) error
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db}
}

func (cr *categoryRepository) CreateCategory(category *entity.Category) (*entity.Category, error) {
	if err := cr.db.Create(category).Error; err != nil {
		return nil, err
	}
	return category, nil
}

func (cr *categoryRepository) GetCategoryByID(userID int, categoryID int) (*entity.Category, error) {
	category := &entity.Category{}
	if err := cr.db.Where("id = ? AND user_id = ?", categoryID, userID).First(category).Error; err != nil {
		return nil, err
	}
	return category, nil
}

func (cr *categoryRepository) GetCategoriesByUserID(userID int) ([]entity.Category, error) {
	var categories []entity.Category
	if err := cr.db.Where("user_id = ?", userID).Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (cr *categoryRepository) UpdateCategory(category *entity.Category) (*entity.Category, error) {
	// 既存データの取得
	selectedCategory, err := cr.GetCategoryByID(category.UserID, category.ID)
	if err != nil {
		return nil, err
	}

	// フィールドをコピー（空の値を無視）
	if err := copier.CopyWithOption(selectedCategory, category, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return nil, err
	}

	// 更新
	if err := cr.db.Save(selectedCategory).Error; err != nil {
		return nil, err
	}

	return selectedCategory, nil
}

func (cr *categoryRepository) DeleteCategory(userID int, categoryID int) error {
	if err := cr.db.Where("id = ? AND user_id = ?", categoryID, userID).Delete(&entity.Category{}).Error; err != nil {
		return err
	}
	return nil
}
