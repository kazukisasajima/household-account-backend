package handler

import (
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/oapi-codegen/runtime/types"

	"household-account-backend/adapter/controller/echo/presenter"
	"household-account-backend/entity"
	"household-account-backend/pkg/logger"
	"household-account-backend/usecase"
)

type UserHandler struct {
	userUseCase usecase.UserUseCase
}

func NewUserHandler(userUseCase usecase.UserUseCase) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
	}
}

func userToResponse(user *entity.User) *presenter.UserResponse {
	return &presenter.UserResponse{
		Id:        user.ID,
		Email:     types.Email(user.Email),
		Name:      user.Name,
	}
}

func (u *UserHandler) Signup(c echo.Context) error {
	var requestBody presenter.CreateUserJSONRequestBody
	if err := c.Bind(&requestBody); err != nil {
		logger.Warn(err.Error())
		return c.JSON(http.StatusBadRequest, &presenter.ErrorResponse{Message: err.Error()})
	}

	user := &entity.User{
		Email:    string(requestBody.Email),
		Password: requestBody.Password,
		Name:     requestBody.Name,
	}

	createdUser, err := u.userUseCase.Signup(user)
	if err != nil {
		logger.Error(err.Error())
		return c.JSON(http.StatusInternalServerError, &presenter.ErrorResponse{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, userToResponse(createdUser))
}

func (u *UserHandler) Login(c echo.Context) error {
	logger.Info("Loginが呼ばれた")
	var credentials entity.Credentials
	if err := c.Bind(&credentials); err != nil {
		logger.Warn(err.Error())
		return c.JSON(http.StatusBadRequest, &presenter.ErrorResponse{Message: err.Error()})
	}

	tokenString, err := u.userUseCase.Login(&credentials)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	// Set JWT token as a secure cookie
	cookie := new(http.Cookie)
	cookie.Name = "auth_token"
	cookie.Value = tokenString
	cookie.Expires = time.Now().Add(24 * time.Hour)
	cookie.Path = "/"
	cookie.Domain = os.Getenv("API_DOMAIN")
	// cookie.Secure = true
	cookie.HttpOnly = true
	cookie.SameSite = http.SameSiteNoneMode
	c.SetCookie(cookie)
	logger.Info("cookieの中身: " + cookie.Value)

	return c.NoContent(http.StatusOK)
}

func (u *UserHandler) Logout(c echo.Context) error {
	cookie := new(http.Cookie)
	cookie.Name = "auth_token"
	cookie.Value = ""
	cookie.Expires = time.Now().Add(-1 * time.Hour)
	cookie.Path = "/"
	cookie.Domain = os.Getenv("API_DOMAIN")
	// cookie.Secure = trueas
	cookie.HttpOnly = true
	cookie.SameSite = http.SameSiteNoneMode
	c.SetCookie(cookie)

	return c.NoContent(http.StatusOK)
}

func (u *UserHandler) CsrfToken(c echo.Context) error {
	secret := os.Getenv("SECRET")
	logger.Info("secret: " + secret)

	token := c.Get("csrf").(string)
	return c.JSON(http.StatusOK, echo.Map{
		"csrf_token": token,
	})
}

func (u *UserHandler) GetCurrentUser(c echo.Context) error {
	logger.Info("GetCurrentUserが呼ばれた")
	// JWTトークンからユーザーIDを取得
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userId := int(claims["user_id"].(float64)) // トークンから取得した user_id

	// ユースケースからユーザー情報を取得
	userEntity, err := u.userUseCase.GetCurrentUser(userId)
	if err != nil {
		logger.Error(err.Error())
		return c.JSON(http.StatusInternalServerError, &presenter.ErrorResponse{Message: "Failed to retrieve user"})
	}

	// レスポンスとしてユーザー情報を返す
	return c.JSON(http.StatusOK, userToResponse(userEntity))
}

func (u *UserHandler) UpdateUser(c echo.Context) error {
	// リクエストボディを取得
	var requestBody presenter.UpdateCurrentUserJSONRequestBody
	if err := c.Bind(&requestBody); err != nil {
		logger.Warn(err.Error())
		return c.JSON(http.StatusBadRequest, &presenter.ErrorResponse{Message: "Invalid request format"})
	}

	// JWTトークンからユーザーIDを取得
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userId := int(claims["user_id"].(float64))

	// 更新対象ユーザーのエンティティを作成
	userEntity := &entity.User{
		ID:    userId,
		Email: string(requestBody.Email),
		Name:  requestBody.Name,
	}

	// パスワードが提供されている場合はハッシュ化
	if requestBody.Password != "" {
		hashedPassword, err := usecase.HashPassword(requestBody.Password)
		if err != nil {
			logger.Error("Password hashing failed: " + err.Error())
			return c.JSON(http.StatusInternalServerError, &presenter.ErrorResponse{Message: "Failed to hash password"})
		}
		userEntity.Password = hashedPassword
	}

	// ユーザーを更新
	updatedUser, err := u.userUseCase.UpdateUser(userEntity)
	if err != nil {
		logger.Error("Failed to update user: " + err.Error())
		return c.JSON(http.StatusInternalServerError, &presenter.ErrorResponse{Message: "Failed to update user"})
	}

	// 更新後のユーザー情報をレスポンスとして返す
	return c.JSON(http.StatusOK, userToResponse(updatedUser))
}

func (u *UserHandler) DeleteUser(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userId := int(claims["user_id"].(float64))

	if err := u.userUseCase.DeleteUser(userId); err != nil {
		logger.Error(err.Error())
		return c.JSON(http.StatusInternalServerError, &presenter.ErrorResponse{Message: err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}
