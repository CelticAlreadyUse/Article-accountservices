package httphandler

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/CelticAlreadyUse/Article-accountservices/internal/helper"
	"github.com/CelticAlreadyUse/Article-accountservices/internal/model"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type AccountHandler struct {
	accountUsecase model.AccountUsecase
	otpUsecase     model.OTPUsecase
}

var jwtKey = []byte("access-secret-key")
var refreshKey = []byte("refresh-secret-key")

func InitAccountHandler(accountUsecase model.AccountUsecase, otpUsecase model.OTPUsecase) *AccountHandler {
	return &AccountHandler{accountUsecase: accountUsecase, otpUsecase: otpUsecase}
}

var validate = validator.New()

func (handler *AccountHandler) RegisterAccountHandler(e *echo.Echo) {
	g := e.Group("/v1/auth")
	g.POST("/register", handler.register)
	g.GET("/account/:id", handler.show, AuthMiddleWare)
	g.POST("/login", handler.login)
	g.PUT("/update/:id", handler.update, AuthMiddleWare)
	g.GET("/search", handler.findUsername, AuthMiddleWare)
	g.GET("/email/verify/request", handler.requestEmailOTP)
	g.POST("/email/verify/confirm", handler.confirmEmailOTP)
	g.POST("/password/reset/request", handler.sendOTPPass)
	g.POST("/password/reset/verify", handler.verifyOTPPass)
	g.POST("/password/reset/confirm", handler.resetPassword, OTPMiddleWare)
	g.POST("/refresh", handler.refreshToken)
}
func (handler *AccountHandler) login(e echo.Context) error {
	var body *model.Login
	err := e.Bind(&body)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	err = validate.Struct(body)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	login, err := handler.accountUsecase.Login(e.Request().Context(), *body)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}
	refreshToken, err := helper.GenerateRefreshToken(login.ID)
	if err !=nil{
		return echo.NewHTTPError(http.StatusInternalServerError,Response{Error: "Could not generate refresh token"})
	}
	e.SetCookie(&http.Cookie{
		Name:     "access_token",
		Value:    login.Token,
		Path:     "/",
		Expires:  time.Now().Add(15 * time.Minute),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
	e.SetCookie(&http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
	return e.JSON(http.StatusOK, Response{
		Message: "login success",
	})
}
func (handler *AccountHandler) refreshToken(e echo.Context) error {
	refreshCookies, err := e.Cookie("refresh_token")
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, Response{Error: "refresh token missing"})
	}
	claims, err := helper.ValidateToken(refreshCookies.Value, model.ConfigJWT{})
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, Response{Error: "invalid refresh token"})
	}
	newAccesToken, err := helper.GenerateAccessToken(claims.UserID)
	e.SetCookie(&http.Cookie{
		Name:     "access_token",
		Value:    newAccesToken,
		Path:     "/",
		Expires:  time.Now().Add(15 * time.Minute),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
	return e.JSON(http.StatusOK,Response{
		Message: "refreshed",
	})
}
func (handler *AccountHandler) me (e echo.Context) error{
	cookie,err  := e.Cookie("access_token")
	if err !=nil{
		return echo.NewHTTPError(http.StatusUnauthorized, Response{
			Error: "Unauthorized token not found",
		})
	}
	claims,err :=helper.VerifyAccessToken(cookie.Value)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, Response{
			Error: "unauthorized: invalid token",
		})
	}
	userID,ok := claims["sub"].(string)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized: invalid user ID in token")
	}
	user,err := handler.accountUsecase.FindByID(e.Request().Context(),userID)
	if err !=nil{
		return echo.NewHTTPError(http.StatusBadRequest, Response{
			Message: "user Data wasn't found",
		})
	}
	return e.JSON(http.StatusOK,Response{
		Data: user,
		Message: "sucessfully get user ID1",
	})
}
func (handler *AccountHandler) show(e echo.Context) error {
	id := e.Param("id")
	claim, ok := e.Request().Context().Value(model.BearerAuthKey).(model.CustomClaims)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}
	log.Printf("Authenticated User ID : %d", claim.UserID)

	account, err := handler.accountUsecase.FindByID(e.Request().Context(),id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}
	return e.JSON(http.StatusOK, Response{
		Data: account,
	})
}
func (handler *AccountHandler) register(e echo.Context) error {
	var data *model.Register
	if err := e.Bind(&data); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	err := validate.Struct(data)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	token, err := handler.accountUsecase.Create(e.Request().Context(), *data)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	response := Response{
		AccesToken: token,
	}
	return e.JSON(http.StatusAccepted, response)
}
func (handler *AccountHandler) update(e echo.Context) error {
	id := e.Param("id")
	var data *model.Account
	err := e.Bind(&data)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	account, err := handler.accountUsecase.Update(e.Request().Context(), *data, id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return e.JSON(http.StatusAccepted, Response{
		Data: account,
	})
}
func (handler *AccountHandler) findUsername(c echo.Context) error {
	var param model.SearchParam
	if limitParam := c.QueryParam("limit"); limitParam != " " {
		intLimit, err := strconv.Atoi(limitParam)
		if err != nil || intLimit <= 0 {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid limit value")
		}
		param.Limit = int64(intLimit)
	}

	if searchParam := c.QueryParam("username"); searchParam != " " {
		param.Username = searchParam
	}
	account := handler.accountUsecase.Search(c.Request().Context(), param)
	if account == nil {
		return echo.NewHTTPError(http.StatusNotFound, "User not found")
	}
	return c.JSON(http.StatusOK, Response{
		Data: account,
	})
}
