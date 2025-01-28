package usecase

import (
	"errors"
	"household-account-backend/adapter/gateway"
	"household-account-backend/entity"
	"household-account-backend/pkg/logger"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type UserUseCase interface {
	Signup(user *entity.User) (*entity.User, error)
	Login(user *entity.Credentials) (string, error)
	GetCurrentUser(userId int) (*entity.User, error)
	UpdateUser(*entity.User) (*entity.User, error)
	DeleteUser(userId int) error
}

type userUseCase struct {
	userRepository gateway.UserRepository
}

func NewUserUseCase(userRepository gateway.UserRepository) UserUseCase {
	return &userUseCase{
		userRepository: userRepository,
	}
}

func (uu *userUseCase) Signup(user *entity.User) (*entity.User, error) {
	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		return nil, err
	}
	user.Password = hashedPassword

	return uu.userRepository.Signup(user)
}

func (uu *userUseCase) Login(user *entity.Credentials) (string, error) {
	// メールアドレスでユーザーを検索
	// 入力されたパスワードとDBに保存されているハッシュ化されたパスワードを比較
	// 以下のコードはテストの時エラーになる
	// user, err := uu.userRepository.GetUserByEmail(user.Email)
	// if err != nil || !CheckPasswordHash(user.Password, user.Password) {
	// 	return "", errors.New("invalid credentials")
	// }

	storedUser, err := uu.userRepository.GetUserByEmail(user.Email)
	if err != nil || !CheckPasswordHash(user.Password, storedUser.Password) {
		return "", errors.New("invalid credentials")
	}

	// ペイロードの作成
	claims := &jwt.MapClaims{
		"user_id": storedUser.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}

	// トークン生成
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	
	secret := os.Getenv("SECRET")
	logger.Info("os.Getenv(SECRET): " + secret)
	
	// トークンに署名を付与
	tokenString, err := token.SignedString([]byte(secret))
	logger.Info("tokenString: " + tokenString)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (uu *userUseCase) GetCurrentUser(userId int) (*entity.User, error) {
	return uu.userRepository.GetCurrentUser(userId)
}

func (uu *userUseCase) UpdateUser(user *entity.User) (*entity.User, error) {
	return uu.userRepository.UpdateUser(user)
}

func (uu *userUseCase) DeleteUser(userId int) error {
	return uu.userRepository.DeleteUser(userId)
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}