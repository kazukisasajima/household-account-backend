package tester

import (
	"context"
	"fmt"
	"household-account-backend/entity"
	"household-account-backend/infrastructure/database"
	"household-account-backend/pkg"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/gorm"
)

type DBMySQLSuite struct {
	suite.Suite
	mySQLContainer testcontainers.Container
	ctx            context.Context
	DB             *gorm.DB
}

func (suite *DBMySQLSuite) SetupTestContainers() (err error) {
	configs := database.NewConfigMySQL()
	pkg.WaitForPort(configs.Database, configs.Port, 10*time.Second)
	suite.ctx = context.Background()
	req := testcontainers.ContainerRequest{
		Image: "mysql:8.2",
		Env: map[string]string{
			"MYSQL_DATABASE":             configs.Database,
			"MYSQL_USER":                 configs.User,
			"MYSQL_PASSWORD":             configs.Password,
			"MYSQL_ALLOW_EMPTY_PASSWORD": "yes",
		},
		ExposedPorts: []string{fmt.Sprintf("%s:3306/tcp", configs.Port)},
		WaitingFor:   wait.ForLog("port: 3306  MySQL Community Server"),
	}

	suite.mySQLContainer, err = testcontainers.GenericContainer(suite.ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	suite.Assert().Nil(err)
	return nil
}

func (suite *DBMySQLSuite) SetupSuite() {
	err := suite.SetupTestContainers()
	suite.Assert().Nil(err)

	db, err := database.NewDatabaseSQLFactory(database.InstanceMySQL)
	suite.Assert().Nil(err)
	suite.DB = db
	for _, model := range entity.NewDomains() {
		err = suite.DB.AutoMigrate(model)
		suite.Assert().Nil(err)
	}
}

func (suite *DBMySQLSuite) TearDownSuite() {
	if suite.mySQLContainer == nil {
		return
	}
	err := suite.mySQLContainer.Terminate(suite.ctx)
	suite.Assert().Nil(err)
}