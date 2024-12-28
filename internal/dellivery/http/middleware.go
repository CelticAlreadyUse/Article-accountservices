package httphandler

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/CelticAlreadyUse/Article-accountservices/internal/helper"
	"github.com/CelticAlreadyUse/Article-accountservices/internal/model"
	"github.com/labstack/echo/v4"
)



func MiddleWare(next echo.HandlerFunc) echo.HandlerFunc{
	return func (c echo.Context) error {
		authHandler := c.Request().Header.Get(echo.HeaderAuthorization)
		splitAuth := strings.Split(authHandler," ")
		if len(splitAuth) !=2 || splitAuth[0] != "Bearer" {
			return echo.NewHTTPError(http.StatusUnauthorized,"invalid token")
		}

		accessToken := splitAuth[1]
		log.Println(accessToken)

		var claim model.CustomClaims
		err := helper.DecodeToken(accessToken,&claim)
		if err !=nil{
			log.Println(claim)
			return echo.NewHTTPError(http.StatusUnauthorized,"invalid token")		
		}
			ctx := context.WithValue(c.Request().Context(),model.BearerAuthKey,claim)
			req := c.Request().WithContext(ctx)
			c.SetRequest(req)
			
			return next(c)
		}

}