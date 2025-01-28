package router

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"github.com/swaggo/swag"
	"gorm.io/gorm"

	mymiddleware "household-account-backend/adapter/controller/echo/middleware"
	"household-account-backend/adapter/controller/echo/handler"
	"household-account-backend/adapter/controller/echo/presenter"
	"household-account-backend/adapter/gateway"
	"household-account-backend/pkg"
	"household-account-backend/pkg/logger"
	"household-account-backend/usecase"
)

// Swagger の設定
func setupSwagger(router *echo.Echo) (*openapi3.T, error) {
	swagger, err := presenter.GetSwagger()
	if err != nil {
		return nil, err
	}

	env := pkg.GetEnvDefault("APP_ENV", "development")
	if env == "development" {
		swaggerJson, _ := json.Marshal(swagger)
		var SwaggerInfo = &swag.Spec{
			InfoInstanceName: "swagger",
			SwaggerTemplate:  string(swaggerJson),
		}
		swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
		router.GET("/swagger/*", echoSwagger.WrapHandler)
	}

	return swagger, nil
}

// Echo 用のルータを作成。
func NewEchoRouter(db *gorm.DB) *echo.Echo {
	router := echo.New()

	// ミドルウェア設定
	router.Use(mymiddleware.CustomRequestLogger())
	router.Use(mymiddleware.CustomRecovery())
	router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000", os.Getenv("FE_URL")},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAccessControlAllowHeaders, echo.HeaderXCSRFToken},
		AllowMethods:     []string{"GET", "PUT", "PATCH", "POST", "DELETE"},
		// AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		AllowCredentials: true,
	}))
	router.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		CookiePath:     "/",
		CookieDomain:   os.Getenv("API_DOMAIN"),
		CookieHTTPOnly: true,
		// CookieSameSite: http.SameSiteNoneMode,
		CookieSameSite: http.SameSiteDefaultMode,
		// CookieMaxAge:   60,
	}))

	// Swagger の設定
	_, err := setupSwagger(router)
	if err != nil {
		logger.Warn("Swagger setup error: " + err.Error())
	}

	userRepository := gateway.NewUserRepository(db)
	userUseCase := usecase.NewUserUseCase(userRepository)
	userHandler := handler.NewUserHandler(userUseCase)

	categoryRepository := gateway.NewCategoryRepository(db)
	categoryUseCase := usecase.NewCategoryUseCase(categoryRepository)
	categoryHandler := handler.NewCategoryHandler(categoryUseCase)

	transactionRepository := gateway.NewTransactionRepository(db)
	transactionUseCase := usecase.NewTransactionUseCase(transactionRepository)
	transactionHandler := handler.NewTransactionHandler(transactionUseCase)

	monthlySummaryRepository := gateway.NewMonthlySummaryRepository(db)
	monthlySummaryUseCase := usecase.NewMonthlySummaryUseCase(monthlySummaryRepository)
	monthlySummaryHandler := handler.NewMonthlySummaryHandler(monthlySummaryUseCase)

	// ユーザー用エンドポイント
	users := router.Group("/api/v1/users")
	users.Use(mymiddleware.JWTMiddleware())
	users.GET("", userHandler.GetCurrentUser)
	users.PATCH("", userHandler.UpdateUser)
	users.DELETE("", userHandler.DeleteUser)

	// 認証用エンドポイント
	auth := router.Group("/api/v1/auth")
	auth.POST("/login", userHandler.Login)
	auth.POST("/signup", userHandler.Signup)
	auth.POST("/logout", userHandler.Logout)
	auth.GET("/csrf", userHandler.CsrfToken)

	// カテゴリー用エンドポイント
	categories := router.Group("/api/v1/categories")
	categories.Use(mymiddleware.JWTMiddleware())
	categories.GET("", categoryHandler.GetCategoriesByUserID)
	categories.POST("", categoryHandler.CreateCategory)
	categories.GET("/:id", categoryHandler.GetCategoryByID)
	categories.PATCH("/:id", categoryHandler.UpdateCategory)
	categories.DELETE("/:id", categoryHandler.DeleteCategory)

	// 取引用エンドポイント
	transactions := router.Group("/api/v1/transactions")
	transactions.Use(mymiddleware.JWTMiddleware())
	transactions.GET("", transactionHandler.GetTransactionsByUserID)
	transactions.POST("", transactionHandler.CreateTransaction)
	transactions.GET("/:id", transactionHandler.GetTransactionByID)
	transactions.PATCH("/:id", transactionHandler.UpdateTransaction)
	transactions.DELETE("/:id", transactionHandler.DeleteTransaction)

	// 月次集計用エンドポイント
	monthlySummaries := router.Group("/api/v1/monthly_summaries")
	monthlySummaries.Use(mymiddleware.JWTMiddleware())
	monthlySummaries.GET("", monthlySummaryHandler.GetMonthlySummariesByUserID)
	monthlySummaries.POST("", monthlySummaryHandler.CreateMonthlySummary)
	monthlySummaries.GET("/:id", monthlySummaryHandler.GetMonthlySummaryByID)
	monthlySummaries.PATCH("/:id", monthlySummaryHandler.UpdateMonthlySummary)
	monthlySummaries.DELETE("/:id", monthlySummaryHandler.DeleteMonthlySummary)

	// Swagger やその他のルート
	// router.GET("/", handler.Index)
	router.GET("/health", handler.Health)

	return router
}