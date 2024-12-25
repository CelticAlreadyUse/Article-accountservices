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
}

func InitAccountHandler(usecase model.AccountUsecase) *AccountHandler {
	return &AccountHandler{accountUsecase: usecase}
}
var validate = validator.New()
func (handler *AccountHandler) RegisterAccountHandler(e *echo.Echo) {
	g := e.Group("/v1/auth")
	g.GET("/account/:id", handler.show, MiddleWare)
	g.POST("/register", handler.register)
	g.POST("/login",handler.login)
	g.POST("/update/:id",handler.update)

}

func (handler *AccountHandler) login(e echo.Context) error {
	var body  *model.Login
	err := e.Bind(&body)
	if err !=nil{
	return echo.NewHTTPError(http.StatusBadRequest,err.Error())
	}	
	err = validate.Struct(body)
	if err !=nil{
		return echo.NewHTTPError(http.StatusBadRequest,err.Error())
	}
	token,err := handler.accountUsecase.Login(e.Request().Context(),*body)
	if err !=nil{
		return echo.NewHTTPError(http.StatusUnauthorized,err.Error())
	}
	response := Response{
		AccesToken: token,
	}
	return e.JSON(http.StatusAccepted,response)

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
	if err !=nil{
		return echo.NewHTTPError(http.StatusBadRequest,err.Error())
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

func (handler *AccountHandler)update(e echo.Context)error{
	panic("implement me")
}