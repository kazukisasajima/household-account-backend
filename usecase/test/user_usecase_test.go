package usecase_test

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"household-account-backend/entity"
	"household-account-backend/usecase"
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
	userUseCase usecase.UserUseCase
}

func TestUserUseCaseSuite(t *testing.T) {
	suite.Run(t, new(UserUseCaseSuite))
}

func (suite *UserUseCaseSuite) TestSignup() {
	email := "test@example.com"
	password := "password123"
	hashedPassword, _ := usecase.HashPassword(password)
	mockRepo := NewMockUserRepository()
	suite.userUseCase = usecase.NewUserUseCase(mockRepo)

	user := &entity.User{
		Email:    email,
		Password: password,
		Name:     "John",
	}

	mockRepo.On("Signup", mock.AnythingOfType("*entity.User")).Return(&entity.User{
		ID:       1,
		Email:    email,
		Password: hashedPassword,
		Name:     "John",
	}, nil)

	createdUser, err := suite.userUseCase.Signup(user)
	suite.Assert().Nil(err)
	suite.Assert().Equal(email, createdUser.Email)
	suite.Assert().True(usecase.CheckPasswordHash(password, createdUser.Password))
}

func (suite *UserUseCaseSuite) TestLogin_Success() {
	email := "test@example.com"
	password := "password123"
	hashedPassword, _ := usecase.HashPassword(password)
	mockRepo := NewMockUserRepository()
	suite.userUseCase = usecase.NewUserUseCase(mockRepo)

	mockRepo.On("GetUserByEmail", email).Return(&entity.User{
		ID:       1,
		Email:    email,
		Password: hashedPassword,
	}, nil)

	token, err := suite.userUseCase.Login(&entity.Credentials{Email: email, Password: password})
	suite.Assert().Nil(err)
	suite.Assert().NotEmpty(token)
}

func (suite *UserUseCaseSuite) TestLogin_InvalidCredentials() {
	email := "test@example.com"
	password := "wrongpassword"
	mockRepo := NewMockUserRepository()
	suite.userUseCase = usecase.NewUserUseCase(mockRepo)

	mockRepo.On("GetUserByEmail", email).Return(&entity.User{
		ID:       1,
		Email:    email,
		Password: "hashedpassword",
	}, nil)

	token, err := suite.userUseCase.Login(&entity.Credentials{Email: email, Password: password})
	suite.Assert().NotNil(err)
	suite.Assert().Empty(token)
	suite.Assert().EqualError(err, "invalid credentials")
}

func (suite *UserUseCaseSuite) TestGetCurrentUser() {
	userID := 1
	email := "test@example.com"
	name := "John"
	mockRepo := NewMockUserRepository()
	suite.userUseCase = usecase.NewUserUseCase(mockRepo)

	mockRepo.On("GetCurrentUser", userID).Return(&entity.User{
		ID:    userID,
		Email: email,
		Name:  name,
	}, nil)

	user, err := suite.userUseCase.GetCurrentUser(userID)
	suite.Assert().Nil(err)
	suite.Assert().Equal(userID, user.ID)
	suite.Assert().Equal(email, user.Email)
	suite.Assert().Equal(name, user.Name)
}

func (suite *UserUseCaseSuite) TestUpdateUser() {
	userID := 1
	email := "test@example.com"
	name := "John"
	mockRepo := NewMockUserRepository()
	suite.userUseCase = usecase.NewUserUseCase(mockRepo)

	mockRepo.On("UpdateUser", mock.AnythingOfType("*entity.User")).Return(&entity.User{
		ID:    userID,
		Email: email,
		Name:  name,
	}, nil)

	user := &entity.User{
		ID:    userID,
		Email: email,
		Name:  name,
	}

	updatedUser, err := suite.userUseCase.UpdateUser(user)
	suite.Assert().Nil(err)
	suite.Assert().Equal(userID, updatedUser.ID)
	suite.Assert().Equal(email, updatedUser.Email)
	suite.Assert().Equal(name, updatedUser.Name)
}

func (suite *UserUseCaseSuite) TestDeleteUser() {
	userID := 1
	mockRepo := NewMockUserRepository()
	suite.userUseCase = usecase.NewUserUseCase(mockRepo)

	mockRepo.On("DeleteUser", userID).Return(nil)

	err := suite.userUseCase.DeleteUser(userID)
	suite.Assert().Nil(err)
}
