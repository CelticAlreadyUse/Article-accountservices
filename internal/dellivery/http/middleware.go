package httphandler

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/CelticAlreadyUse/Article-accountservices/internal/config"
	redisClient "github.com/CelticAlreadyUse/Article-accountservices/internal/databases/redis"
	"github.com/CelticAlreadyUse/Article-accountservices/internal/helper"
	"github.com/CelticAlreadyUse/Article-accountservices/internal/model"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func AuthMiddleWare(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get(echo.HeaderAuthorization)
		splitAuth := strings.Split(authHeader, " ")
		if len(splitAuth) != 2 || splitAuth[0] != "Bearer" {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
		}
		accessToken := splitAuth[1]
		log.Println("Token:", accessToken)
		var claim model.CustomClaims
		token, err := jwt.ParseWithClaims(accessToken, &claim, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
			}
			return []byte(config.JWTSigningKey()), nil
		})
		if err != nil {
			if strings.Contains(err.Error(), "token is expired") {
				return echo.NewHTTPError(http.StatusUnauthorized, "token has expired")
			}
			log.Printf("Token Error: %v", err)
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
		}

		if !token.Valid {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
		}
		if helper.IsTokenExpired(claim.ExpiresAt) {
			return echo.NewHTTPError(http.StatusUnauthorized, "token has expired")
		}
		ctx := context.WithValue(c.Request().Context(), model.BearerAuthKey, claim)
		req := c.Request().WithContext(ctx)
		c.SetRequest(req)
		return next(c)
	}
}
func OTPMiddleWare(next echo.HandlerFunc)echo.HandlerFunc{
	return func(c echo.Context) error {
		resetToken := c.Request().Header.Get("Authorization")
		if resetToken != "" {
			echo.NewHTTPError(http.StatusUnauthorized,c.JSON(http.StatusUnauthorized,"Reset Token not found"))
		}

		email,err := redisClient.InitRedisClient().Get(context.Background(),"reset"+resetToken).Result()
		if err !=nil{
			return echo.NewHTTPError(http.StatusUnauthorized,"Reset token not valid or have been expired")
		}	

		c.Set("email",email)
		return next(c)

	}
}