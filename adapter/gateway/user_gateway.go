package gateway

import (
	"github.com/jinzhu/copier"
	"gorm.io/gorm"

	"household-account-backend/entity"
)

// entityを使用してデータベースにアクセスする
// 内側の層から外側の層のメソッドを呼び出そうとする場合、依存性逆転の原則に違反する
// インターフェースを通すことで内側の層からアクセスできるようにする

type UserRepository interface {
	Signup(user *entity.User) (*entity.User, error)
	GetCurrentUser(userId int) (*entity.User, error)
	UpdateUser(*entity.User) (*entity.User, error)
	DeleteUser(userId int) error
	GetUserByEmail(email string) (*entity.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (ur *userRepository) Signup(user *entity.User) (*entity.User, error) {
	if err := ur.db.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (ur *userRepository) GetCurrentUser(userId int) (*entity.User, error) {
	user := &entity.User{}
	if err := ur.db.First(user, userId).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (ur *userRepository) DeleteUser(userId int) error {
	if err := ur.db.Delete(&entity.User{}, userId).Error; err != nil {
		return err
	}
	return nil
}

func (ur *userRepository) UpdateUser(user *entity.User) (*entity.User, error) {
	selectedUser, err := ur.GetCurrentUser(user.ID)
	if err != nil {
		return nil, err
	}

	if err := copier.CopyWithOption(selectedUser, user, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return nil, err
	}
	if err := ur.db.Save(&selectedUser).Error; err != nil {
		return nil, err
	}

	return selectedUser, nil
}

func (ur *userRepository) GetUserByEmail(email string) (*entity.User, error) {
	user := entity.User{}
	if err := ur.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
