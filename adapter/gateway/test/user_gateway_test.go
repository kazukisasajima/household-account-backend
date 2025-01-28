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

type UserRepositorySuite struct {
	tester.DBSQLiteSuite
	repository gateway.UserRepository
}

func TestUserRepositorySuite(t *testing.T) {
	suite.Run(t, new(UserRepositorySuite))
}

func (suite *UserRepositorySuite) SetupSuite() {
	suite.DBSQLiteSuite.SetupSuite()
	suite.repository = gateway.NewUserRepository(suite.DB)
}

func (suite *UserRepositorySuite) MockDB() sqlmock.Sqlmock {
	mock, mockGormDB := tester.MockDB()
	suite.repository = gateway.NewUserRepository(mockGormDB)
	return mock
}

func (suite *UserRepositorySuite) AfterTest(suiteName, testName string) {
	suite.repository = gateway.NewUserRepository(suite.DB)
}

func (suite *UserRepositorySuite) TestUserRepositoryCRUD() {
	user := &entity.User{
		Email:    "test@example.com",
		Password: "password",
		Name: 	  "Jhone",
	}
	createdUser, err := suite.repository.Signup(user)
	suite.Assert().Nil(err)
	suite.Assert().NotZero(createdUser.ID)
	suite.Assert().Equal("test@example.com", createdUser.Email)
	suite.Assert().Equal("password", createdUser.Password)
	suite.Assert().Equal("Jhone", createdUser.Name)

	getUser, err := suite.repository.GetCurrentUser(createdUser.ID)
	suite.Assert().Nil(err)
	suite.Assert().Equal("test@example.com", getUser.Email)
	suite.Assert().Equal("password", getUser.Password)
	suite.Assert().Equal("Jhone", getUser.Name)

	getUser.Name = "Doe"
	updatedUser, err := suite.repository.UpdateUser(getUser)
	suite.Assert().Nil(err)
	suite.Assert().Equal("Doe", updatedUser.Name)
	suite.Assert().Equal("test@example.com", updatedUser.Email)
	suite.Assert().Equal("password", updatedUser.Password)

	err = suite.repository.DeleteUser(createdUser.ID)
	suite.Assert().Nil(err)
	deletedUser, err := suite.repository.GetCurrentUser(createdUser.ID)
	suite.Assert().Nil(deletedUser)
	suite.Assert().Equal("record not found", err.Error())
}

func (suite *UserRepositorySuite) TestUserCreateFailure() {
	mockDB := suite.MockDB()
	mockDB.ExpectBegin()
	mockDB.ExpectExec(regexp.QuoteMeta("INSERT INTO `users` (`email`,`password`,`name`) VALUES (?,?,?)")).
		WithArgs("fail@example.com", "password", "Jhon").
		WillReturnError(errors.New("create error"))
	mockDB.ExpectRollback()

	user := &entity.User{
		Email: "fail@example.com", 
		Password: "password",
		Name: "Jhon",
	}

	createdUser, err := suite.repository.Signup(user)
	suite.Assert().Nil(createdUser)
	suite.Assert().NotNil(err)
	suite.Assert().Equal("create error", err.Error())
}

func (suite *UserRepositorySuite) TestUserGetFailure() {
	mockDB := suite.MockDB()
	mockDB.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE `users`.`id` = ? ORDER BY `users`.`id` LIMIT ?")).
		WithArgs(1, 1).
		WillReturnError(errors.New("get error"))

	user, err := suite.repository.GetCurrentUser(1)
	suite.Assert().Nil(user)
	suite.Assert().NotNil(err)
	suite.Assert().Equal("get error", err.Error())
}

func (suite *UserRepositorySuite) TestUserUpdateFailure() {
	mockDB := suite.MockDB()
	mockDB.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE `users`.`id` = ? ORDER BY `users`.`id` LIMIT ?")).
		WithArgs(1, 1).
		WillReturnError(errors.New("update error"))

	user := &entity.User{
		ID: 1,
		Email: "test@example.com",
		Password: "password",
		Name: "Jhon",
	}

	user, err := suite.repository.UpdateUser(user)
	suite.Assert().Nil(user)
	suite.Assert().NotNil(err)
	suite.Assert().Equal("update error", err.Error())
}

func (suite *UserRepositorySuite) TestUserDeleteFailure() {
	mockDB := suite.MockDB()
	mockDB.ExpectBegin()
	mockDB.ExpectExec(regexp.QuoteMeta("DELETE FROM `users` WHERE `users`.`id` = ?")).WithArgs(1).
		WillReturnError(errors.New("delete error"))
	mockDB.ExpectRollback()

	err := suite.repository.DeleteUser(1)
	suite.Assert().NotNil(err)
	suite.Assert().Equal("delete error", err.Error())
}
