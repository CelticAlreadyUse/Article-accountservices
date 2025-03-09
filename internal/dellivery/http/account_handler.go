package httphandler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/CelticAlreadyUse/Article-accountservices/internal/model"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type AccountHandler struct {
	accountUsecase model.AccountUsecase
	otpUsecase     model.OTPUsecase
}

func InitAccountHandler(accountUsecase model.AccountUsecase, otpUsecase model.OTPUsecase) *AccountHandler {
	return &AccountHandler{accountUsecase: accountUsecase, otpUsecase: otpUsecase}
}

var validate = validator.New()

func (handler *AccountHandler) RegisterAccountHandler(e *echo.Echo) {
	g := e.Group("/v1/auth")
	g.GET("/account/:id", handler.show, AuthMiddleWare)
	g.POST("/register", handler.register)
	g.POST("/login", handler.login)
	g.POST("/update/:id", handler.update, AuthMiddleWare)
	g.GET("/search", handler.findUsername, AuthMiddleWare)
	g.GET("/otp/request", handler.requestOTP)
	g.POST("/otp/validate", handler.validateOTP)
	g.POST("/send-otp/password", handler.sendOTPPass)
	g.POST("/verify-otp/password", handler.verifyOTPPass)
	g.POST("/reset-password", handler.resetPassword, OTPMiddleWare)
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
	var User = model.Login{
		ID: login.ID,
		Email: login.Email,
		Username: login.Username,
		Role: login.Role,
	}
	response := Response{
		AccesToken: login.Token,
		Message: "login sucessfully",
		Data: User,
	}
	return e.JSON(http.StatusAccepted, response)

}
func (handler *AccountHandler) show(e echo.Context) error {
	idParam := e.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id format")
	}
	claim, ok := e.Request().Context().Value(model.BearerAuthKey).(model.CustomClaims)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}
	log.Printf("Authenticated User ID : %d", claim.UserID)

	account, err := handler.accountUsecase.FindByID(e.Request().Context(), model.Account{}, id)
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
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	var data *model.Account
	err = e.Bind(&data)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	account, err := handler.accountUsecase.Update(e.Request().Context(), *data, int64(idInt))
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
