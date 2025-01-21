package usecase

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"household-account-backend/entity"
)

type mockUserRepository struct {
	mock.Mock
}

func NewMockUserRepository() *mockUserRepository {
	return new(mockUserRepository)
}

func (m *mockUserRepository) Signup(user *entity.User) (*entity.User, error) {
	args := m.Called(user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *mockUserRepository) GetCurrentUser(ID int) (*entity.User, error) {
	args := m.Called(ID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *mockUserRepository) UpdateUser(user *entity.User) (*entity.User, error) {
	args := m.Called(user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *mockUserRepository) DeleteUser(ID int) error {
	args := m.Called(ID)
	return args.Error(0)
}

func (m *mockUserRepository) GetUserByEmail(email string) (*entity.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

type UserUseCaseSuite struct {
	suite.Suite
	userUseCase UserUseCase
}

func TestUserUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(UserUseCaseSuite))
}

func (suite *UserUseCaseSuite) TestSignup() {
	email := "test@example.com"
	password := "password123"
	hashedPassword, _ := HashPassword(password)
	mockUserRepository := NewMockUserRepository()
	suite.userUseCase = NewUserUseCase(mockUserRepository)

	user := &entity.User{
		Email:    email,
		Password: password,
		Name:    "Jhone",
	}

	mockUserRepository.On("Signup", mock.AnythingOfType("*entity.User")).Return(&entity.User{
		ID:       1,
		Email:    email,
		Password: hashedPassword,
		Name:    "Jhone",
	}, nil)

	createdUser, err := suite.userUseCase.Signup(user)
	suite.Assert().Nil(err)
	suite.Assert().Equal(email, createdUser.Email)
	suite.Assert().True(CheckPasswordHash(password, createdUser.Password))
}

func (suite *UserUseCaseSuite) TestLogin() {
	email := "test@example.com"
	password := "password123"
	hashedPassword, _ := HashPassword(password)
	mockUserRepository := NewMockUserRepository()
	suite.userUseCase = NewUserUseCase(mockUserRepository)

	users := &entity.User{
		Email:    email,
		Password: password,
	}

	mockUserRepository.On("GetUserByEmail", email).Return(&entity.User{
		ID:       1,
		Email:    email,
		Password: hashedPassword,
	}, nil)

	jwt, err := suite.userUseCase.Login(users)
	suite.Assert().Nil(err)
	suite.Assert().NotEmpty(jwt)
}

func (suite *UserUseCaseSuite) TestGetCurrentUser() {
	userID := 1
	email := "test@example.com"
	password := "password123"
	name := "Jhone"
	mockUserRepository := NewMockUserRepository()
	suite.userUseCase = NewUserUseCase(mockUserRepository)

	mockUserRepository.On("GetCurrentUser", userID).Return(&entity.User{
		ID:       userID,
		Email:    email,
		Password: password,
		Name:     name,
	}, nil)

	user, err := suite.userUseCase.GetCurrentUser(userID)
	suite.Assert().Nil(err)
	suite.Assert().Equal(userID, user.ID)
	suite.Assert().Equal(email, user.Email)
	suite.Assert().Equal(password, user.Password)
	suite.Assert().Equal(name, user.Name)
}

func (suite *UserUseCaseSuite) TestUpdateUser() {
	userID := 1
	email := "test@example.com"
	password := "password123"
	name := "Jhone"
	mockUserRepository := NewMockUserRepository()
	suite.userUseCase = NewUserUseCase(mockUserRepository)

	mockUserRepository.On("UpdateUser", mock.AnythingOfType("*entity.User")).Return(&entity.User{
		ID:       userID,
		Email:    email,
		Password: password,
		Name:     name,
	}, nil)
		
	user := &entity.User{
		ID: userID,
		Email: email,
		Password: password,
		Name: name,
	}

	updatedUser, err := suite.userUseCase.UpdateUser(user)
	suite.Assert().Nil(err)
	suite.Assert().Equal(userID, updatedUser.ID)
	suite.Assert().Equal(email, updatedUser.Email)
	suite.Assert().Equal(password, updatedUser.Password)
	suite.Assert().Equal(name, updatedUser.Name)
}

func (suite *UserUseCaseSuite) TestDeleteUser() {
	userID := 1
	mockUserRepository := NewMockUserRepository()
	suite.userUseCase = NewUserUseCase(mockUserRepository)

	mockUserRepository.On("DeleteUser", userID).Return(nil)

	err := suite.userUseCase.DeleteUser(userID)
	suite.Assert().Nil(err)
}
