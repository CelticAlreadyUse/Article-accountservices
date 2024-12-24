package httphandler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/CelticAlreadyUse/Article-accountservices/internal/model"
	"github.com/labstack/echo/v4"
)

type AccountHandler struct {
	accountUsecase model.AccountUsecase
}

func InitAccountHandler(usecase model.AccountUsecase) *AccountHandler {
	return &AccountHandler{accountUsecase: usecase}
}
func (handler *AccountHandler) RegisterAccountHandler(e *echo.Echo) {
	g := e.Group("/v1/auth")
	g.GET("/account/:id", handler.show, MiddleWare)
	g.POST("/register", handler.register)
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
		echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	token, err := handler.accountUsecase.Create(e.Request().Context(), *data)
	if err != nil {
		echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
		response := Response{
		AccesToken: token,
	}
	return e.JSON(http.StatusAccepted, response)
}
